#!/usr/bin/env bash

if [ ! -f make ]; then
    echo 'make must be run within its container folder' 1>&2
    exit 1
fi

export CGO_ENABLED=0
export GOROOT=/usr/local/go

OLDPATH="$PATH"
export PATH=$PATH:$GOROOT/bin
OLDGOPATH="$GOPATH"
export GOPATH=`pwd`

go install code.google.com/p/freetype-go/freetype
go install utility/process
go install signer
go install example/png

export GOPATH="$OLDGOPATH"
export PATH="$OLDPATH"

echo 'finished.'
