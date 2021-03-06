package main

import (
	"flag"
	_ "github.com/ribice/gorsk/cmd/api/docs"
	"github.com/ribice/gorsk/pkg/api"
	"github.com/ribice/gorsk/pkg/utl/config"
	"log"
)

//go:generate swagger generate spec
func main() {

	cfgPath := flag.String("p", "./cmd/api/conf.local.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	checkErr(err)

	checkErr(api.Start(cfg))
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
