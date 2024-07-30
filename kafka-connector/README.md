# kafka Connect Nebula

kafka-connect-nebula is a Kafka sink connector for capturing all row from kafka and streaming the
rows to NebulaGraph.

## Target node/edge types for Kafka topics

Kafka topics can be mapped to existing NebulaGraph node/edge types in the Kafka configuration.

The connector converts the row data with json string format, Struct format, Map format into a Node
record or a Edge record
according to the sink config file, and then write the node/edge record into NebulaGraph.

## Configs for Nebula kafka connect.

You need to specify configuration settings for your connector.
These can be found in the config/connect-nebula-sink.properties file.

|        config name         | required |                                                                             description                                                                              | default value |                   available values                   |
|:--------------------------:|:--------:|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------:|:-------------:|:----------------------------------------------------:|
|            name            |   true   |                                                                  A unique name for the con nector.                                                                   |       _       |                  kafka-nebula-sink                   |
|      connector.class       |   true   |                                                                 The entry class for this connector.                                                                  |       _       | com.vesoft.nebula.connector.sink.NebulaSinkConnector |
|         tasks.max          |   true   |     The maximum number of tasks that should be created for this connector. The connector may create fewer tasks if it cannot achieve this level of parallelism.      |       -       |                      Int number                      |
|           topics           |   true   |                                                            kafka topic need to sink data into NebulaGraph                                                            |       _       |                          _                           |
|    nebula.graph.servers    |   true   |                                            NebulaGraph graphd servers address, split multiple addresses by English comma                                             |       _       |                          _                           |
|        nebula.user         |   true   |                                                                           NebulaGraph user                                                                           |       _       |                          _                           |
|       nebula.passwd        |  false   |                                                                    NebulaGraph password for user                                                                     |       _       |                          _                           |
|     nebula.antuOptions     |  false   |                                                                    NebulaGraph other auth options                                                                    |       _       |                          _                           |
|   nebula.request.timeout   |  false   |                                                              The request timeout for NebulaGraph Client                                                              |   3000(ms)    |                          -                           |
|   nebula.sin.partitions    |  false   |                                                                  The partitions for sink connector                                                                   |      10       |                          -                           |
|     nebula.graph.name      |   true   |                                                               NebulaGraph graph name to sink data into                                                               |       _       |                          _                           |
|      nebula.data.type      |  false   | NebulaGraph data type to sink data into.NODE:sink kafka data into node type;EDGE:sink kafka data into edge type;BOTH:sink kafka data into both one node and one edge |     NODE      |                    NODE,EDGE,BOTH                    |
|    nebula.node.typeName    |  false   |                                        The node type name in NebulaGraph,it's required when nebula.data.type is NODE or BOTH                                         |       _       |                          _                           |
|    nebula.edge.typeName    |  false   |                                               The edge type name, it's required when nebula.data.type is EDGE or BOTH                                                |       -       |                          -                           |
|   nebula.node.primaryKey   |  false   |                           The kafka property name to sink the data as node primarykey, it's required when nebula.data.type is NODE or BOTH                           |       _       |                          _                           |
| nebula.node.property.names |  false   |                                      The node properties need to sink into, it's required when nebula.data.type is NODE or BOTH                                      |       -       |                          -                           |
| kafka.node.property.names  |  false   |                                The kafka properties need to write to NebulaGraph, it's required when nebula.data.type is NODE or BOTH                                |       -       |                          -                           |
|     nebula.edge.srcPk      |  false   |                     The kafka property name to sink the data as edge's src node primary key,it's required when nebula.data.type is EDGE or BOTH                      |       -       |                          -                           |
|     nebula.edge.dstPk      |  false   |                     The kafka property name to sink the data as edge's dst node primary key,it's required when nebula.data.type is EDGE or BOTH                      |       -       |                          -                           |
| nebula.edge.property.names |  false   |                                      The node properties need to sink into, it's required when nebula.data.type is NODE or BOTH                                      |       -       |                          -                           |
| kafka.edge.property.names  |  false   |                                The kafka properties need to write to NebulaGraph, it's required when nebula.data.type is EDGE or BOTH                                |       -       |                          -                           |
|      nebula.batchSize      |  false   |                                                            The batch size when sink data into NebulaGraph                                                            |     2000      |                          -                           |
|      nebula.sink.mode      |  false   |                                                                    The sink mode for NebulaGraph                                                                     |    INSERT     |                 INSERT,UPDATE,DELETE                 |
|       key.converter        |  false   |                                          Use this parameter to override the default key converter class set by the worker.                                           |       -       |   org.apache.kafka.connect.storage.StringConverter   |
|      value.converter       |  false   |                                         Use this parameter to override the default value converter class set by the worker.                                          |               |   org.apache.kafka.connect.storage.StringConverter   |

## Standalone Quickstart

> NOTE: You must have the Confluent Kafka Platform installed in order to run the example.

### 1. get the connector jar

```agsl
git clone https://github.com/vesoft-inc/nebula-ng-tools.git
cd nebula-ng-tools/kafka-connector/connector
mvn clean package -Dmaven.test.skip=true
```

after the command finished, there will be `kafka-connect-nebula-5.0-SNAPSHOT.jar` in
kafka-connector/connector/target.

### 2. Put the connector jar to Kafka lib

put the kafka-connect-nebula-5.0-SNAPSHOT.jar into your kafka env: KAFKA_HOME/lib

### 3. config the connect-nebula-sink.properties
you can get the demo config in quickstart/connect-nebula-sink.properties.

Please update the config value according to your Kafka environment and NebulaGraph environment.

### 4. producer data into Kafka
you can producer the json data with quickstart/produce_data.sh

### 5. consumer data from kafka into NebulaGraph

```agsl
${KAFKA_HOME}/bin/connect-standalone.sh ${KAFKA_HOME}/config/standalone.properties connect-nebula-sink.properties
```


