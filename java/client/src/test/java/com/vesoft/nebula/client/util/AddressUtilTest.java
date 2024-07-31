package com.vesoft.nebula.client.util;

import com.vesoft.nebula.client.graph.data.HostAddress;
import com.vesoft.nebula.client.graph.utils.AddressUtil;
import java.net.UnknownHostException;
import java.util.List;
import org.junit.Test;

public class AddressUtilTest {


    @Test
    public void testIpv4Adderss() {
        try {
            // test ipv4
            List<HostAddress> addresses = AddressUtil.validateAddress("localhost:9669");
            assert (addresses.size() == 1);
            assert (addresses.get(0).toString().equals("localhost:9669"));

            addresses = AddressUtil.validateAddress(",127.0.0.1:9669 ");
            assert (addresses.size() == 1);
            assert (addresses.get(0).toString().equals("127.0.0.1:9669"));
        } catch (UnknownHostException e) {
            assert (false);
        }
    }

    @Test
    public void testIpv6Address() {
        try {
            // test ipv6
            List<HostAddress> addresses = AddressUtil
                    .validateAddress("fe80::64eb:eff:fe32:f7fe:9669");
            assert (addresses.size() == 1);
            assert (addresses.get(0).toString().equals("fe80::64eb:eff:fe32:f7fe:9669"));

            addresses = AddressUtil.validateAddress("[fe80::64eb:eff:fe32:f7fe]:9669");
            assert (addresses.get(0).toString().equals("[fe80::64eb:eff:fe32:f7fe]:9669"));
        } catch (UnknownHostException e) {
            assert (false);
        }
    }
}
