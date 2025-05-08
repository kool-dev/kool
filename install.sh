#!/usr/bin/env bash

set -euo pipefail

echo -e "Hello, we are gonna install the \033[33mlatest stable\033[39m version of Kool!"

DEFAULT_DOWNLOAD_URL="https://github.com/kool-dev/kool/releases/latest/download"
if [ -z "${DOWNLOAD_URL:-}" ]; then
	DOWNLOAD_URL=$DEFAULT_DOWNLOAD_URL
fi

DEFAULT_BIN="/usr/local/bin/kool"
if [ -z "${BIN_PATH:-}" ]; then
	BIN_PATH=$DEFAULT_BIN
fi

is_darwin() {
	case "$(uname -s)" in
	*darwin* ) true ;;
	*Darwin* ) true ;;
	* ) false;;
	esac
}

do_install () {
	ARCH=$(uname -m)
	PLAT="linux"

	if is_darwin; then
		PLAT="darwin"
	fi

	if [ "$ARCH" == "x86_64" ]; then
		ARCH="amd64"
	fi

	if [ "$ARCH" == "aarch64" ]; then
		ARCH="arm64"
	fi

	echo "Downloading latest binary (kool-$PLAT-$ARCH)..."

    rm -f /tmp/kool_binary
	
	# fallback to wget if no curl available
    if command -v curl &> /dev/null; then
        curl -fsSL "$DOWNLOAD_URL/kool-$PLAT-$ARCH" -o /tmp/kool_binary
    elif command -v wget &> /dev/null; then
        wget -qO /tmp/kool_binary "$DOWNLOAD_URL/kool-$PLAT-$ARCH"
    else
        echo -e "\033[31;31mError: Neither curl nor wget is available. Please install one of them to proceed.\033[0m"
        exit 1
    fi

	# check for running kool process which would prevent
	# replacing existing version under Linux.
	if [ command -v kool &> /dev/null ]; then
		if [ ! is_darwin ]; then
			running=$(ps aux | grep kool | grep -v grep | wc -l | awk '{ print $1 }')
			if [ "$running" != "0" ]; then
				echo -e "\033[31;31mThere is a kool process still running. You might need to stop them before we replace the current binary.\033[0m"
			fi
		fi
	fi

	echo -e "Moving kool binary to $BIN_PATH..."
	if [ -w $(dirname $BIN_PATH) ]; then
		mv -f /tmp/kool_binary $BIN_PATH
		chmod +x $BIN_PATH
	else
		echo "(requires sudo)"
		sudo mv -f /tmp/kool_binary $BIN_PATH
		sudo chmod +x $BIN_PATH
	fi

	start_success="\033[0;32m"
	end_success="\033[0m"
	start_error="\033[1;31m"
	end_error="\033[0m"

	if ! command -v docker &> /dev/null; then
		builtin echo -e "${start_error}We could not identify the Docker installed.${end_error}"
		builtin echo -e "Please refer to the official documentation to get it: https://docs.docker.com/get-docker/"
		exit
	fi

	composeVersion=$(docker compose version || true)
	if [[ ! "$composeVersion" == *"Docker Compose version v2"* ]]; then
		builtin echo -e "${start_error}We could not identify Composer V2 installed.${end_error}"
		builtin echo -e "Please make sure you are running an updated Docker version that includes Compose V2:"
		builtin echo -e "  Official Docker installation documentation: https://docs.docker.com/get-docker/"
		builtin echo -e "  Official Docker Compose V2 documentation: https://docs.docker.com/compose/reference/"
		exit
	fi

	# success
	builtin echo -e "${start_success}$(kool -v) installed successfully.${end_success}"
}

do_install
