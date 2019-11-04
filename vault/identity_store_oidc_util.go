// +build !enterprise

package vault

import (
	"github.com/abhishekpadadale/vault/helper/namespace"
)

func (i *IdentityStore) listNamespacePaths() []string {
	return []string{namespace.RootNamespace.Path}
}
