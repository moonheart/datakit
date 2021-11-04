# DataKit install script for UNIX-like OS
# Wed Aug 11 11:35:28 CST 2021
# Author: tanb@jiagouyun.com

# https://stackoverflow.com/questions/19339248/append-line-to-etc-hosts-file-with-shell-script/37824076
# usage: updateHosts ip domain1 domain2 domain3 ...
function updateHosts() {
    for n in $@
    do
        if [ "$n" != "$1" ]; then
            # echo $n
            ip_address=$1
            host_name=$n
            # find existing instances in the host file and save the line numbers
            matches_in_hosts="$(grep -n $host_name /etc/hosts | cut -f1 -d:)"
            host_entry="${ip_address} ${host_name}"

            if [ ! -z "$matches_in_hosts" ]
            then
                # iterate over the line numbers on which matches were found
                while read -r line_number; do
                    # replace the text of each line with the desired host entry
                    if [[ "$OSTYPE" == "darwin"* ]]; then
                        sudo sed -i '' "${line_number}s/.*/${host_entry} /" /etc/hosts
                    else
                        sudo sed -i "${line_number}s/.*/${host_entry} /" /etc/hosts
                    fi
                done <<< "$matches_in_hosts"
            else
                echo "$host_entry" | sudo tee -a /etc/hosts > /dev/null
            fi
        fi
    done
}

set -e

domain=(
    "static.guance.com"
    "openway.guance.com"
    "dflux-dial.guance.com"

    "static.dataflux.cn"
    "openway.dataflux.cn"
    "dflux-dial.dataflux.cn"

    "zhuyun-static-files-production.oss-cn-hangzhou.aliyuncs.com"
)

# detect root user
if [ "$(echo "UID")" = "0" ]; then
	sudo_cmd=''
else
	sudo_cmd='sudo'
fi

##################
# Global variables
##################
RED="\033[31m"
CLR="\033[0m"
BLU="\033[34m"

##################
# Set Variables
##################

# Detect OS/Arch

arch=
case $(uname -m) in

	"x86_64")
		arch="amd64"
		;;

	"i386" | "i686")
		arch="386"
		;;

	"aarch64")
		arch="arm64"
		;;

	"arm" | "armv7l")
		arch="arm"
		;;

	*)
		# shellcheck disable=SC2059
		printf "${RED}[E] Unknown arch $(uname -m) ${CLR}\n"
		exit 1
		;;
esac

os=
if [[ "$OSTYPE" == "darwin"* ]]; then
	if [[ $arch != "amd64" ]] && [[ $arch != "arm64" ]]; then # Darwin only support amd64 and arm64, for arm64, use amd64 instead
		# shellcheck disable=SC2059
		printf "${RED}[E] Darwin only support amd64/arm64.${CLR}\n"
		exit 1;
	fi

	os="darwin"
else
	os="linux"
fi

# Select installer
installer_base_url="https://static.guance.com/datakit"
if [ -n "$DK_INSTALLER_BASE_URL" ]; then
	installer_base_url=$DK_INSTALLER_BASE_URL
fi

installer_file="installer-${os}-${arch}"
# shellcheck disable=SC2059
printf "${BLU} Detect installer ${installer_file}${CLR}\n"

installer_url="${installer_base_url}/${installer_file}"
installer=/tmp/dk-installer

dataway=
if [ -n "$DK_DATAWAY" ]; then
	dataway=$DK_DATAWAY
fi

upgrade=
if [ -n "$DK_UPGRADE" ]; then
	upgrade=$DK_UPGRADE
fi

if [ ! "$dataway" ]; then # check dataway on new install
	if [ ! "$upgrade" ]; then
		# shellcheck disable=SC2059
		printf "${RED}[E] DataWay not set in DK_DATAWAY.${CLR}\n"
		exit 1;
	fi
fi

def_inputs=
if [ -n "$DK_DEF_INPUTS" ]; then
	# shellcheck disable=SC2034
	def_inputs=$DK_DEF_INPUTS
fi

global_tags=
if [ -n "$DK_GLOBAL_TAGS" ]; then
	global_tags=$DK_GLOBAL_TAGS
fi

cloud_provider=
if [ -n "$DK_CLOUD_PROVIDER" ]; then
	cloud_provider=$DK_CLOUD_PROVIDER
fi

namespace=
if [ -n "$DK_NAMESPACE" ]; then
	namespace=$DK_NAMESPACE
fi

