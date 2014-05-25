if [ ! -f tools/windowsconstgen.go ]; then
	echo error: $0 must be run from the package source root 1>&2
	exit 1
fi
set -e
if [ x$GOOS = xwindows ]; then
	# have to invoke go run with the host $GOOS/$GOARCH
	GOOS= GOARCH= go run tools/windowsconstgen.go . 386 "$@"
	GOOS= GOARCH= go run tools/windowsconstgen.go . amd64 "$@"
fi
cd test
go build "$@"
