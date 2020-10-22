wget -q -O dotnet-sdk.tar.gz https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz
sudo tar zxf dotnet-sdk.tar.gz -C /usr/local
export PATH=/usr/local/go/bin:$PATH
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
go version