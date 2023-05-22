
# LDBC SNB Interactive Nebula5.0 implementation
## User's guide

## Setup
* Bash
* Java 11
1、 Since the nebula 5.0 kernel and surrounding tools have not been released, so temporary java client alternative
      git clone -b  temp_java_client_using_old_thrift https://github.com/Nicole00/nebula-ng-tools.git
      cd nebula-ng-tools/java
      mvn clean install -Dmaven.javadoc.skip=true -Dmaven.test.skip=true -Dgpg.skip
2、 Java environment requirements
      Download Address https://www.oracle.com/java/technologies/javase/jdk11-archive-downloads.html
      Install JDK. under the ubuntu system. sudo dpkg -i jdk-11.0.13_linux-x64_bin.deb
      Configure environment variables
      export JAVA_HOME=/usr/lib/jvm/jdk-11.0.13
      export PATH=$PATH:$JAVA_HOME/bin
      Check whether the JDK is configured successfully: Enter java -version in the terminal, if the output version number is the installed JDK version number, it means the installation is successful
3、 LDBC driver
      DownLoad Address https://github.com/ldbc/ldbc_snb_interactive_driver
      cd ldbc_snb_driver
      mvn clean package -DskipTests (Need to modify java version in pom.xml to 11)

some problems encountered when generating cypher data:
   pip install backports.zoneinfo (ldbc_snb_driver/scripts/create_update_stream.py and ldbc_snb_driver/paramgen/paramgen.py  from zoneinfo import ZoneInfo -> from backports.zoneinfo import ZoneInfo)
   pip install duckdb
   pip install networkit


3. Load data into nebula graph

```
# by default, load the data in `127.0.0.1:9669` with space `sf0_1`.
./nebula/scripts/load.sh

# or pass the environments.

SCALE_FACTOR=10 \
NEBULA_ADDRESS=192.168.15.8:9669 \
NEBULA_SPACE=sf10 \
./nebula/scripts/load.sh

```

4. For each implementation, it is possible to (1) create validation parameters, (2) validate against an existing validation parameters, and (3) run the benchmark. Set the parameters according to your system configuration in the appropriate `.properties` file and run the driver with one of the following scripts:

   ```bash
   # Interactive workload - note that if the workload contains updates, the database needs to be re-loaded between steps
   # ./interactive-create-validation-parameters.sh
   # ./interactive-validate.sh
   ./interactive-benchmark.sh
   ```
   
   - configs
   ```bash
   endpoint=192.168.15.3:9669,192.168.15.5:9669,192.168.15.6:9669
   user=root
   password=nebula
   spaceName=sf100

   ## 路径都是基于执行路径，最好写绝对路径
   queryDir=../queries/
   ldbc.snb.interactive.parameters_dir=../test-data/substitution_parameters/
   ldbc.snb.interactive.updates_dir=../test-data/social_network/
   
   ldbc.snb.interactive.LdbcQuery10_enable=false
   ldbc.snb.interactive.LdbcQuery14_enable=false
   ```