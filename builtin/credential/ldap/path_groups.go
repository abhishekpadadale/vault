package ldap

import (
	"context"
	"strings"

	"github.com/abhishekpadadale/vault/sdk/framework"
	"github.com/abhishekpadadale/vault/sdk/helper/policyutil"
	"github.com/abhishekpadadale/vault/sdk/logical"
)

func pathGroupsList(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "groups/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathGroupList,
		},

		HelpSynopsis:    pathGroupHelpSyn,
		HelpDescription: pathGroupHelpDesc,
		DisplayAttrs: &framework.DisplayAttributes{
			Navigation: true,
			ItemType:   "Group",
		},
	}
}

func pathGroups(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `groups/(?P<name>.+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the LDAP group.",
			},

			"policies": {
				Type:        framework.TypeCommaStringSlice,
				Description: "Comma-separated list of policies associated to the group.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: b.pathGroupDelete,
			logical.ReadOperation:   b.pathGroupRead,
			logical.UpdateOperation: b.pathGroupWrite,
		},

		HelpSynopsis:    pathGroupHelpSyn,
		HelpDescription: pathGroupHelpDesc,
		DisplayAttrs: &framework.DisplayAttributes{
			Action:   "Create",
			ItemType: "Group",
		},
	}
}

func (b *backend) Group(ctx context.Context, s logical.Storage, n string) (*GroupEntry, error) {
	entry, err := s.Get(ctx, "group/"+n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result GroupEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathGroupDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, "group/"+d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathGroupRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	groupname := d.Get("name").(string)

	cfg, err := b.Config(ctx, req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return logical.ErrorResponse("ldap backend not configured"), nil
	}
	if !*cfg.CaseSensitiveNames {
		groupname = strings.ToLower(groupname)
	}

	group, err := b.Group(ctx, req.Storage, groupname)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"policies": group.Policies,
		},
	}, nil
}

func (b *backend) pathGroupWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	groupname := d.Get("name").(string)

	cfg, err := b.Config(ctx, req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return logical.ErrorResponse("ldap backend not configured"), nil
	}
	if !*cfg.CaseSensitiveNames {
		groupname = strings.ToLower(groupname)
	}

	// Store it
	entry, err := logical.StorageEntryJSON("group/"+groupname, &GroupEntry{
		Policies: policyutil.ParsePolicies(d.Get("policies")),
	})
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathGroupList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	keys, err := logical.CollectKeysWithPrefix(ctx, req.Storage, "group/")
	if err != nil {
		return nil, err
	}
	for i := range keys {
		keys[i] = strings.TrimPrefix(keys[i], "group/")
	}
	return logical.ListResponse(keys), nil
}

type GroupEntry struct {
	Policies []string
}

const pathGroupHelpSyn = `
Manage additional groups for users allowed to authenticate.
`

const pathGroupHelpDesc = `
This endpoint allows you to create, read, update, and delete configuration
for LDAP groups that are allowed to authenticate, and associate policies to
them.

Deleting a group will not revoke auth for prior authenticated users in that
group. To do this, do a revoke on "login/<username>" for
the usernames you want revoked.
`
