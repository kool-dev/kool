#!/usr/bin/env bash
set -e

# curl -fsSL https://kool.dev/install -o get-kool.sh
# sudo sh get-kool.sh

DEFAULT_DOWNLOAD_URL="https://github.com/kool-dev/kool/releases/latest/download"
if [ -z "$DOWNLOAD_URL" ]; then
	DOWNLOAD_URL=$DEFAULT_DOWNLOAD_URL
fi

DEFAULT_BIN="/usr/local/bin/kool"
if [ -z "$BIN_PATH" ]; then
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

	# wget -O $BIN_PATH "$DOWNLOAD_URL/kool-$PLAT-$ARCH"
	curl -fsSL "$DOWNLOAD_URL/kool-$PLAT-$ARCH" -o $BIN_PATH
	chmod +x $BIN_PATH
}

do_install
