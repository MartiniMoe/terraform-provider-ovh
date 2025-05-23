package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceCloudProjectDatabaseMongodbUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudProjectDatabaseMongodbUserCreate,
		ReadContext:   resourceCloudProjectDatabaseMongodbUserRead,
		DeleteContext: resourceCloudProjectDatabaseMongodbUserDelete,
		UpdateContext: resourceCloudProjectDatabaseMongodbUserUpdate,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectDatabaseMongodbUserImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Description: "Id of the database cluster",
				ForceNew:    true,
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the user",
				ForceNew:    true,
				Required:    true,
				StateFunc: func(new interface{}) string {
					return new.(string) + "@admin"
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new+"@admin" == old
				},
			},
			"password_reset": {
				Type:        schema.TypeString,
				Description: "Arbitrary string to change to trigger a password update",
				Optional:    true,
			},
			"roles": {
				Type:        schema.TypeSet,
				Description: "Roles the user belongs to with the authentication database",
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateCloudProjectDatabaseMongodbUserAuthenticationDatabase,
				},
			},

			//Computed
			"created_at": {
				Type:        schema.TypeString,
				Description: "Date of the creation of the user",
				Computed:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "Password of the user",
				Sensitive:   true,
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Current status of the user",
				Computed:    true,
			},
		},

		CustomizeDiff: customdiff.ComputedIf(
			"password",
			func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
				return d.HasChange("password_reset")
			},
		),
	}
}

func resourceCloudProjectDatabaseMongodbUserImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return importCloudProjectDatabaseUser(d, meta)
}

func resourceCloudProjectDatabaseMongodbUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	if name == "admin" {
		diags := dataSourceCloudProjectDatabaseMongodbUserRead(ctx, d, meta)
		if diags.HasError() {
			return diags
		}
		return resourceCloudProjectDatabaseMongodbUserUpdate(ctx, d, meta)
	}

	serviceName := d.Get("service_name").(string)
	clusterID := d.Get("cluster_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/database/mongodb/%s/user",
		url.PathEscape(serviceName),
		url.PathEscape(clusterID),
	)
	params := (&CloudProjectDatabaseMongodbUserCreateOpts{}).fromResource(d)
	res := &CloudProjectDatabaseUserResponse{}

	return diag.FromErr(
		retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate),
			func() *retry.RetryError {
				log.Printf("[DEBUG] Will create user: %+v for cluster %s from project %s", params, clusterID, serviceName)
				rErr := postFuncCloudProjectDatabaseUser(ctx, d, meta, "mongodb", endpoint, params, res, schema.TimeoutCreate)
				if rErr != nil {
					if errOvh, ok := rErr.(*ovh.APIError); ok && (errOvh.Code == 409) {
						return retry.RetryableError(rErr)
					}
					return retry.NonRetryableError(rErr)
				}

				d.SetId(res.Id)
				readDiags := resourceCloudProjectDatabaseMongodbUserRead(ctx, d, meta)
				rErr = diagnosticsToError(readDiags)
				if rErr != nil {
					return retry.NonRetryableError(rErr)
				}
				return nil
			},
		),
	)
}

func resourceCloudProjectDatabaseMongodbUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterID := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/mongodb/%s/user/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterID),
		url.PathEscape(id),
	)
	res := &CloudProjectDatabaseMongodbUserResponse{}

	log.Printf("[DEBUG] Will read user %s from cluster %s from project %s", id, clusterID, serviceName)
	if err := config.OVHClient.GetWithContext(ctx, endpoint, res); err != nil {
		return diag.FromErr(helpers.CheckDeleted(d, err, endpoint))
	}

	for k, v := range res.toMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	log.Printf("[DEBUG] Read user %+v", res)
	return nil
}

func resourceCloudProjectDatabaseMongodbUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterID := d.Get("cluster_id").(string)
	isAdmin := d.Get("name").(string) == "admin@admin"
	passwordReset := d.HasChange("password_reset") || (d.IsNewResource() && isAdmin)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/mongodb/%s/user/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterID),
		url.PathEscape(id),
	)
	params := (&CloudProjectDatabaseMongodbUserUpdateOpts{}).fromResource(d)

	return diag.FromErr(
		retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate),
			func() *retry.RetryError {
				log.Printf("[DEBUG] Will update user: %+v from cluster %s from project %s", params, clusterID, serviceName)
				rErr := config.OVHClient.PutWithContext(ctx, endpoint, params, nil)
				if rErr != nil {
					if errOvh, ok := rErr.(*ovh.APIError); ok && (errOvh.Code == 409) {
						return retry.RetryableError(rErr)
					}
					return retry.NonRetryableError(fmt.Errorf("calling Put %s with params %s:\n\t %q", endpoint, params, rErr))
				}

				log.Printf("[DEBUG] Waiting for user %s to be READY", id)
				rErr = waitForCloudProjectDatabaseUserReady(ctx, config.OVHClient, serviceName, "mongodb", clusterID, id, d.Timeout(schema.TimeoutUpdate))
				if rErr != nil {
					return retry.NonRetryableError(fmt.Errorf("timeout while waiting user %s to be READY: %w", id, rErr))
				}
				log.Printf("[DEBUG] user %s is READY", id)

				if passwordReset {
					pwdResetEndpoint := endpoint + "/credentials/reset"
					res := &CloudProjectDatabaseUserResponse{}
					log.Printf("[DEBUG] Will update user password for cluster %s from project %s", clusterID, serviceName)
					err := postFuncCloudProjectDatabaseUser(ctx, d, meta, "mongodb", pwdResetEndpoint, nil, res, schema.TimeoutUpdate)
					if err != nil {
						if errOvh, ok := err.(*ovh.APIError); ok && (errOvh.Code == 409) {
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
				}

				readDiags := resourceCloudProjectDatabaseMongodbUserRead(ctx, d, meta)
				rErr = diagnosticsToError(readDiags)
				if rErr != nil {
					return retry.NonRetryableError(rErr)
				}
				return nil
			},
		),
	)
}

func resourceCloudProjectDatabaseMongodbUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	if name == "admin@admin" {
		d.SetId("")
		return nil
	}

	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	clusterID := d.Get("cluster_id").(string)
	id := d.Id()

	endpoint := fmt.Sprintf("/cloud/project/%s/database/mongodb/%s/user/%s",
		url.PathEscape(serviceName),
		url.PathEscape(clusterID),
		url.PathEscape(id),
	)

	return diag.FromErr(
		retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete),
			func() *retry.RetryError {
				log.Printf("[DEBUG] Will delete user %s from cluster %s from project %s", id, clusterID, serviceName)
				err := config.OVHClient.DeleteWithContext(ctx, endpoint, nil)
				if err != nil {
					if errOvh, ok := err.(*ovh.APIError); ok && (errOvh.Code == 409) {
						return retry.RetryableError(err)
					}
					err = helpers.CheckDeleted(d, err, endpoint)
					if err != nil {
						return retry.NonRetryableError(err)
					}
					return nil
				}

				log.Printf("[DEBUG] Waiting for user %s to be DELETED", id)
				err = waitForCloudProjectDatabaseUserDeleted(ctx, config.OVHClient, serviceName, "mongodb", clusterID, id, d.Timeout(schema.TimeoutDelete))
				if err != nil {
					return retry.NonRetryableError(fmt.Errorf("timeout while waiting user %s to be DELETED: %w", id, err))
				}
				log.Printf("[DEBUG] user %s is DELETED", id)

				d.SetId("")

				return nil
			},
		),
	)
}
