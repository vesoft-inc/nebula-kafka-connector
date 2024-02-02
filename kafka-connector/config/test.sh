${KAFKA_HOME}/bin/kafka-topics.sh --create --topic connect-test --broker-list "192.168.8.171:9092" --replication-factor 1 --partitions 1
cat test.json | ${KAFKA_HOME}/bin/kafka-console-producer.sh --broker-list "192.168.8.171:9092" --topic connect-test
${KAFKA_HOME}/bin/connect-standalone.sh ../connector/src/main/resources/standalone.properties connect-nebula-sink.properties
