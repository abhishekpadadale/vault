package main

import (
	"log"
	"os"

	"github.com/abhishekpadadale/vault/api"
	"github.com/abhishekpadadale/vault/plugins/database/mysql"
)

func main() {
	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	err := mysql.Run(apiClientMeta.GetTLSConfig())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
