#!/bin/sh

# create topic in kafka
${KAFKA_HOME}/bin/kafka-topics.sh --create --topic actor --broker-list "127.0.0.1:9092" --replication-factor 1 --partitions 1
${KAFKA_HOME}/bin/kafka-topics.sh --create --topic director --broker-list "127.0.0.1:9092" --replication-factor 1 --partitions 1
${KAFKA_HOME}/bin/kafka-topics.sh --create --topic genre --broker-list "127.0.0.1:9092" --replication-factor 1 --partitions 1
${KAFKA_HOME}/bin/kafka-topics.sh --create --topic movie --broker-list "127.0.0.1:9092" --replication-factor 1 --partitions 1
${KAFKA_HOME}/bin/kafka-topics.sh --create --topic user --broker-list "127.0.0.1:9092" --replication-factor 1 --partitions 1
${KAFKA_HOME}/bin/kafka-topics.sh --create --topic act --broker-list "127.0.0.1:9092" --replication-factor 1 --partitions 1
${KAFKA_HOME}/bin/kafka-topics.sh --create --topic direct --broker-list "127.0.0.1:9092" --replication-factor 1 --partitions 1
${KAFKA_HOME}/bin/kafka-topics.sh --create --topic withGenre --broker-list "127.0.0.1:9092" --replication-factor 1 --partitions 1
${KAFKA_HOME}/bin/kafka-topics.sh --create --topic watched --broker-list "127.0.0.1:9092" --replication-factor 1 --partitions 1

# producer the test.json into kafka
cat movie_data/actor.json | ${KAFKA_HOME}/bin/kafka-console-producer.sh --broker-list "127.0.0.1:9092" --topic actor
cat movie_data/director.json | ${KAFKA_HOME}/bin/kafka-console-producer.sh --broker-list "127.0.0.1:9092" --topic director
cat movie_data/genre.json | ${KAFKA_HOME}/bin/kafka-console-producer.sh --broker-list "127.0.0.1:9092" --topic genre
cat movie_data/movie.json | ${KAFKA_HOME}/bin/kafka-console-producer.sh --broker-list "127.0.0.1:9092" --topic movie
cat movie_data/user.json | ${KAFKA_HOME}/bin/kafka-console-producer.sh --broker-list "127.0.0.1:9092" --topic user
cat movie_data/actor_act_movie.json | ${KAFKA_HOME}/bin/kafka-console-producer.sh --broker-list "127.0.0.1:9092" --topic act
cat movie_data/director_direct_movie.json | ${KAFKA_HOME}/bin/kafka-console-producer.sh --broker-list "127.0.0.1:9092" --topic direct
cat movie_data/movie_withgenre_genre.json | ${KAFKA_HOME}/bin/kafka-console-producer.sh --broker-list "127.0.0.1:9092" --topic withgenre
cat movie_data/user_watched_movies.json | ${KAFKA_HOME}/bin/kafka-console-producer.sh --broker-list "127.0.0.1:9092" --topic watched

