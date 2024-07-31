#!/usr/bin/env bash

# Always use bash
shell=$(basename $(readlink /proc/$$/exe))
if [ ! x$shell = x"bash" ]
then
    bash $0 $@
    exit $?
fi

set -eu

function usage {
    cat <<EOF
Description:
  Build packages for NebulaGraph Tools
  This script is assumed to be placed at the 'package' directory of the source tree.
  The result packages will be put in ${PWD}/packages.

Usage: package.sh [--option]... [--option=argument]...
  --help|-h                     Print this message
  --product=NAME                     Product in nebula-ng-tools, e.g. golang.
  --package-name=NAME           Package name, e.g. nebula-graph
  --release-name=NAME           Release name, e.g. 5.0.0
  --exclude=DIR1,DIR2           Exclude directories, e.g. build,bin

Example:
  ./package.sh --product=golang --package-name=nebula-golang --release-name=5.0.0 --exclude=e2e   
	output: nebula-golang-5.0.0.tar.gz
EOF
}

# <cmd>...
function require_cmd {
    for cmd in $@
    do
        hash ${cmd} &>/dev/null || { echo "${cmd}: command not found" 1>&2; exit 1; }
    done
}

require_cmd tar

opt_product=""
opt_package_name=""
opt_release_name=""
opt_exclude=""
output_dir=packages


for opt in $@
do
    case ${opt} in
        --product=*)
						opt_product=$(echo "${opt}" | sed -E 's/[^=]+=?//')
            ;;
				--package-name=*)
						opt_package_name=$(echo "${opt}" | sed -E 's/[^=]+=?//')
						;;
				--release-name=*)
						opt_release_name=$(echo "${opt}" | sed -E 's/[^=]+=?//')
						;;
				--exclude=*)
						opt_exclude=$(echo "${opt}" | sed -E 's/[^=]+=?//')
						;;
        --help|-h)
            usage
            exit 0
            ;;
        *)
            echo "${opt}: Unknown option" 1>&2 && exit 1
            ;;
    esac
done

function check_option {
		if [ -z "${opt_product}" ] || [ -z "${opt_package_name}" ] || [ -z "${opt_release_name}" ]
		then
				echo "product, package-name, release-name are required." 1>&2
				usage
				exit 1
		fi
}

check_option

function package {
		local product=${opt_product}
		local package_name=${opt_package_name}
		local release_name=${opt_release_name}
		local exclude=${opt_exclude}

		local package_dir=${output_dir}/${product}
		local package_file=${package_name}-${release_name}.tar.gz

		if [ -d ${package_dir} ]
		then
				rm -rf ${package_dir}
		fi

		mkdir -p ${package_dir}

		cp -r ../${product} ${output_dir}
		rm -rf ${package_dir}/.git
		rm -rf ${package_dir}/.github
		rm -rf ${package_dir}/.gitignore

		if [ -n "${exclude}" ]
		then
				for dir in $(echo ${exclude} | tr ',' ' ')
				do
						rm -rf ${package_dir}/${dir}
				done
		fi
		tar -czvf ${output_dir}/${package_file} -C ${output_dir} ${product} > /dev/null\
				&& echo "Package ${package_file} is created successfully." \
				|| echo "Failed to create package ${package_file}."
		rm -rf ${package_dir} 
}

package
