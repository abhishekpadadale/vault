package cf

import (
	"context"

	"github.com/abhishekpadadale/vault/sdk/framework"
	"github.com/abhishekpadadale/vault/sdk/logical"
)

const (
	// These env vars are used frequently to pull the client certificate and private key
	// from CF containers; thus are placed here for ease of discovery and use from
	// outside packages.
	EnvVarInstanceCertificate = "CF_INSTANCE_CERT"
	EnvVarInstanceKey         = "CF_INSTANCE_KEY"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := &backend{}
	b.Backend = &framework.Backend{
		AuthRenew: b.pathLoginRenew,
		Help:      backendHelp,
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{"config"},
			Unauthenticated: []string{"login"},
		},
		Paths: []*framework.Path{
			b.pathConfig(),
			b.pathListRoles(),
			b.pathRoles(),
			b.pathLogin(),
		},
		BackendType: logical.TypeCredential,
	}
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

type backend struct {
	*framework.Backend
}

const backendHelp = `
The CF auth backend supports logging in using CF's identity service.
Once a CA certificate is configured, and Vault is configured to consume
CF's API, CF's instance identity credentials can be used to authenticate.'
`
