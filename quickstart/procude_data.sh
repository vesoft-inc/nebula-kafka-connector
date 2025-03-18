#!/bin/sh

# create topic in kafka
${KAFKA_HOME}/bin/kafka-topics.sh --create --topic connect-test --broker-list "127.0.0.1:9092" --replication-factor 1 --partitions 1

# producer the test.json into kafka
cat test.json | ${KAFKA_HOME}/bin/kafka-console-producer.sh --broker-list "127.0.0.1:9092" --topic connect-test

