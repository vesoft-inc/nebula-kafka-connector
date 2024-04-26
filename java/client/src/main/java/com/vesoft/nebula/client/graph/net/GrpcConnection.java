package com.vesoft.nebula.client.graph.net;

import com.alibaba.fastjson.JSON;
import com.google.common.base.Charsets;
import com.google.protobuf.ByteString;
import com.vesoft.nebula.client.graph.ErrorCode;
import com.vesoft.nebula.client.graph.data.HostAddress;
import com.vesoft.nebula.client.graph.exception.AuthFailedException;
import com.vesoft.nebula.client.graph.exception.IOErrorException;
import com.vesoft.nebula.client.graph.utils.ClientVersion;
import com.vesoft.nebula.proto.common.ClientInfo;
import com.vesoft.nebula.proto.common.Common;
import com.vesoft.nebula.proto.graph.AuthRequest;
import com.vesoft.nebula.proto.graph.AuthResponse;
import com.vesoft.nebula.proto.graph.ExecuteRequest;
import com.vesoft.nebula.proto.graph.ExecuteResponse;
import com.vesoft.nebula.proto.graph.GraphServiceGrpc;
import io.grpc.Deadline;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import java.nio.charset.Charset;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.locks.ReadWriteLock;
import java.util.concurrent.locks.ReentrantReadWriteLock;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class GrpcConnection extends Connection {

    private static final Logger LOGGER = LoggerFactory.getLogger(GrpcConnection.class);

    private static final ConcurrentHashMap<HostAddress, ManagedChannel> channels =
            new ConcurrentHashMap<>();
    private GraphServiceGrpc.GraphServiceBlockingStub stub;
    private long connectTimeout = 0;
    private long requestTimeout = 0;

    private final Charset charset = Charsets.UTF_8;

    private static final ReadWriteLock lock = new ReentrantReadWriteLock();

    @Override
    public void open(HostAddress address, long connectTimeout, long requestTimeout) {
        this.serverAddr = address;
        this.connectTimeout = connectTimeout;
        this.requestTimeout = requestTimeout;
        lock.readLock().lock();
        try {
            channels.computeIfAbsent(serverAddr, key -> createChannel());
        } finally {
            lock.readLock().unlock();
        }
        stub = GraphServiceGrpc.newBlockingStub(channels.get(serverAddr));
    }

    @Override
    public void close() {
        if (!channels.isEmpty()) {
            closeChannel();
        }
        stub = null;
    }

    @Override
    public boolean ping(long sessionID) throws IOErrorException {
        ExecuteResponse response = execute(sessionID, "RETURN 1");
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

            getChannel();
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
        getChannel();
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

    private void getChannel() {
        lock.readLock().lock();
        try {
            channels.computeIfAbsent(serverAddr, key -> {
                ManagedChannel channel = createChannel();
                stub = GraphServiceGrpc.newBlockingStub(channel)
                        .withDeadline(Deadline.after(requestTimeout, TimeUnit.MILLISECONDS));
                return channel;
            });
        } finally {
            lock.readLock().unlock();
        }
    }

    private ManagedChannel createChannel() {
        return ManagedChannelBuilder
                .forAddress(serverAddr.getHost(), serverAddr.getPort()).usePlaintext()
                .build();
    }

    private static void closeChannel() {
        lock.writeLock().lock();
        try {
            for (ManagedChannel channel : channels.values()) {
                if (channel != null && !channel.isShutdown()) {
                    try {
                        channel.shutdownNow().awaitTermination(5, TimeUnit.SECONDS);
                    } catch (InterruptedException e) {
                        LOGGER.warn("close grpc connection is interrupted.", e);
                    }
                }
            }
            channels.clear();
        } finally {
            lock.writeLock().unlock();
        }
    }
}
