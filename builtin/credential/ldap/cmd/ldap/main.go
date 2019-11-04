package main

import (
	"os"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/abhishekpadadale/vault/api"
	"github.com/abhishekpadadale/vault/builtin/credential/ldap"
	"github.com/abhishekpadadale/vault/sdk/plugin"
)

func main() {
	apiClientMeta := &api.PluginAPIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := api.VaultPluginTLSProvider(tlsConfig)

	if err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: ldap.Factory,
		TLSProviderFunc:    tlsProviderFunc,
	}); err != nil {
		logger := hclog.New(&hclog.LoggerOptions{})

		logger.Error("plugin shutting down", "error", err)
		os.Exit(1)
	}
}
