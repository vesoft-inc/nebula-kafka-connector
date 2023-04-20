package nrpc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

const (
	kHeaderLen      = 24
	kConnectTimeout = 2000 * time.Millisecond
	kMagicCode      = 0x414C424E
)

type Error interface {
	error
	Timeout() bool
	BadChannel() bool
}

type ErrTimeout struct {
	err error
}

type ErrBadChannel struct {
	err error
}

type ErrBadArg struct {
	err error
}

type Client struct {
	tcp  *net.TCPConn
	chid uint64
	lock sync.Mutex
	addr string
}

/**
 *  Create a Client object, which represents a TCP connection
 *  @addr   Target server address, following formats as:
 *          "localhost:9669",
 *          "127.0.0.1:9669",
 *          "[::ffff:127.0.0.1]:9669"
 */
func NewClient(addr string) *Client {
	client := &Client{
		addr: addr,
	}
	return client
}

/**
 *  Reconnect this client, should be invoked when some connection failure happened
 *  @attempts  Number of times to retry to connect
 *  @delay  Time duration to wait between each retry
 */
func (this *Client) Reconnect(attempts int, delay time.Duration) error {
	if attempts <= 0 {
		return &ErrBadArg{
			err: errors.New(fmt.Sprintf("Invalid argument: `attempts': %d", attempts)),
		}
	}
	this.lock.Lock()
	defer this.lock.Unlock()

	var err error
	for i := 0; i < attempts; i++ {
		time.Sleep(delay)
		if err = this.connect(); err == nil {
			break
		}
	}

	return err
}

/**
 * Close the client, should be invoked when this client is no longer needed
 */
func (this *Client) Close() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.tcp != nil {
		this.tcp.Close()
		this.tcp = nil
	}
}

/**
 * Send a request buffer and receive a response buffer.
 * Since a connection or protocol error is handled only when
 * performing I/O operations, and cannot be recovered automatically,
 * `Send' will retry the I/O operations upon ErrBadChannel, but only once.
 *
 * During network I/O, `Send' holds a lock to prevent multiple goroutines
 * accessing. If there are requirements for multiplexing, we could reimplement
 * it by using background goroutines for I/O
 *
 * @req     Request buffer
 * @timeout Time duration to timeout
 * @return  Response buffer and error
 *          When the error is an instance of ErrBadChannel,
 *          the caller should invoke `Reconnect' to switch to
 *          another underlying connection. Otherwise if it is
 *          a ErrTimeout, `Reconnect' is not necessary.
 */
func (this *Client) Send(req []byte, timeout time.Duration) ([]byte, error) {
	if req == nil || len(req) == 0 {
		return nil, &ErrBadArg{
			err: errors.New("`req' is nil or empty"),
		}
	}
	this.lock.Lock()
	defer this.lock.Unlock()

	retry := true

	if this.tcp == nil {
		err := this.connect()
		if err != nil {
			return nil, err
		}
		retry = false
	}

	this.chid++

RETRY:
	this.tcp.SetDeadline(time.Now().Add(timeout))

	if err := this.writeRequest(req, timeout); err != nil {
		if err.BadChannel() && retry {
			if err = this.connect(); err != nil {
				return nil, err
			}
			goto RETRY
		}
		return nil, err
	}

	if resp, err := this.readResponse(); err != nil {
		if err.BadChannel() && retry {
			retry = false
			if err = this.connect(); err != nil {
				return nil, err
			}
			goto RETRY
		}
		return nil, err
	} else {
		return resp, err
	}
}

func (this *Client) connect() Error {
	conn, err := net.DialTimeout("tcp", this.addr, kConnectTimeout)
	if err != nil {
		return mapNetError(err)
	}

	tcp := conn.(*net.TCPConn)

	tcp.SetNoDelay(true)
	tcp.SetKeepAlive(true)
	tcp.SetKeepAlivePeriod(60 * time.Second)

	if this.tcp != nil {
		this.tcp.Close()
	}
	this.tcp = tcp
	return nil
}

func (this *Client) writeRequest(req []byte, timeout time.Duration) Error {
	header := Header{
		magic:   kMagicCode,
		timeout: uint32(timeout.Milliseconds()),
		chid:    this.chid,
		bits:    uint64(len(req)) << 24,
	}
	hbuf := header.encode()

	buffers := net.Buffers(make([][]byte, 2))
	buffers[0] = hbuf
	buffers[1] = req

	if _, err := buffers.WriteTo(this.tcp); err != nil {
		return mapNetError(err)
	}

	return nil
}

func (this *Client) readResponse() ([]byte, Error) {
	hbuf := make([]byte, kHeaderLen)
	for {
		if _, err := io.ReadFull(this.tcp, hbuf); err != nil {
			return nil, mapNetError(err)
		}

		var header Header
		header.decode(hbuf)

		if header.magic != kMagicCode {
			return nil, &ErrBadChannel{
				err: errors.New("Corrupted Response: " + header.String()),
			}
		}

		resp := make([]byte, header.size())
		if _, err := io.ReadFull(this.tcp, resp); err != nil {
			return nil, mapNetError(err)
		}

		if header.chid == this.chid {
			return resp, nil
		}
	}
}

type Header struct {
	magic   uint32
	timeout uint32
	chid    uint64
	bits    uint64
}

func (this *Header) size() uint64 {
	return this.bits >> 24
}

func (this *Header) more() bool {
	return this.bits&0x1 != 0
}

func (this *Header) compressed() bool {
	return this.bits&0x2 != 0
}

func (this *Header) cmd() uint32 {
	return uint32((this.bits >> 2) & 0x7)
}

func (this *Header) String() string {
	return fmt.Sprintf("magic: 0x%x, timeout: %d, chid: %d, size: %d, bits: 0x%x",
		this.magic, this.timeout, this.chid, this.size(), this.bits)
}

func (this *Header) encode() []byte {
	bytes := make([]byte, kHeaderLen)
	binary.LittleEndian.PutUint32(bytes[0:], this.magic)
	binary.LittleEndian.PutUint32(bytes[4:], this.timeout)
	binary.LittleEndian.PutUint64(bytes[8:], this.chid)
	binary.LittleEndian.PutUint64(bytes[16:], this.bits)
	return bytes
}

func (this *Header) decode(bytes []byte) {
	this.magic = binary.LittleEndian.Uint32(bytes[0:])
	this.timeout = binary.LittleEndian.Uint32(bytes[4:])
	this.chid = binary.LittleEndian.Uint64(bytes[8:])
	this.bits = binary.LittleEndian.Uint64(bytes[16:])
}

func (this *ErrTimeout) Error() string {
	return this.err.Error()
}

func (this *ErrTimeout) Timeout() bool {
	return true
}

func (this *ErrTimeout) BadChannel() bool {
	return false
}

func (this *ErrBadChannel) Error() string {
	return this.err.Error()
}

func (this *ErrBadChannel) Timeout() bool {
	return false
}

func (this *ErrBadChannel) BadChannel() bool {
	return true
}

func (this *ErrBadArg) Error() string {
	return this.err.Error()
}

func (this *ErrBadArg) Timeout() bool {
	return false
}

func (this *ErrBadArg) BadChannel() bool {
	return false
}

func mapNetError(err error) Error {
	if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
		return &ErrTimeout{err: err}
	}
	return &ErrBadChannel{err: err}
}
