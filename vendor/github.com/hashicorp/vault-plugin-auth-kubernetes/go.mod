module github.com/abhishekpadadale/vault-plugin-auth-kubernetes

go 1.12

require (
	github.com/briankassouf/jose v0.9.2-0.20180619214549-d2569464773f
	github.com/go-test/deep v1.0.2-0.20181118220953-042da051cf31
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/go-cleanhttp v0.5.1
	github.com/hashicorp/go-hclog v0.8.0
	github.com/hashicorp/go-multierror v1.0.0
	github.com/hashicorp/go-sockaddr v1.0.2
	github.com/abhishekpadadale/vault/api v1.0.5-0.20190814205728-e9c5cd8aca98
	github.com/abhishekpadadale/vault/sdk v0.1.14-0.20190814205504-1cad00d1133b
	github.com/mitchellh/mapstructure v1.1.2
	k8s.io/api v0.0.0-20190409092523-d687e77c8ae9
	k8s.io/apimachinery v0.0.0-20190409092423-760d1845f48b
)
