APP_NAME="envar"
APP_OS="linux darwin windows"
APP_ARCH="386 amd64"

pushd ~/
curl -s -o go.tar.gz https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
tar xzf go.tar.gz
export GOROOT=~/go
export PATH=$GOROOT/bin:$PATH
go version
popd

go get github.com/mitchellh/gox
gox -build-toolchain -os="$APP_OS" -arch="$APP_ARCH"

APP_VERSION=$(git describe --tags)
echo APP_VERSION is $APP_VERSION

gox -os="$APP_OS" -arch="$APP_ARCH" -output="artifacts/{{.OS}}-{{.Arch}}/$APP_NAME" -ldflags "-X main.Version '$APP_VERSION'"
find artifacts
