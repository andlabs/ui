if [ ! -f tools/windowsconstgen.go ]; then
	echo error: $0 must be run from the package source root 1>&2
	exit 1
fi
set -e
if [ x$GOOS = xwindows ]; then
	# have to build windowsconstgen as the host, otherwise weird things happen
	wcg=`mktemp /tmp/windowsconstgenXXXXXXXXXXXX`
	GOOS= GOARCH= go build -o $wcg tools/windowsconstgen.go
	# but we can run it regardless of $GOOS/$GOARCH
	$wcg . 386 "$@"
	$wcg . amd64 "$@"
	rm $wcg
fi
cd test
go build "$@"
