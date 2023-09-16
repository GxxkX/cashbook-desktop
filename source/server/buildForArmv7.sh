sudo apt-get update
sudo apt-get install -y gcc-arm-linux-gnueabihf
export GOOS=linux
export GOARCH=arm
export GOARM=7
export CC=arm-linux-gnueabihf-gcc
go build -o cashbook *.go