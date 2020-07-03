package main

// DefaultEnv is supposed to contain all default
// .env file injected upon build time
const DefaultEnv string = `
# Only print out the generated commands instead of running them
# KOOL_DEBUG=0
# Prints out all commands and their output
# KOOL_VERBOSE=0

# Optional, default is folder name
# KOOL_NAME=custom_name

# Host user for mapping execution within the containers using kool images
# KOOL_ASUSER=$UID

# Docker Run flags
# KOOL_DISABLE_TTY=1

# Docker Composer
# KOOL_GLOBAL_NETWORK="custom-network" # by default will be generated on the fly
`
