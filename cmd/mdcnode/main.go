package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mediacoin-pro/core/chain/bcstore"
	"github.com/mediacoin-pro/core/chain/replication"
	"github.com/mediacoin-pro/core/common/xlog"
	"github.com/mediacoin-pro/node/rest/restsrv"
)

const applicationName = "Mediacoin Blockchain Node v1.2"

func main() {
	var (
		argHelp     = flag.Bool("help", false, "Show this help")
		argVersion  = flag.Bool("version", false, "Show software version")
		argDataDir  = flag.String("dir", os.Getenv("HOME")+"/mdc", "Node data dir")
		argLogLevel = flag.Int("loglevel", xlog.LevelInfo, "Log level (1-fatal, 2-error, 3-warning, 4-info, 5-debug, 6-trace)")
		restCfg     = restsrv.NewConfig()
	)
	flag.Parse()

	if *argHelp {
		fmt.Println(applicationName + "\n\nUsage:\n")
		flag.PrintDefaults()
		return
	}
	if *argVersion {
		fmt.Println(applicationName)
		return
	}

	xlog.SetLogLevel(*argLogLevel)

	//---- start node --------
	if err := os.Mkdir(*argDataDir, 0755); err != nil && !os.IsExist(err) {
		log.Panic(err)
	}
	var bc = bcstore.NewChainStorage(*argDataDir+"/bc", nil)

	go restsrv.StartServer(restCfg, bc)
	go replication.Start(bc)

	select {}
}
