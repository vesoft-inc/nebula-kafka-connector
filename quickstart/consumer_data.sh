#!/bin/sh
# please make sure you already installed NebulaGraph, and update the user and password to connect-nebula-sink_actor.properties.
# create graph type and graph in NebulaGraph with qgl in nebula.gql.

# sink kafka data info NebulaGraph
${KAFKA_HOME}/bin/connect-standalone.sh standalone.properties connect-nebula-sink_actor.properties
${KAFKA_HOME}/bin/connect-standalone.sh standalone.properties connect-nebula-sink_director.properties
${KAFKA_HOME}/bin/connect-standalone.sh standalone.properties connect-nebula-sink_user.properties
${KAFKA_HOME}/bin/connect-standalone.sh standalone.properties connect-nebula-sink_movie.properties
${KAFKA_HOME}/bin/connect-standalone.sh standalone.properties connect-nebula-sink_genre.properties
${KAFKA_HOME}/bin/connect-standalone.sh standalone.properties connect-nebula-sink_act.properties
${KAFKA_HOME}/bin/connect-standalone.sh standalone.properties connect-nebula-sink_direct.properties
${KAFKA_HOME}/bin/connect-standalone.sh standalone.properties connect-nebula-sink_watched.properties
${KAFKA_HOME}/bin/connect-standalone.sh standalone.properties connect-nebula-sink_withgenre.properties

