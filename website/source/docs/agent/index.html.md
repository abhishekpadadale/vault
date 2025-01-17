---
layout: "docs"
page_title: "Vault Agent"
sidebar_title: "Vault Agent"
sidebar_current: "docs-agent"
description: |-
  Vault Agent is a client-side daemon that can be used to perform some Vault
  functionality automatically.
---

# Vault Agent

Vault Agent is a client daemon that provides the following features:

- <tt>[Auto-Auth][autoauth]</tt> - Automatically authenticate to Vault and manage the token renewal process for locally-retrieved dynamic secrets.
- <tt>[Caching][caching]</tt> - Allows client-side caching of responses containing newly created tokens and responses containing leased secrets generated off of these newly created tokens.

To get help, run:

```text
$ vault agent -h
```

## Auto-Auth

Vault Agent allows for easy authentication to Vault in a wide variety of
environments. Please see the [Auto-Auth docs][autoauth]
for information.

Auto-Auth functionality takes place within an `auto_auth` configuration stanza.

## Caching

Vault Agent allows client-side caching of responses containing newly created tokens
and responses containing leased secrets generated off of these newly created tokens.
Please see the [Caching docs][caching] for information.

## Configuration

These are the currently-available general configuration option:

- `vault` <tt>([vault][vault]: \<optional\>)</tt> - Specifies the remote Vault server the Agent connects to.

- `auto_auth` <tt>([auto_auth][autoauth]: \<optional\>)</tt> - Specifies the method and other options used for Auto-Auth functionality.

- `cache` <tt>([cache][caching]: \<optional\>)</tt> - Specifies options used for Caching functionality.

- `listener` <tt>([listener][listener]: \<optional\>)</tt> - Specifies the addresses and ports on which the Agent will respond to requests.

- `pid_file` `(string: "")` - Path to the file in which the agent's Process ID
  (PID) should be stored

- `exit_after_auth` `(bool: false)` - If set to `true`, the agent will exit
  with code `0` after a single successful auth, where success means that a
  token was retrieved and all sinks successfully wrote it

- `template` <tt>([template][template`]: \<optional\>)</tt> - Specifies options used for templating Vault secrets to files.

### vault Stanza

There can at most be one top level `vault` block and it has the following
configuration entries:

- `address (string: optional)` - The address of the Vault server. This should
  be a complete URL such as `https://127.0.0.1:8200`. This value can be
  overridden by setting the `VAULT_ADDR` environment variable.

- `ca_cert (string: optional)` - Path on the local disk to a single PEM-encoded
  CA certificate to verify the Vault server's SSL certificate. This value can
  be overridden by setting the `VAULT_CACERT` environment variable.

- `ca_path (string: optional)` - Path on the local disk to a directory of
  PEM-encoded CA certificates to verify the Vault server's SSL certificate.
  This value can be overridden by setting the `VAULT_CAPATH` environment
  variable.

- `client_cert (string: option)` - Path on the local disk to a single
  PEM-encoded CA certificate to use for TLS authentication to the Vault server.
  This value can be overridden by setting the `VAULT_CLIENT_CERT` environment
  variable.

- `client_key (string: option)` - Path on the local disk to a single
  PEM-encoded private key matching the client certificate from `client_cert`.
  This value can be overridden by setting the `VAULT_CLIENT_KEY` environment
  variable.

- `tls_skip_verify (string: optional)` - Disable verification of TLS
  certificates. Using this option is highly discouraged as it decreases the
  security of data transmissions to and from the Vault server. This value can
  be overridden by setting the `VAULT_SKIP_VERIFY` environment variable.

- `tls_server_name (string: optional)` - Name to use as the SNI host when
  connecting via TLS. This value can be overridden by setting the
  `VAULT_TLS_SERVER_NAME` environment variable.

### listener Stanza

Agent supports one or more [listener][listener_main] stanzas.  In addition to
the standard listener configuration, an Agent's listener configuration also
supports an additional optional entry:

- `require_request_header (bool: false)` - Require that all incoming HTTP
  requests on this listener must have an `X-Vault-Request: true` header entry.
  Using this option offers an additional layer of protection from Service Side
  Request Forgery attacks.  Requests on the listener that do not have the proper
  `X-Vault-Request` header will fail, with a HTTP response status code of `412:
  Precondition Failed`.

## Example Configuration

An example configuration, with very contrived values, follows:

```python
pid_file = "./pidfile"

vault {
        address = "https://127.0.0.1:8200"
}

auto_auth {
        method "aws" {
                mount_path = "auth/aws-subaccount"
                config = {
                        type = "iam"
                        role = "foobar"
                }
        }

        sink "file" {
                config = {
                        path = "/tmp/file-foo"
                }
        }

        sink "file" {
                wrap_ttl = "5m"
                aad_env_var = "TEST_AAD_ENV"
                dh_type = "curve25519"
                dh_path = "/tmp/file-foo-dhpath2"
                config = {
                        path = "/tmp/file-bar"
                }
        }
}

cache {
        use_auto_auth_token = true
}

listener "unix" {
         address = "/path/to/socket"
         tls_disable = true
}

listener "tcp" {
         address = "127.0.0.1:8100"
         tls_disable = true
}

template {
  source      = "/etc/vault/server.key.ctmpl"
  destination = "/etc/vault/server.key"
}

template {
  source      = "/etc/vault/server.crt.ctmpl"
  destination = "/etc/vault/server.crt"
}
```

[vault]: /docs/agent/index.html#vault-stanza
[autoauth]: /docs/agent/autoauth/index.html
[caching]: /docs/agent/caching/index.html
[template]: /docs/agent/template/index.html
[listener]: /docs/agent/index.html#listener-stanza
[listener_main]: /docs/configuration/listener/tcp.html
