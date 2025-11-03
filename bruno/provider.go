package bruno

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"collection_path": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap:         map[string]*schema.Resource{},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	//host := d.Get("host").(string)
	//port := d.Get("port").(int)
	//defaultDatabase := d.Get("default_database").(string)
	//username := d.Get("username").(string)
	//password := d.Get("password").(string)
	//
	var diags diag.Diagnostics
	//c, err := client.NewClient(host, port, defaultDatabase, username, password)
	//if err != nil {
	//	return nil, diag.FromErr(err)
	//}
	return nil, diags
}
