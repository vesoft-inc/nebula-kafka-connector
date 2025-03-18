#!/bin/sh
# please make sure you already installed NebulaGraph, and update the user and password to connect-nebula-sink.properties.
# create graph type and graph in NebulaGraph with qgl in nebula.gql.

# sink kafka data info NebulaGraph
${KAFKA_HOME}/bin/connect-standalone.sh standalone.properties connect-nebula-sink.properties

