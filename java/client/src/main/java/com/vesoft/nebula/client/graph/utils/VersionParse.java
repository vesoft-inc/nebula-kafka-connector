/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.utils;

import com.google.common.base.Charsets;
import com.google.protobuf.ByteString;
import com.google.protobuf.Descriptors;
import com.vesoft.nebula.proto.version.Version;
import java.util.Map;

public class VersionParse {
    public static ByteString getProtoVersion() {
        ByteString version = null;
        Map<Descriptors.FieldDescriptor, Object> descriptors =
                Version.getDescriptor().getOptions().getAllFields();
        for (Map.Entry<Descriptors.FieldDescriptor, Object> entry : descriptors.entrySet()) {
            if (entry.getKey().getFullName().equals("nebula.proto.version.protocol_version")) {
                version = (ByteString) entry.getValue();
            }
        }
        return version;
    }
}
