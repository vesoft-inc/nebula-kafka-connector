/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package org.ldbcouncil.snb.impls.workloads.nebula;

import com.vesoft.nebula.driver.graph.net.NebulaPool;

public class NewNebulaPool {
    private NebulaPool pool;
    private String addr;

    public NewNebulaPool(NebulaPool pool, String  addr){
        this.pool = pool;
        this.addr = addr;
    }

    public NebulaPool getPool() {
        return pool;
    }

    public void setClient(NebulaPool pool) {
        this.pool = pool;
    }

    public String getAddr() {
        return addr;
    }

    public void setAddr(String addr) {
        this.addr = addr;
    }
}
