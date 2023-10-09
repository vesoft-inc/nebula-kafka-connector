package com.vesoft.nebula.client.graph.net;

import com.vesoft.nebula.client.graph.data.HostAddress;
import java.util.List;

public interface LoadBalancer {
    HostAddress getAddress();

    void close();

    void updateServersStatus();

    boolean isServersOK();

    List<HostAddress> getGoodAddresses();
}
