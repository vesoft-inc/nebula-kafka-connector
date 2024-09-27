/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.driver.graph.decode.struct;

import com.google.protobuf.ByteString;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;

/**
 * PathHeader is stored in Vector data, path header is 32 bytes.
 * headNodeId: 8 bytes
 * tailNodeId: 8 bytes
 * size: 4 bytes
 * length: 4 bytes
 * headOffset: 4 bytes
 * tailOffset: 4 bytes
 */
public class PathHeader {
    private long headNodeId;
    private long tailNodeId;
    private int  size; // numNodes + numEdges
    private int  length; // numEdges

    // head node offset in the vector
    private int headOffset;
    // tail node offset in the vector
    private int tailOffset;

    public PathHeader(ByteString byteString, ByteOrder order) {
        ByteBuffer buffer = ByteBuffer
                .wrap(byteString.toByteArray())
                .order(order);
        this.headNodeId = buffer.getLong();
        this.tailNodeId = buffer.getLong();
        this.size = buffer.getInt();
        this.length = buffer.getInt();
        this.headOffset = buffer.getInt();
        this.tailOffset = buffer.getInt();
    }

    public long getHeadNodeId() {
        return headNodeId;
    }

    public long getTailNodeId() {
        return tailNodeId;
    }

    public int getSize() {
        return size;
    }

    public int getLength() {
        return length;
    }

    public int getHeadOffset() {
        return headOffset;
    }

    public int getTailOffset() {
        return tailOffset;
    }
}
