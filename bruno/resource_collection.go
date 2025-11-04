package bruno

import (
	"context"
	"encoding/json"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-bruno/bruno/client"
	"github.com/scastria/terraform-provider-bruno/bruno/client/dsl"
)

func resourceCollection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCollectionCreate,
		ReadContext:   resourceCollectionRead,
		DeleteContext: resourceCollectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCollectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	body := map[string]interface{}{
		"version": "1",
		"name":    d.Get("name").(string),
		"type":    "collection",
		"ignore": []string{
			"node_modules",
			".git",
		},
	}
	bytes, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	err = os.WriteFile(c.GetAbsolutePath("bruno.json"), bytes, 0644)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	meta := dsl.BruDict{
		Tag: "meta",
		Data: map[string]interface{}{
			"name": d.Get("name").(string),
		},
	}
	bd := dsl.BruDoc{
		Data: []dsl.BruBlock{
			&meta,
		},
	}
	relativePath := "collection.bru"
	err = bd.ExportDoc(c.GetAbsolutePath(relativePath))
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(relativePath)
	return diags
}

func resourceCollectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	collectionSchema := map[string]string{
		"meta": dsl.DICT_TAG,
	}
	doc, err := dsl.ImportDoc(c.GetAbsolutePath(d.Id()), collectionSchema)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	metaBlock, err := doc.GetBlock("meta")
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	metaDict := metaBlock.(*dsl.BruDict)
	name := metaDict.Data["name"].(string)
	d.Set("name", name)
	return diags
}

func resourceCollectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	err := os.Remove(c.GetAbsolutePath("bruno.json"))
	if err != nil {
		return diag.FromErr(err)
	}
	err = os.Remove(c.GetAbsolutePath(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
