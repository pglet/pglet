curl -fsSL https://appveyordownloads.blob.core.windows.net/misc/goreleaser-m1.tar.gz -o /tmp/goreleaser-m1.tar.gz
sudo tar zxf /tmp/goreleaser-m1.tar.gz -C /usr/local
ls -al /usr/local
/usr/local/goreleaser --version