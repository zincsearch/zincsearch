#!/usr/bin/env bash
#
# Installs ZincSearch the quick way, for adventurers that want to spend more time
# grooming their cats.
#
# Requires curl, grep, cut, tar, uname, chmod, mv, rm.

[[ $- = *i* ]] && echo "Don't source this script!" && return 10

header() {
    cat 1>&2 <<EOF
ZincSearch Installer

Website: https://zinc.dev
Docs: https://docs.zinc.dev
Repo: https://github.com/zinclabs/zincsearch

EOF
}

check_cmd() {
	command -v "$1" > /dev/null 2>&1
}

check_tools() {
	Tools=("curl" "grep" "cut" "tar" "uname" "chmod" "mv" "rm")

	for tool in ${Tools[*]}; do
		if ! check_cmd $tool; then
			echo "Aborted, missing $tool, sorry!"
			exit 6
		fi
	done
}

install_zincsearch()
{
	trap 'echo -e "Aborted, error $? in command: $BASH_COMMAND"; trap ERR; exit 1' ERR
	install_path="/usr/local/bin"
	zinc_os="unsupported"
	zinc_arch="unknown"
	zinc_arm=""

	header
	check_tools

	if [[ -n "$PREFIX" ]]; then
		install_path="$PREFIX/bin"
	fi

	# Fall back to /usr/bin if necessary
	if [[ ! -d $install_path ]]; then
		install_path="/usr/bin"
	fi

	# Not every platform has or needs sudo (https://termux.com/linux.html)
	((EUID)) && sudo_cmd="sudo"

	#########################
	# Which OS and version? #
	#########################

	zinc_bin="zincsearch"
	zinc_dl_ext=".tar.gz"

	# NOTE: `uname -m` is more accurate and universal than `arch`
	# See https://en.wikipedia.org/wiki/Uname
	unamem="$(uname -m)"
	if [[ $unamem == *aarch64* ]]; then
		zinc_arch="arm64"
    elif [[ $unamem == *arm64* ]]; then
		zinc_arch="arm64"
	elif [[ $unamem == *64* ]]; then
		zinc_arch="x86_64"
	elif [[ $unamem == *armv6l* ]]; then
		zinc_arch="arm"
		zinc_arm="v6"
	elif [[ $unamem == *armv7l* ]]; then
		zinc_arch="arm"
		zinc_arm="v7"
	else
		echo "Aborted, unsupported or unknown architecture: $unamem"
		return 2
	fi

	unameu="$(tr '[:lower:]' '[:upper:]' <<<$(uname))"
	if [[ $unameu == *DARWIN* ]]; then
		zinc_os="Darwin"
		version=${vers##*ProductVersion:}
	elif [[ $unameu == *LINUX* ]]; then
		zinc_os="Linux"
	elif [[ $unameu == *FREEBSD* ]]; then
		zinc_os="freebsd"
	elif [[ $unameu == *OPENBSD* ]]; then
		zinc_os="openbsd"
	elif [[ $unameu == *WIN* || $unameu == MSYS* ]]; then
		# Should catch cygwin
		sudo_cmd=""
		zinc_os="Windows"
		zinc_bin=$zinc_bin.exe
	else
		echo "Aborted, unsupported or unknown os: $uname"
		return 6
	fi

	########################
	# Download and extract #
	########################

	echo "Downloading ZincSearch for ${zinc_os}/${zinc_arch}${zinc_arm}..."
	zinc_file="zinc_${zinc_os}_${zinc_arch}${zinc_arm}${zinc_dl_ext}"

	if [[ "$#" -eq 0 ]]; then
		# get latest release
		zinc_tag=$(curl -s https://api.github.com/repos/zinclabs/zincsearch/releases/latest | grep 'tag_name' | cut -d\" -f4)
		zinc_version=$(echo ${zinc_tag} | cut -c2-)
	elif [[ "$#" -gt 1 ]]; then
		echo "Too many arguments."
		exit 1
	elif [ -n $1  ]; then
		# try to get passed version
		zinc_tag="v$1"
		zinc_version=$1
	fi

	zinc_url="https://github.com/zinclabs/zincsearch/releases/download/${zinc_tag}/zinc_${zinc_version}_${zinc_os}_${zinc_arch}${zinc_arm}.tar.gz"

	dl="/tmp/$zinc_file"
	rm -rf -- "$dl"

    echo "Downloading $zinc_url"

	curl -fsSL "$zinc_url" -o "$dl"

	echo "Extracting..."
	case "$zinc_file" in
		*.tar.gz) tar -xzf "$dl" -C "$PREFIX/tmp/" "$zinc_bin" ;;
	esac
	chmod +x "$PREFIX/tmp/$zinc_bin"

	echo "Putting ZincSearch in $install_path (may require password)"
	$sudo_cmd mv "$PREFIX/tmp/$zinc_bin" "$install_path/$zinc_bin"
	$sudo_cmd rm -- "$dl"

	# check installation
	# $zinc_bin -version

	echo "Successfully installed ZincSearch"
	trap ERR
	return 0
}

install_zincsearch $@
