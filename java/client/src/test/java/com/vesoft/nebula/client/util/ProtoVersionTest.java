/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.util;

import com.google.common.base.Charsets;
import com.google.protobuf.ByteString;
import com.vesoft.nebula.client.graph.utils.VersionParse;
import org.junit.Assert;
import org.junit.Test;

public class ProtoVersionTest {

    @Test
    public void testProtoVersion() {
        Assert.assertEquals(ByteString.copyFrom("5.0.0", Charsets.UTF_8),
                VersionParse.getProtoVersion());
    }
}
