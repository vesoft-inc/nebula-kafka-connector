package com.vesoft.nebula.driver.graph.net;

import com.vesoft.nebula.driver.graph.data.HostAddress;
import com.vesoft.nebula.driver.graph.exception.AuthFailedException;
import java.util.List;

public interface LoadBalancer {
    HostAddress getAddress();

    void close();

    void updateServersStatus() throws AuthFailedException;

    boolean isServersOK() throws AuthFailedException;

    List<HostAddress> getGoodAddresses();
}
