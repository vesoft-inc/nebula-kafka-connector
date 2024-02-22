/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package org.ldbcouncil.snb.impls.workloads.nebula;

import com.vesoft.nebula.client.graph.net.NebulaClient;

public class NebulaNewClient {
    private NebulaClient client;
    private String addr;

    public NebulaNewClient(NebulaClient client, String  addr){
        this.client = client;
        this.addr = addr;
    }

    public NebulaClient getClient() {
        return client;
    }

    public void setClient(NebulaClient client) {
        this.client = client;
    }

    public String getAddr() {
        return addr;
    }

    public void setAddr(String addr) {
        this.addr = addr;
    }
}
