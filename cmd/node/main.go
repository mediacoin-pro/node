package main

import (
	"flag"
	"fmt"

	"github.com/mediacoin-pro/core/chain"
	"github.com/mediacoin-pro/core/chain/bcstore"
	"github.com/mediacoin-pro/core/chain/replication"
	"github.com/mediacoin-pro/node/rest/restsrv"
)

const applicationName = "Mediacoin Blockchain Node v1.0"

func main() {
	var (
		argHelp    = flag.Bool("help", false, "Show this help")
		argVersion = flag.Bool("version", false, "Show software version")

		bcCfg   = chain.NewConfig()
		restCfg = restsrv.NewConfig()
	)
	flag.Parse()

	if *argHelp {
		fmt.Println(applicationName)
		flag.PrintDefaults()
		return
	}
	if *argVersion {
		fmt.Println(applicationName)
		return
	}

	//xlog.SetLogLevel(xlog.LevelInfo)

	//---- start node --------
	var bc = bcstore.NewChainStorage(bcCfg)

	go restsrv.StartServer(restCfg, bc)
	go replication.Start(bc)

	select {}
}
