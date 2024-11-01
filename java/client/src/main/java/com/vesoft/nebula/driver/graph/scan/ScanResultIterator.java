package com.vesoft.nebula.driver.graph.scan;

import com.vesoft.nebula.driver.graph.data.ExtraInfo;
import com.vesoft.nebula.driver.graph.data.HostAddress;
import com.vesoft.nebula.driver.graph.data.ResultSet;
import com.vesoft.nebula.driver.graph.net.NebulaClient;
import java.io.Serializable;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.stream.Collectors;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class ScanResultIterator implements Serializable {
    private static final Logger logger = LoggerFactory.getLogger(ScanResultIterator.class);

    String serversAddress;
    String userName;
    Map<String, Object> authOptions;
    long requestTimeout;

    protected boolean hasNext = true;

    protected final Map<Integer, String> partCursor = new HashMap<>();

    protected final String graphName;
    protected final String labelName;
    protected List<String> propNames;
    protected int batchSize;
    protected final ExecutorService threadPool;
    protected final ConcurrentHashMap<Integer, NebulaClient> partClient = new ConcurrentHashMap<>();

    protected ScanResultIterator(String graphName,
                                 String labelName,
                                 List<String> propNames,
                                 List<Integer> parts,
                                 int batchSize,
                                 int parallel,
                                 List<HostAddress> servers,
                                 String userName,
                                 Map<String, Object> authOptions,
                                 long requestTimeout) {
        this.graphName = graphName;
        this.labelName = labelName;
        this.propNames = propNames;
        this.batchSize = batchSize;
        this.threadPool = Executors.newFixedThreadPool(parallel);
        this.serversAddress = servers
                .stream()
                .map(HostAddress::toString)
                .collect(Collectors.joining(","));
        for (int part : parts) {
            partCursor.put(part, "");
            try {
                partClient.put(part, NebulaClient.builder(serversAddress, userName)
                        .withAuthOptions(authOptions)
                        .build());
            } catch (Exception e) {
                throw new RuntimeException(e);
            }
        }
        this.userName = userName;
        this.authOptions = authOptions;
        this.requestTimeout = requestTimeout;
    }

    /**
     * if iter has more data
     *
     * @return true if the scan cursor is not at end.
     */
    public boolean hasNext() {
        return hasNext;
    }

    protected String getPropertyListString() {
        StringBuilder properties = new StringBuilder();
        String propertyListPrefix = "list[";
        properties.append(propertyListPrefix);
        for (String column : propNames) {
            properties.append("\"");
            properties.append(column);
            properties.append("\"");
            properties.append(",");
        }
        if (properties.length() > propertyListPrefix.length()) {
            properties.deleteCharAt(properties.length() - 1);
        }
        String propertyListSuffix = "]";
        properties.append(propertyListSuffix);
        return properties.toString();
    }

    protected ResultSet scan(String scanTemplate, Map.Entry<Integer, String> partCur)
            throws Exception {
        // construct the scan producer
        String producer = String.format(scanTemplate, graphName, graphName, labelName,
                getPropertyListString(), partCur.getKey(), partCur.getValue(), batchSize);
        NebulaClient client = partClient.get(partCur.getKey());
        return  client.execute(producer);
    }

    protected String getCursor(ResultSet resultSet) {
        ExtraInfo extraInfo = resultSet.getExtraInfo();
        if (extraInfo.getCursor() == null) {
            throw new RuntimeException("result does not contain cursor in extra info.");
        }
        return extraInfo.getCursor();
    }

    protected void close() {
        for (NebulaClient client: partClient.values()) {
            client.close();
        }
    }
}
