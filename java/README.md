#NebulaGraph Java SDK for NebulaGraph 5.0

NebulaGraph Java SDK is a Java client for developers to connect their projects to Nebula Graph.

## How to package
* dependency: your environment must have JDK8.


```agsl
./mvnw clean package -Dmaven.test.skip=true
```
the sdk jar will be generated in java/client/target/client-5.0.0.jar

## Example to use Java SDK
There are two ways to use Java SDK: get a NebulaClient from NebulaPool or get a NebulaClient by yourself.

* Using NebulaClient through the pool:
```agsl
       NebulaPool pool = null;
        try {
            pool = NebulaPool
                    .builder(addresses, userName, password)
                    .withMaxClientSize(10)
                    .withMinClientSize(1)
                    .withConnectTimeoutMills(1000)
                    .withRequestTimeoutMills(30000)
                    .withBlockWhenExhausted(true)
                    .withMaxWaitMills(Long.MAX_VALUE)
                    .build();
            NebulaClient client = pool.getClient();
            client.execute("USE nba MATCH (v:player) RETURN v.id, v.name, v.score, v.gender, v.rate");
            pool.returnClient(client);
        } catch (Exception e) {
            throw e;
        } finally {
            if (pool != null) {
                pool.close();
            }
        }
```

* Using NebulaClient by yourself
```agsl
        NebulaClient client = null;
        try {
            client = NebulaClient.builder(host, user, passwd)
                    .withAuthOptions(Collections.emptyMap())
                    .withConnectTimeoutMills(1000)
                    .withRequestTimeoutMills(3000)
                    .build();
            client.execute("USE nba MATCH (v:player) RETURN v.id, v.name, v.score, v.gender, v.rate");
        } catch (Exception e) {
            throw e;
        } finally {
            if (client != null) {
                client.close();
            }
        }
```
