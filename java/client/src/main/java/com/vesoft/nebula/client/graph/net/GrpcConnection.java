package com.vesoft.nebula.client.graph.net;

import com.google.common.base.Charsets;
import com.google.protobuf.ByteString;
import com.vesoft.nebula.client.graph.ErrorCode;
import com.vesoft.nebula.client.graph.data.HostAddress;
import com.vesoft.nebula.client.graph.exception.AuthFailedException;
import com.vesoft.nebula.client.graph.exception.ClientServerIncompatibleException;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.proto.graph.AuthRequest;
import com.vesoft.nebula.proto.graph.AuthResponse;
import com.vesoft.nebula.proto.graph.ExecuteRequest;
import com.vesoft.nebula.proto.graph.ExecuteResponse;
import com.vesoft.nebula.proto.graph.GraphServiceGrpc;
import com.vesoft.nebula.proto.graph.GraphServiceGrpc;
import com.vesoft.nebula.proto.graph.SignoutRequest;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import java.nio.charset.Charset;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class GrpcConnection extends Connection {

    private static final Logger LOGGER = LoggerFactory.getLogger(GrpcConnection.class);

    private static ConcurrentHashMap<HostAddress, ManagedChannel> channels =
            new ConcurrentHashMap<>();
    private GraphServiceGrpc.GraphServiceBlockingStub stub;
    private int connTimeout = 0;
    private int requestTimeout = 0;

    private Charset charset = Charsets.UTF_8;

    @Override
    public void open(HostAddress address, int connTimeout, int requestTimeout)
            throws IOErrorException {
        this.serverAddr = address;
        this.connTimeout = connTimeout <= 0 ? Integer.MAX_VALUE : connTimeout;
        this.requestTimeout = requestTimeout <= 0 ? Integer.MAX_VALUE : requestTimeout;
        synchronized (GrpcConnection.class) {
            if (!channels.containsKey(address)) {
                ManagedChannel channel = ManagedChannelBuilder
                        .forAddress(address.getHost(), address.getPort()).usePlaintext()
                        .build();
                channels.put(address, channel);
            }
        }
        stub = GraphServiceGrpc.newBlockingStub(channels.get(address));
    }

    @Override
    public void close() {
        for (ManagedChannel channel : channels.values()) {
            if (channel != null && !channel.isShutdown()) {
                try {
                    channel.shutdownNow().awaitTermination(5, TimeUnit.SECONDS);
                } catch (InterruptedException e) {
                    LOGGER.warn("close grpc connection is interrupted.", e);
                }
            }
        }
        stub = null;
        channels.clear();
    }

    @Override
    public boolean ping(long sessionID) throws IOErrorException {
        ExecuteResponse response = execute(sessionID, "RETURN 1");
        return response.hasExecutionOutcome()
                && response.getExecutionOutcome().hasGqlStatus()
                && ErrorCode.SUCCESSFUL_COMPLETION.code.equals(
                response.getExecutionOutcome().getGqlStatus().getCode().toString(charset));
    }

    public AuthResult authenticate(String user, String password)
            throws AuthFailedException, IOErrorException {
        try {
            AuthRequest authReq = AuthRequest.newBuilder()
                    .setUsername(ByteString.copyFrom(user, charset))
                    .setPassword(ByteString.copyFrom(password, charset))
                    .setClientType(ByteString.copyFrom("java", charset))
                    .setClientVersion(ByteString.copyFrom("v5.0.0", charset))
                    .build();
            AuthResponse resp = stub.authenticate(authReq);
            String code = resp.getGqlStatus().getCode().toString(charset);
            if (!ErrorCode.SUCCESSFUL_COMPLETION.code.equals(code)) {
                throw new AuthFailedException(resp.getGqlStatus().getMessage().toString());
            }
            return new AuthResult(resp.getSessionId());
        } catch (Exception e) {
            // TODO
            throw e;
        }
    }

    public ExecuteResponse execute(long sessionID, String stmt) throws IOErrorException {
        try {
            ExecuteRequest request = ExecuteRequest.newBuilder()
                    .setSessionId(sessionID)
                    .setStmt(ByteString.copyFrom(stmt, charset))
                    .build();
            return stub.execute(request);
        } catch (Exception e) {
            // TODO
            throw e;
        }
    }

    public void signout(long sessionID) {
        SignoutRequest request = SignoutRequest.newBuilder().setSessionId(sessionID).build();
        stub.signout(request);
    }
}
