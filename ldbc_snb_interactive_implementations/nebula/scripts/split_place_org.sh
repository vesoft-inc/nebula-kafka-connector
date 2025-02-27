#!/bin/bash

set -eu
# Directory of this script
SCRIPT_DIR=$(dirname $(readlink -f "$0"))
# Directory of this project
PROJECT_DIR=$(dirname ${SCRIPT_DIR})

source ${SCRIPT_DIR}/env.sh

echo "split data"

# Split node data
cd ${LDBC_HOME}/static
# Split place to city, country, continent
head -1 place.csv > city.csv
head -1 place.csv > country.csv
head -1 place.csv > continent.csv

grep -e '|city$' place.csv >> city.csv
grep -e '|country$' place.csv >> country.csv
grep -e '|continent$' place.csv >> continent.csv

# Split organisation to university, company
head -1 organisation.csv > university.csv
head -1 organisation.csv > company.csv
grep -e '^[0-9]\+|university' organisation.csv >> university.csv
grep -e '^[0-9]\+|company' organisation.csv >> company.csv

# Split place_isPartOf_place.csv to city_isPartOf_country.csv、country_isPartOf_continent.csv
head -1 place_isPartOf_place.csv > city_isPartOf_country.csv
head -1 place_isPartOf_place.csv > country_isPartOf_continent.csv
file_name=("city.csv,city_isPartOf_country.csv" "country.csv,country_isPartOf_continent.csv")
for file in ${file_name[@]};do
        array=(`echo $file | tr ',' ' '`)
        for line in `awk -F '|' '{print $1}' ${array[0]}`;do
               echo `grep "^$line|" place_isPartOf_place.csv` >> ${array[1]}
        done
done

# Split organisation_isLocatedIn_place.csv to university_isLocatedIn_city.csv、company_isLocatedIn_country.csv
head -1 organisation_isLocatedIn_place.csv > university_isLocatedIn_city.csv
head -1 organisation_isLocatedIn_place.csv > company_isLocatedIn_country.csv
file_name=("university.csv,university_isLocatedIn_city.csv" "company.csv,company_isLocatedIn_country.csv")
for file in ${file_name[@]};do
        array=(`echo $file | tr ',' ' '`)
        for line in `awk -F '|' '{print $1}' ${array[0]}`;do
               echo `grep "^$line|" organisation_isLocatedIn_place.csv` >> ${array[1]}
        done
done

echo "Finish"
