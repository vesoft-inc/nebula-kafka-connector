package com.vesoft.nebula.client.graph.net;

import com.google.common.base.Charsets;
import com.google.protobuf.ByteString;
import com.vesoft.nebula.client.graph.ErrorCode;
import com.vesoft.nebula.client.graph.data.HostAddress;
import com.vesoft.nebula.client.graph.exception.AuthFailedException;
import com.vesoft.nebula.client.graph.exception.ClientServerIncompatibleException;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.proto.AuthRequest;
import com.vesoft.nebula.proto.AuthResponse;
import com.vesoft.nebula.proto.ExecuteRequest;
import com.vesoft.nebula.proto.ExecuteResponse;
import com.vesoft.nebula.proto.GraphGrpc;
import com.vesoft.nebula.proto.SignoutRequest;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import java.nio.charset.Charset;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class GrpcConnection extends Connection {

    private static final Logger LOGGER = LoggerFactory.getLogger(GrpcConnection.class);

    private ManagedChannel channel;
    private GraphGrpc.GraphBlockingStub stub;
    private int connTimeout = 0;
    private int requestTimeout = 0;

    private Charset charset = Charsets.UTF_8;

    @Override
    public void open(HostAddress address, int connTimeout, int requestTimeout)
            throws IOErrorException {
        this.serverAddr = address;
        this.connTimeout = connTimeout <= 0 ? Integer.MAX_VALUE : connTimeout;
        this.requestTimeout = requestTimeout <= 0 ? Integer.MAX_VALUE : requestTimeout;
        channel = ManagedChannelBuilder
                .forAddress(address.getHost(), address.getPort()).usePlaintext()
                .build();
        stub = GraphGrpc.newBlockingStub(channel);
    }

    @Override
    public void reopen() throws IOErrorException, ClientServerIncompatibleException {
        // TODO
    }

    @Override
    public void close() {
        channel.shutdownNow();
    }

    @Override
    public boolean ping(long sessionID) throws IOErrorException {
        return false;
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
            String s = resp.getGqlStatus().getCode().toString(charset);
            if (!s.equals(ErrorCode.SUCCESSFUL_COMPLETION.code)) {
                throw new AuthFailedException(resp.getGqlStatus().toString());
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
        try {
            SignoutRequest request = SignoutRequest.newBuilder().setSessionId(sessionID).build();
            stub.signout(request);
        } catch (Exception e) {
            // TODO
            throw e;
        }
    }
}
