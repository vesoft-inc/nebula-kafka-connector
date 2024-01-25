#NebulaGraph Java SDK for NebulaGraph 5.0

NebulaGraph Java SDK is a Java client for developers to connect their projects to Nebula Graph.

## How to package
* dependency: your environment must have JDK8.


```agsl
./mvnw clean package
```
the sdk jar will be generated in java/client/target/client-5.0.0.jar

## Example to use Java SDK
```agsl
       NebulaClient client = null;
        try {
            client = NebulaClient.builder(host, user, passwd)
                    .setConnectTimeoutMills(1000)
                    .setRequestTimeoutMills(3000)
                    .setMaxSessionSize(10)
                    .setMinSessionSize(1)
                    .setRetryTimes(3)
                    .setIntervalTimeMills(1000)
                    .setReconnect(true)
                    .setBlockWhenExhausted(true)
                    .setMaxWaitMills(1000)
                    .setStrictlyServerHealthy(true)
                    .build();
            client.execute("USE nba MATCH ()-[e:follow]->() RETURN e.followness, e.likeness");
        } catch (Exception e) {
            e.printStackTrace();
            System.exit(1);
        } finally {
            if (client != null) {
                client.close();
            }
        }
```
