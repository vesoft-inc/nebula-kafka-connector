package com.vesoft.nebula.driver.graph.net;

import com.alibaba.fastjson.JSON;
import com.google.common.base.Charsets;
import com.google.protobuf.ByteString;
import com.vesoft.nebula.driver.graph.ErrorCode;
import com.vesoft.nebula.driver.graph.data.HostAddress;
import com.vesoft.nebula.driver.graph.exception.AuthFailedException;
import com.vesoft.nebula.driver.graph.exception.IOErrorException;
import com.vesoft.nebula.driver.graph.utils.ClientVersion;
import com.vesoft.nebula.driver.graph.utils.TlsUtil;
import com.vesoft.nebula.proto.common.ClientInfo;
import com.vesoft.nebula.proto.common.Common;
import com.vesoft.nebula.proto.graph.AuthRequest;
import com.vesoft.nebula.proto.graph.AuthResponse;
import com.vesoft.nebula.proto.graph.ExecuteRequest;
import com.vesoft.nebula.proto.graph.ExecuteResponse;
import com.vesoft.nebula.proto.graph.GraphServiceGrpc;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.netty.shaded.io.grpc.netty.NegotiationType;
import io.grpc.netty.shaded.io.grpc.netty.NettyChannelBuilder;
import io.grpc.netty.shaded.io.netty.handler.ssl.SslContext;
import java.nio.charset.Charset;
import java.util.Map;
import java.util.concurrent.TimeUnit;
import javax.net.ssl.SSLException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class GrpcConnection extends Connection {

    private static final Logger LOGGER = LoggerFactory.getLogger(GrpcConnection.class);

    private ManagedChannel                            channel;
    private GraphServiceGrpc.GraphServiceBlockingStub stub;
    private long                                      connectTimeout = 0;
    private long                                      requestTimeout = 0;

    private final Charset charset = Charsets.UTF_8;

    @Override
    public void open(HostAddress address,
                     NebulaClient.Builder builder) throws IOErrorException {
        this.serverAddr = address;
        this.connectTimeout = builder.connectTimeoutMills;
        this.requestTimeout = builder.requestTimeoutMills;
        if (builder.enableTls) {
            NettyChannelBuilder channelBuilder = NettyChannelBuilder
                    .forAddress(address.getHost(), address.getPort())
                    .useTransportSecurity()
                    .sslContext(TlsUtil.getSslContext(builder.tlsCa,
                                                      builder.tlsCert,
                                                      builder.tlsKey))
                    .maxInboundMessageSize(Integer.MAX_VALUE);
            if (builder.tlsPeerNameVerify) {
                channelBuilder.overrideAuthority(builder.tlsPeerName);
            }
            channel = channelBuilder.build();
        } else {
            channel = NettyChannelBuilder
                    .forAddress(address.getHost(), address.getPort())
                    .usePlaintext()
                    .maxInboundMessageSize(Integer.MAX_VALUE)
                    .build();
        }
        stub = GraphServiceGrpc.newBlockingStub(channel);
    }

    @Override
    public void close() {
        if (channel != null && !channel.isShutdown()) {
            channel.shutdown();
        }
        stub = null;
    }

    @Override
    public boolean ping(long sessionID, long timeoutMs) throws IOErrorException {
        ExecuteResponse response = execute(sessionID, "RETURN 1", timeoutMs);
        return ErrorCode.SUCCESSFUL_COMPLETION.code
                .equals(response.getStatus().getCode().toString(charset));
    }

    public AuthResult authenticate(String user, Map<String, Object> authOptions)
            throws AuthFailedException {
        try {
            ClientInfo clientInfo = ClientInfo.newBuilder()
                    .setLang(ClientInfo.Language.JAVA)
                    .setProtocolVersion(Common
                                                .getDescriptor()
                                                .getOptions()
                                                .getExtension(Common.protocolVersion))
                    .setVersion(ByteString.copyFrom(ClientVersion.clientVersion, charset))
                    .build();
            String authInfoString = JSON.toJSONString(authOptions);
            AuthRequest authReq = AuthRequest.newBuilder()
                    .setUsername(ByteString.copyFrom(user, charset))
                    .setAuthInfo(ByteString.copyFrom(authInfoString, charset))
                    .setClientInfo(clientInfo)
                    .build();

            AuthResponse resp = stub
                    .withDeadlineAfter(connectTimeout, TimeUnit.MILLISECONDS)
                    .authenticate(authReq);
            String code = resp.getStatus().getCode().toString(charset);
            if (!ErrorCode.SUCCESSFUL_COMPLETION.code.equals(code)) {
                throw new AuthFailedException(resp.getStatus().getMessage().toString(charset));
            }
            return new AuthResult(resp.getSessionId());
        } catch (Exception e) {
            // TODO
            throw e;
        }
    }

    public ExecuteResponse execute(long sessionID, String stmt, long timeout)
            throws IOErrorException {
        if (stmt == null) {
            throw new NullPointerException("statement is null.");
        }
        try {
            ExecuteRequest request = ExecuteRequest.newBuilder()
                    .setSessionId(sessionID)
                    .setStmt(ByteString.copyFrom(stmt, charset))
                    .build();

            return stub.withDeadlineAfter(timeout, TimeUnit.MILLISECONDS).execute(request);
        } catch (Exception e) {
            // TODO
            throw e;
        }
    }

    public ExecuteResponse execute(long sessionID, String stmt) throws IOErrorException {
        return execute(sessionID, stmt, this.requestTimeout);
    }
}
