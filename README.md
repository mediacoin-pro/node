# Mediacoin Blockchain Node


## Install Node
##### Install Golang (version ≥ 1.11)
for Linux 
``` shell
apt-get install golang
```
or
``` shell
wget https://dl.google.com/go/go1.11.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.11.5.linux-amd64.tar.gz
```
or
follow the installation instructions: https://golang.org/doc/install

##### Install or update node 
``` shell
go get -u -v github.com/mediacoin-pro/node/
go build -o mdcnode github.com/mediacoin-pro/node/cmd/mdcnode
``` 

##### Show Node version, arguments
``` shell
./mdcnode -version
./mdcnode -help
```


## Start Node
``` shell
nohup ./mdcnode -http=127.0.0.1:8888 -dir=$HOME/mdc.db < /dev/null >/var/log/mdcnode.log 2>&1 &
``` 



