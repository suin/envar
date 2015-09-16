APP_NAME="envar"
APP_OS="linux darwin windows"
APP_ARCH="386 amd64"

# define app version
git fetch --tags
APP_VERSION=$(git describe --tags)
echo APP_VERSION is $APP_VERSION

# download go
pushd $HOME
curl -s -o go.tar.gz https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
tar xzf go.tar.gz
export GOROOT=$HOME/go
export PATH=$GOROOT/bin:$PATH
go version
popd

# get dependencies
go get ./...

# build artifacts
go get github.com/mitchellh/gox
gox -os="$APP_OS" -arch="$APP_ARCH" -output="artifacts/{{.OS}}-{{.Arch}}/$APP_NAME" -ldflags "-X main.Version '$APP_VERSION'"
find artifacts
