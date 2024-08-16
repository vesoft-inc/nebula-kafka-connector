package com.vesoft.nebula.driver.graph.net;

import com.vesoft.nebula.driver.graph.data.HostAddress;
import java.util.List;

public interface LoadBalancer {
    HostAddress getAddress();

    void close();

    void updateServersStatus();

    boolean isServersOK();

    List<HostAddress> getGoodAddresses();
}
