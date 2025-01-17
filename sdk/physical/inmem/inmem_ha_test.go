package inmem

import (
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/abhishekpadadale/vault/sdk/helper/logging"
	"github.com/abhishekpadadale/vault/sdk/physical"
)

func TestInmemHA(t *testing.T) {
	logger := logging.NewVaultLogger(log.Debug)

	inm, err := NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	// Use the same inmem backend to acquire the same set of locks
	physical.ExerciseHABackend(t, inm.(physical.HABackend), inm.(physical.HABackend))
}
