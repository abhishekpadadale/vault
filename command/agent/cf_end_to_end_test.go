package agent

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	log "github.com/hashicorp/go-hclog"
	credCF "github.com/abhishekpadadale/vault-plugin-auth-cf"
	"github.com/abhishekpadadale/vault-plugin-auth-cf/testing/certificates"
	cfAPI "github.com/abhishekpadadale/vault-plugin-auth-cf/testing/cf"
	"github.com/abhishekpadadale/vault/api"
	"github.com/abhishekpadadale/vault/command/agent/auth"
	agentcf "github.com/abhishekpadadale/vault/command/agent/auth/cf"
	"github.com/abhishekpadadale/vault/command/agent/sink"
	"github.com/abhishekpadadale/vault/command/agent/sink/file"
	vaulthttp "github.com/abhishekpadadale/vault/http"
	"github.com/abhishekpadadale/vault/sdk/helper/logging"
	"github.com/abhishekpadadale/vault/sdk/logical"
	"github.com/abhishekpadadale/vault/vault"
)

func TestCFEndToEnd(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)

	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{
			"cf": credCF.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)
	client := cores[0].Client
	if err := client.Sys().EnableAuthWithOptions("cf", &api.EnableAuthOptions{
		Type: "cf",
	}); err != nil {
		t.Fatal(err)
	}

	testIPAddress := "127.0.0.1"

	// Generate some valid certs that look like the ones we get from CF.
	testCFCerts, err := certificates.Generate(cfAPI.FoundServiceGUID, cfAPI.FoundOrgGUID, cfAPI.FoundSpaceGUID, cfAPI.FoundAppGUID, testIPAddress)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := testCFCerts.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	// Start a mock server representing their API.
	mockCFAPI := cfAPI.MockServer(false)
	defer mockCFAPI.Close()

	// Configure a CA certificate like a Vault operator would in setting up CF.
	if _, err := client.Logical().Write("auth/cf/config", map[string]interface{}{
		"identity_ca_certificates": testCFCerts.CACertificate,
		"cf_api_addr":             mockCFAPI.URL,
		"cf_username":             cfAPI.AuthUsername,
		"cf_password":             cfAPI.AuthPassword,
	}); err != nil {
		t.Fatal(err)
	}

	// Configure a role to be used for logging in, another thing a Vault operator would do.
	if _, err := client.Logical().Write("auth/cf/roles/test-role", map[string]interface{}{
		"bound_instance_ids":     cfAPI.FoundServiceGUID,
		"bound_organization_ids": cfAPI.FoundOrgGUID,
		"bound_space_ids":        cfAPI.FoundSpaceGUID,
		"bound_application_ids":  cfAPI.FoundAppGUID,
	}); err != nil {
		t.Fatal(err)
	}

	os.Setenv(credCF.EnvVarInstanceCertificate, testCFCerts.PathToInstanceCertificate)
	os.Setenv(credCF.EnvVarInstanceKey, testCFCerts.PathToInstanceKey)

	ctx, cancelFunc := context.WithCancel(context.Background())
	timer := time.AfterFunc(30*time.Second, func() {
		cancelFunc()
	})
	defer timer.Stop()

	am, err := agentcf.NewCFAuthMethod(&auth.AuthConfig{
		MountPath: "auth/cf",
		Config: map[string]interface{}{
			"role": "test-role",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ahConfig := &auth.AuthHandlerConfig{
		Logger: logger.Named("auth.handler"),
		Client: client,
	}

	ah := auth.NewAuthHandler(ahConfig)
	go ah.Run(ctx, am)
	defer func() {
		<-ah.DoneCh
	}()

	tmpFile, err := ioutil.TempFile("", "auth.tokensink.test.")
	if err != nil {
		t.Fatal(err)
	}
	tokenSinkFileName := tmpFile.Name()
	tmpFile.Close()
	os.Remove(tokenSinkFileName)
	t.Logf("output: %s", tokenSinkFileName)

	config := &sink.SinkConfig{
		Logger: logger.Named("sink.file"),
		Config: map[string]interface{}{
			"path": tokenSinkFileName,
		},
		WrapTTL: 10 * time.Second,
	}

	fs, err := file.NewFileSink(config)
	if err != nil {
		t.Fatal(err)
	}
	config.Sink = fs

	ss := sink.NewSinkServer(&sink.SinkServerConfig{
		Logger: logger.Named("sink.server"),
		Client: client,
	})
	go ss.Run(ctx, ah.OutputCh, []*sink.SinkConfig{config})
	defer func() {
		<-ss.DoneCh
	}()

	if stat, err := os.Lstat(tokenSinkFileName); err == nil {
		t.Fatalf("expected err but got %s", stat)
	} else if !os.IsNotExist(err) {
		t.Fatal("expected notexist err")
	}

	// Wait 2 seconds for the env variables to be detected and an auth to be generated.
	time.Sleep(time.Second * 2)

	token, err := readToken(tokenSinkFileName)
	if err != nil {
		t.Fatal(err)
	}

	if token.Token == "" {
		t.Fatal("expected token but didn't receive it")
	}
}