http_listen="localhost"
if [ -n "$DK_HTTP_LISTEN" ]; then
	http_listen=$DK_HTTP_LISTEN
fi

http_port=9529
if [ -n "$DK_HTTP_PORT" ]; then
	http_port=$DK_HTTP_PORT
fi

install_only=
if [ -n "$DK_INSTALL_ONLY" ]; then
	install_only=$DK_INSTALL_ONLY
fi

dca_white_list=""
if [ -n "$DK_DCA_WHITE_LIST" ]; then
	dca_white_list=$DK_DCA_WHITE_LIST
fi

dca_listen=""
if [ -n "$DK_DCA_LISTEN" ]; then
	dca_listen=$DK_DCA_LISTEN
fi

dca_enable=""
if [ -n "$DK_DCA_ENABLE" ]; then
	dca_enable=$DK_DCA_ENABLE
	if [ -z "$dca_white_list" ]; then
		printf "${RED}[E] DCA service is enabled, but white list is not set in DK_DCA_WHITE_LIST!${CLR}\n"
		exit 1;
	fi
fi

if [ -n "$HTTP_PROXY" ]; then
	proxy=$HTTP_PROXY
fi

if [ -n "$HTTPS_PROXY" ]; then
	proxy=$HTTPS_PROXY
fi

# check nginx proxy
proxy_type=""
if [ -n "$DK_PROXY_TYPE" ]; then
	proxy_type=$DK_PROXY_TYPE
	proxy_type=$(echo $proxy_type | tr '[:upper:]' '[:lower:]') # to lowercase
	printf "${BLU}\n* found Proxy Type: $proxy_type${CLR}\n"

	if [ "$proxy_type" == "nginx" ]; then
		# env DK_NGINX_IP has highest priority on proxy level
		if [ -n "$DK_NGINX_IP" ]; then
		    proxy=$DK_NGINX_IP
		    if [ "$proxy" != "" ]; then
			    printf "${BLU}\n* got nginx Proxy: $proxy${CLR}\n"

				for i in ${domain[@]}; do
				    updateHosts "$proxy" "$i"
                done
			fi
			proxy=""
		fi
	fi
fi

env_hostname=
if [ -n "$DK_HOSTNAME" ]; then
  env_hostname=$DK_HOSTNAME
fi

install_log=/var/log/datakit/install.log
if [ -n "$DK_INSTALL_LOG" ]; then
	install_log=$DK_INSTALL_LOG
fi

##################
# Try install...
##################
# shellcheck disable=SC2059
printf "${BLU}\n* Downloading installer ${installer}\n${CLR}"

rm -rf $installer

if [ "$proxy" ]; then # add proxy for curl
	# shellcheck disable=SC2086
	curl -x "$proxy" --fail --progress-bar $installer_url > $installer
else
	# shellcheck disable=SC2086
	curl --fail --progress-bar $installer_url > $installer
fi

# Set executable
chmod +x $installer

if [ "$upgrade" ]; then
	# shellcheck disable=SC2059
	printf "${BLU}\n* Upgrading DataKit...${CLR}\n"
    $sudo_cmd $installer --upgrade --proxy="${proxy}" | $sudo_cmd tee ${install_log}
else
	printf "${BLU}\n* Installing DataKit...${CLR}\n"
	if [ "$install_only" ]; then
		$sudo_cmd $installer                   \
			--dataway="${dataway}"               \
			--global-tags="${global_tags}"       \
			--cloud-provider="${cloud_provider}" \
			--namespace="${namespace}"           \
			--listen="${http_listen}"            \
			--port="${http_port}"                \
			--proxy="${proxy}"                   \
			--env_hostname="${env_hostname}"      \
			--dca-enable="${dca_enable}"				 \
			--dca-listen="${dca_listen}"				 \
			--dca-white-list="${dca_white_list}" \
			--install_only | $sudo_cmd tee ${install_log}
	else
		$sudo_cmd $installer                   \
		  --dataway="${dataway}"               \
			--global-tags="${global_tags}"       \
			--cloud-provider="${cloud_provider}" \
			--namespace="${namespace}"           \
			--listen="${http_listen}"            \
			--port="${http_port}"                \
			--env_hostname="${env_hostname}"      \
			--dca-enable="${dca_enable}"				 \
			--dca-listen="${dca_listen}"				 \
			--dca-white-list="${dca_white_list}"	\
			--proxy="${proxy}" | $sudo_cmd tee ${install_log}
	fi
fi
rm -rf $installer
