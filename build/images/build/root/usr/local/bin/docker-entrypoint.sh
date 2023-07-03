#!/bin/bash
set -e
source /usr/local/bin/setup-credentials-helper.sh

if [ -z "$1" ] ; then
    exec "/bin/sh"
else
    exec "$@"
fi
