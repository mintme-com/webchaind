#!/usr/bin/env bash

set -e

OUTPUT="$1"

if [ ! "$OUTPUT" == "build" ] && [ ! "$OUTPUT" == "install" ]; then
	echo "Specify 'install' or 'build' as first argument."
	exit 1
else
	echo "With SputnikVM, running webchaind $OUTPUT ..."
fi

OS=`uname -s`

geth_path="github.com/webchain-network/webchaind"
sputnik_path="github.com/webchain-network/sputnikvm-ffi"
sputnik_dir="$GOPATH/src/$geth_path/vendor/$sputnik_path"

geth_bindir="$GOPATH/src/$geth_path/bin"

echo "Building SputnikVM"
make -C "$sputnik_dir/c"

echo "Doing webchaind $OUTPUT ..."
cd "$GOPATH/src/$geth_path"

LDFLAGS="$sputnik_dir/c/libsputnikvm.a "
case $OS in
	"Linux")
		LDFLAGS+="-ldl"
		;;

	"Darwin")
		LDFLAGS+="-ldl -lresolv"
		;;

    CYGWIN*|MINGW32*|MSYS*)
		LDFLAGS="-Wl,--allow-multiple-definition $sputnik_dir/c/sputnikvm.lib -lws2_32 -luserenv"
		;;
esac



if [ "$OUTPUT" == "install" ]; then
	CGO_CFLAGS_ALLOW='-maes.*' CGO_LDFLAGS=$LDFLAGS go install -ldflags '-X main.Version='$(git describe --tags) -tags="sputnikvm netgo" ./cmd/webchaind
elif [ "$OUTPUT" == "build" ]; then
	mkdir -p "$geth_bindir"
	CGO_CFLAGS_ALLOW='-maes.*' CGO_LDFLAGS=$LDFLAGS go build -ldflags '-X main.Version='$(git describe --tags) -o $geth_bindir/webchaind -tags="sputnikvm netgo" ./cmd/webchaind
fi
