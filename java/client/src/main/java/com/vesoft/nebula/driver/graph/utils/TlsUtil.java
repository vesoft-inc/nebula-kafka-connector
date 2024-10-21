/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.driver.graph.utils;

import com.vesoft.nebula.driver.graph.exception.IOErrorException;
import io.grpc.netty.shaded.io.grpc.netty.GrpcSslContexts;
import io.grpc.netty.shaded.io.netty.handler.ssl.SslContext;
import io.grpc.netty.shaded.io.netty.handler.ssl.SslContextBuilder;
import java.io.File;
import javax.net.ssl.SSLException;

public class TlsUtil {

    public static SslContext getSslContext(String ca,
                                           String cert,
                                           String key) throws IOErrorException {
        File caFile = new File(ca);
        SslContextBuilder builder = GrpcSslContexts
                .forClient()
                .trustManager(caFile);
        if (cert != null && key != null) {
            builder.keyManager(new File(cert), new File(key));
        }
        try {
            return builder.build();
        } catch (SSLException e) {
            throw new IOErrorException(IOErrorException.E_SSL_ERROR, e.getMessage());
        }
    }
}
