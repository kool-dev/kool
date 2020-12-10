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

command_exists() {
	command -v "$@" > /dev/null 2>&1
}

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

	echo "Downloading latest binary (kool-$PLAT-$ARCH)..."

	# TODO: fallback to wget if no curl available
	rm -f /tmp/kool_binary
	curl -fsSL "$DOWNLOAD_URL/kool-$PLAT-$ARCH" -o /tmp/kool_binary

	# check for running kool process which would prevent
	# replacing existing version under Linux.
	if command_exists kool && ! is_darwin; then
		running=$(ps aux | grep kool | grep -v grep | wc -l | awk '{ print $1 }')
		if [ "$running" != "0" ]; then
			echo -e "\033[31;31mThere is a kool process still running. You might need to stop them before we replace the current binary.\033[0m"
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
	builtin echo -e "${start_success}$(kool -v) installed successfully.${end_success}"

	# TODO: use command_exists to check and alert about docker/docker-compose
	# being available.
}

do_install
