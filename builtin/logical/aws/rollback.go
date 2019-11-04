package aws

import (
	"context"
	"fmt"

	"github.com/abhishekpadadale/vault/sdk/framework"
	"github.com/abhishekpadadale/vault/sdk/helper/consts"
	"github.com/abhishekpadadale/vault/sdk/logical"
)

func (b *backend) walRollback(ctx context.Context, req *logical.Request, kind string, data interface{}) error {
	walRollbackMap := map[string]framework.WALRollbackFunc{
		"user": b.pathUserRollback,
	}

	if !b.System().LocalMount() && b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary|consts.ReplicationPerformanceStandby) {
		return nil
	}

	f, ok := walRollbackMap[kind]
	if !ok {
		return fmt.Errorf("unknown type to rollback")
	}

	return f(ctx, req, kind, data)
}
