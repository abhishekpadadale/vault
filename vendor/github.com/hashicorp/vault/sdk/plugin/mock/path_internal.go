package mock

import (
	"context"

	"github.com/abhishekpadadale/vault/sdk/framework"
	"github.com/abhishekpadadale/vault/sdk/logical"
)

// pathInternal is used to test viewing internal backend values. In this case,
// it is used to test the invalidate func.
func pathInternal(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "internal",
		Fields: map[string]*framework.FieldSchema{
			"value": &framework.FieldSchema{Type: framework.TypeString},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathInternalUpdate,
			logical.ReadOperation:   b.pathInternalRead,
		},
	}
}

func (b *backend) pathInternalUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	value := data.Get("value").(string)
	b.internal = value
	// Return the secret
	return nil, nil
}

func (b *backend) pathInternalRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Return the secret
	return &logical.Response{
		Data: map[string]interface{}{
			"value": b.internal,
		},
	}, nil
}
