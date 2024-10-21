package grpcutil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"strconv"
	"time"

	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/internal_error"
	"google.golang.org/grpc"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	grpcstatus "google.golang.org/grpc/status"
)

var defaultMsgSize = math.MaxInt64

func NewGrpcClient(host string, port int, timeout time.Duration, tlsCfg *tls.Config) (*grpc.ClientConn, error) {
	var (
		err      error
		grpcConn *grpc.ClientConn
		cred     grpc.DialOption
	)

	if tlsCfg != nil {
		cred = grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg))
	} else {
		cred = grpc.WithInsecure()
	}
	duration := time.Duration(timeout)
	grpcConn, err = grpc.NewClient(fmt.Sprintf("%s:%d", host, port), cred, grpc.WithBlock(), grpc.WithTimeout(duration),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(defaultMsgSize), grpc.MaxCallRecvMsgSize(defaultMsgSize)))
	if err != nil {
		return nil, internal_error.ErrConnCannotOpen(host, port, err.Error())
	}
	return grpcConn, nil
}

func NewTLSConfig(host string, ca, cert, key string, peerName string, peerNameVerify bool) (*tls.Config, error) {
	if ca == "" {
		return nil, internal_error.ErrTLS("No CA certificate provide")
	}

	peer := peerName
	if !peerNameVerify {
		peer = ""
	} else if peer == "" {
		peer = host
	}

	tlsCfg := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         peer,
		MinVersion:         tls.VersionTLS12,
	}

	CAs := x509.NewCertPool()
	if ca, err := ioutil.ReadFile(ca); err == nil {
		if !CAs.AppendCertsFromPEM(ca) {
			return nil, internal_error.ErrTLS("AppendCertsFromPEM failed")
		}
		tlsCfg.RootCAs = CAs
	} else {
		return nil, internal_error.ErrTLS(err.Error())
	}

	if cert != "" || key != "" {
		if cert, err := tls.LoadX509KeyPair(cert, key); err != nil {
			return nil, internal_error.ErrTLS(err.Error())
		} else {
			tlsCfg.Certificates = []tls.Certificate{cert}
		}
	}

	tlsCfg.VerifyPeerCertificate = func(certificates [][]byte, _ [][]*x509.Certificate) error {
		certs := make([]*x509.Certificate, len(certificates))
		for i, data := range certificates {
			cert, err := x509.ParseCertificate(data)
			if err != nil {
				return internal_error.ErrTLS(err.Error())
			}
			certs[i] = cert
		}

		opts := x509.VerifyOptions{
			Roots:         tlsCfg.RootCAs,
			DNSName:       tlsCfg.ServerName,
			Intermediates: x509.NewCertPool(),
		}

		for _, cert := range certs[1:] {
			opts.Intermediates.AddCert(cert)
		}

		_, err := certs[0].Verify(opts)

		return err
	}

	return tlsCfg, nil
}

func GetGrpcError(address string, err error) error {
	var (
		host string
		port int
	)
	h, p, err2 := net.SplitHostPort(address)
	if err2 != nil {
		// ignore
	}
	port, err2 = strconv.Atoi(p)
	host = h
	if err2 != nil {
		//ignore
	}
	rpcErr, ok := grpcstatus.FromError(err)
	if !ok {
		return err
	}
	switch rpcErr.Code() {
	case grpccodes.DeadlineExceeded, grpccodes.Canceled:
		return internal_error.ErrConnRequestTimeout(host, port)
	case grpccodes.Unavailable:
		return internal_error.ErrConnUnavailable(host, port)
	}
	return err
}
