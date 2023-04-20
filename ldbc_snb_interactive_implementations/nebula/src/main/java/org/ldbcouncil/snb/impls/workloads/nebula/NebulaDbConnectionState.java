package org.ldbcouncil.snb.impls.workloads.nebula;

import com.vesoft.nebula.client.graph.exception.AuthFailedException;
import com.vesoft.nebula.client.graph.exception.ClientServerIncompatibleException;
import com.vesoft.nebula.client.graph.exception.NotValidConnectionException;
import org.ldbcouncil.snb.impls.workloads.BaseDbConnectionState;
import org.ldbcouncil.snb.impls.workloads.QueryStore;

import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

import com.vesoft.nebula.client.graph.NebulaPoolConfig;
import com.vesoft.nebula.client.graph.net.NebulaPool;
import com.vesoft.nebula.client.graph.net.Session;
import com.vesoft.nebula.client.graph.data.HostAddress;

public class NebulaDbConnectionState<TDbQueryStore extends QueryStore> extends BaseDbConnectionState<TDbQueryStore>
{
    private final NebulaPool pool = new NebulaPool();
    private Session session = null;
    private final String username;
    private final String password;

    public NebulaDbConnectionState( Map<String, String> properties, TDbQueryStore store ) throws UnknownHostException {
        super(properties, store);

        final String endpointURI = properties.get( "endpoint" );
        username = properties.get( "user" );
        password = properties.get( "password" );

        List<HostAddress> addressList = new ArrayList<>();
        for (String address : endpointURI.split(",")) {
            String[] ip_port = address.split(":");
            addressList.add(new HostAddress(ip_port[0], Integer.parseInt(ip_port[1])));
        }

        NebulaPoolConfig nebulaPoolConfig = new NebulaPoolConfig();
        // (TODO) jmq
        nebulaPoolConfig.setMaxConnSize(1);

        pool.init(addressList, nebulaPoolConfig);
    }

    public Session getSession() throws AuthFailedException, ClientServerIncompatibleException, NotValidConnectionException {
        if (session == null) {
            try {
                session = pool.getSession(username, password, false);
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
        return session;
    }

    @Override
    public void close() {
        pool.close();
    }
}
