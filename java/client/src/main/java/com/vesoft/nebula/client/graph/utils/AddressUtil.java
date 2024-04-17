/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package com.vesoft.nebula.client.graph.utils;

import com.google.common.net.InetAddresses;
import com.google.common.net.InternetDomainName;
import com.vesoft.nebula.client.graph.data.HostAddress;
import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.List;

public class AddressUtil {

    /**
     * validate the graph addresses
     *
     * @param addresses graph server addresses, multiple addresses are split by comma
     * @return List of HostAddress
     * @throws IllegalArgumentException if address id not split by comma or port is beyond range
     * @throws UnknownHostException     if address host is wrong
     */
    public static List<HostAddress> validateAddress(String addresses) throws UnknownHostException {
        List<HostAddress> newAddrs = new ArrayList<>();
        for (String addr : addresses.split(",")) {
            if (addr.isEmpty()) {
                continue;
            }
            int indexOfLastComa = addr.lastIndexOf(":");
            String host = addr.substring(0, indexOfLastComa);
            int port = Integer.parseInt(addr.substring(indexOfLastComa + 1));

            // get all host name
            InetAddress[] inetAddresses = InetAddress.getAllByName(host);
            for (InetAddress inetAddress : inetAddresses) {
                String ip = inetAddress.getHostAddress();
                if (!(InetAddresses.isInetAddress(ip)
                        || InetAddresses.isUriInetAddress(ip)
                        || InternetDomainName.isValid(ip))
                        || (port <= 0 || port >= 65535)) {
                    throw new IllegalArgumentException(
                            String.format("host %s and port %d is invalid.", ip, port));
                }
            }
            newAddrs.add(new HostAddress(host, port));
        }
        return newAddrs;
    }
}
