#!/usr/bin/env bash
set -e

# curl -fsSL https://kool.dev/install -o get-kool.sh
# sudo sh get-kool.sh

DEFAULT_DOWNLOAD_URL="https://downloads.kool.dev"
if [ -z "$DOWNLOAD_URL" ]; then
	DOWNLOAD_URL=$DEFAULT_DOWNLOAD_URL
fi

DEFAULT_BIN="/usr/local/bin/kool"
if [ -z "$DOWNLOAD_URL" ]; then
	BIN_PATH=$DEFAULT_BIN
fi

command_exists() {
	command -v "$@" > /dev/null 2>&1
}

is_wsl() {
	case "$(uname -r)" in
	*microsoft* ) true ;; # WSL 2
	*Microsoft* ) true ;; # WSL 1
	* ) false;;
	esac
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
	if [ is_darwin ]; then
		PLAT="darwin"
	elif [ is_wsl ]; then
		PLAT="wsl"
	fi

	wget -O $BIN_PATH "$DOWNLOAD_URL/$PLAT-$ARCH-kool"
}

do_install
