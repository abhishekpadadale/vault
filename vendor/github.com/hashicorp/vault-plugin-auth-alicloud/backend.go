package alicloud

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/abhishekpadadale/vault/sdk/framework"
	"github.com/abhishekpadadale/vault/sdk/logical"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	client := cleanhttp.DefaultClient()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	b := newBackend(client)
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// newBackend exists for testability. It allows us to inject a fake client.
func newBackend(client *http.Client) *backend {
	b := &backend{
		identityClient: client,
	}
	b.Backend = &framework.Backend{
		AuthRenew: b.pathLoginRenew,
		Help:      backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
		},
		Paths: []*framework.Path{
			pathLogin(b),
			pathListRole(b),
			pathListRoles(b),
			pathRole(b),
		},
		BackendType: logical.TypeCredential,
	}
	return b
}

type backend struct {
	*framework.Backend

	identityClient *http.Client
}

const backendHelp = `
That AliCloud RAM auth method allows entities to authenticate based on their
identity and pre-configured roles.
`
