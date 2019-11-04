package plugin

import (
	"github.com/abhishekpadadale/vault-plugin-secrets-ad/plugin/client"
)

type configuration struct {
	PasswordConf          *passwordConf
	ADConf                *client.ADConf
	LastRotationTolerance int
}
