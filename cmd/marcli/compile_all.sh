MARCLI_VERSION="$1"
if [ "$MARCLI_VERSION" = "" ]; then
  echo "Must pass a version number, e.g. v1.0.0"
  exit 1
fi

echo "Compiling marcli $MARCLI_VERSION ..."
# See this blog post for information about Go's ldflags parameter
# https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
MARCLI_LDFLAGS="-X 'main.Version=$MARCLI_VERSION'"

GOOS=darwin GOARCH=amd64 go build -ldflags="$MARCLI_LDFLAGS" -o marcli_amd64
GOOS=darwin GOARCH=arm64 go build -ldflags="$MARCLI_LDFLAGS" -o marcli_arm64
GOOS=linux GOARCH=amd64 go build -ldflags="$MARCLI_LDFLAGS" -o marcli_linux
GOOS=windows GOARCH=386 go build -ldflags="$MARCLI_LDFLAGS" -o marcli.exe

# Create the universal executable for Mac
lipo -create -output marcli marcli_amd64 marcli_arm64