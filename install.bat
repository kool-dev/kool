@echo off

setlocal

REM curl -fsSL https://kool.dev/install -o get-kool.sh

set DEFAULT_DOWNLOAD_URL="https://github.com/kool-dev/kool/releases/latest/download"
if "%DOWNLOAD_URL%" == "" set DOWNLOAD_URL=%DEFAULT_DOWNLOAD_URL%

if "%PROCESSOR_ARCHITECTURE%" == "x86" (set DEFAULT_BIN_PATH="%ProgramFiles(x86)%\kool-dev") else (set DEFAULT_BIN_PATH="%ProgramFiles%\kool-dev")
DEFAULT_BIN=%DEFAULT_BIN_PATH%
if "%BIN_PATH%" == "" set BIN_PATH=%DEFAULT_BIN%

:command_exists
	call -v %* > NULL

:do_install
	if "%PROCESSOR_ARCHITECTURE%" == "x86" (set ARCH=386) else (set ARCH=amd64)
	set PLAT="windows"

    mkdir %BIN_PATH%
	bitsadmin /transfer KoolDevDownload /dynamic /download /priority foreground "%DOWNLOAD_URL%/kool-%PLAT%-%ARCH%" "%BIN_PATH%\kool.exe"
	setx path "%BIN_PATH%;%PATH%"

goto do_install
