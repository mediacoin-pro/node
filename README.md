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
nohup ./mdcnode -http=127.0.0.1:8888 -dir=$HOME/mdc </dev/null >/var/log/mdcnode.log 2>&1 &
``` 


## Node REST API
``` 
http://127.0.0.1:8888/<command>? [&pretty] &<param>=<value>.... 
```

##### Get general node and blockchain information
``` 
GET /info 
```

##### Get block 
``` 
GET /block/<blockNum> 
```

##### Get blocks
``` 
GET /blocks?offset=<blockNum>&limit=<countBlocks> 
```

##### Get transaction 
``` 
GET /tx/<txID:hex> 
```

##### Get address info 
``` 
GET /address/<address> 
GET /address/@<username>
GET /address/0x<userID:hex> 
GET /address/?address=<address> 
```

##### Get address info + memo code 
``` 
GET /address/?address=<address>&memo=<memo> 
```



