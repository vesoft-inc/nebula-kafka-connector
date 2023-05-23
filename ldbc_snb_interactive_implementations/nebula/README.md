
# LDBC SNB Interactive Nebula5.0 implementation

## Environmental requirements
1. Since the nebula 5.0 kernel and surrounding tools have not been released, so temporary java client alternative <br>
   ``git clone  https://github.com/vesoft-inc/nebula-ng-tools.git
     cd nebula-ng-tools/java
     mvn clean install -Dmaven.javadoc.skip=true -Dmaven.test.skip=true -Dgpg.skip``
2. Java environment requirements<br>
    ``**requried java11** Download Address https://www.oracle.com/java/technologies/javase/jdk11-archive-downloads.html
      Install JDK. under the ubuntu system. sudo dpkg -i jdk-11.0.13_linux-x64_bin.deb
      Configure environment variables
      export JAVA_HOME=/usr/lib/jvm/jdk-11.0.13
      export PATH=$PATH:$JAVA_HOME/bin
      Check whether the JDK is configured successfully: Enter java -version in the terminal, if the output version number is the installed JDK version number, it means the installation is successful``
3. LDBC driver (Tag v1.2.0)<br>
    `` git clone --b v1.2.0 https://github.com/ldbc/ldbc_snb_interactive_driver
       cd ldbc_snb_driver
       mvn clean package -DskipTests (Need to modify java version in pom.xml to 11)``
4. execute scripts/build.sh  to build nebula implementation's jar package<br>

~some problems encountered when generating cypher data:~
   ~pip install backports.zoneinfo (ldbc_snb_driver/scripts/create_update_stream.py and ldbc_snb_driver/paramgen/paramgen.py  from zoneinfo import ZoneInfo -> from backports.zoneinfo import ZoneInfo)~
   ~pip install duckdb~
   ~pip install networkit~

## Load data into nebula graph
1. by default, load the data in `127.0.0.1:3713` with graphName `sf01`.<br>
   `execute ./nebula/scripts/load.sh`

2. it is recommended to configure env variable by yourself.<br>
-SCALE_FACTOR=0.1
-NEBULA_ADDRESS=192.168.15.8:9669
-GRAPH_NAME=sf01
`then execute ./nebula/scripts/load.sh`

## HOW TO VALIDATE
1. the verification file is stored in the driver/sf0_1_validation_param.csv  which contains the parameters and corresponding results of all queries(IC、IS、update)
2. each line in the verification file represents the input parameters and corresponding results of a query
3. If you do not verify a certain query, you can delete the input parameters and results of the statement from the file

| query | prefix                 |
|-------|------------------------|
| IC1   | personIdQ1"            |
| IC2   | personIdQ2             | 
| IC3   | personIdQ3             | 
| IC4   | personIdQ4             | 
| IC5   | personIdQ5             | 
| IC6   | personIdQ6             | 
| IC7   | personIdQ7             | 
| IC8   | personIdQ8             | 
| IC9   | personIdQ9             | 
| IC10  | personIdQ10            | 
| IC11  | personIdQ11            | 
| IC12  | personIdQ12            | 
| IC13  | person1IdQ13StartNode  | 
| IC14  | person1IdQ14StartNode  | 
| IS1   | personIdSQ1            | 
| IS2   | personIdSQ2            | 
| IS3   | personIdSQ3            | 

**if you want to turn off IC14, execute  sed -i '/person1IdQ14StartNode/d' sf0_1_validation_param.csv**
*the update statement cannot be turned off, and the four query statements of IS4 IS5 IS6 IS7 have not found a way to turn off*

4. Modify the option of `(CONFIG)` in the configuration file validate.properties, then execute validate.sh
5. The result of the verification will be saved in the dataset directory
6. Since the verification contains the update statement, which will change the state of the database, every time the verification is performed.
   the database needs to be cleared, and then the data is re-imported before verification

---
10. For each implementation, it is possible to (1) create validation parameters, (2) validate against an existing validation parameters, and (3) run the benchmark. Set the parameters according to your system configuration in the appropriate `.properties` file and run the driver with one of the following scripts:

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