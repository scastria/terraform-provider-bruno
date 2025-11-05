package bruno

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-bruno/bruno/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"collection_path": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"bruno_collection": resourceCollection(),
			"bruno_folder":     resourceFolder(),
			"bruno_request":    resourceRequest(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	collectionPath := d.Get("collection_path").(string)

	var diags diag.Diagnostics
	c, err := client.NewClient(collectionPath)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diags
}
