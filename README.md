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
go get -u -v github.com/mediacoin-pro/node/cmd/mdcnode
go build -o mdcnode github.com/mediacoin-pro/node/cmd/mdcnode
``` 

##### Show Node version, arguments
``` shell
./mdcnode -version
./mdcnode -help
```


## Start Node
``` shell
nohup ./mdcnode -http=127.0.0.1:8777 -dir=$HOME/mdc </dev/null >/var/log/mdcnode.log 2>&1 &
``` 


## Node REST API
``` 
http://127.0.0.1:8777/<command>? [&pretty] &<param>=<value>.... 
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

##### Generate address with Memo 
``` 
GET /address/?address&memo  
```

##### Get address info + memo code 
``` 
GET /address/?address=<address>&memo=<memo> 
```

##### Get transaction list by address (+memo)
``` 
GET /txs/?address=<address> [&memo=<num|hex>] [&limit=<int>] [&order="asc"|"desc"] [&offset=<hex>]
```

##### Generate new key pair, address by secret-phrase
``` 
GET /new-key?seed=<secret_phrase>
```

##### Register user in blockchain
``` 
POST /new-user?login=<login>&password=<password>
```

##### Transfer founds to address
``` 
POST /new-transfer? &(seed|login&password|private) &address=<address> [&memo=<num|hex>] &amount=<num> [&comment] [&nonce=<num|hex>] 
```

