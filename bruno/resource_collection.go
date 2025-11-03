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
		UpdateContext: resourceCollectionUpdate,
		DeleteContext: resourceCollectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
	filePath := c.GetPath("bruno.json")
	err = os.WriteFile(filePath, bytes, 0644)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	filePath = c.GetPath("collection.bru")
	meta := dsl.BruDict{
		Tag: "meta",
		Data: map[string]interface{}{
			"name": d.Get("name").(string),
		},
	}
	bd := dsl.BruDoc{
		Data: []dsl.BruBlock{
			meta,
		},
	}
	err = os.WriteFile(filePath, []byte(bd.Export()), 0644)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(filePath)
	return diags
}

func resourceCollectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	//workspaceId, id := client.CollectionDecodeId(d.Id())
	//c := m.(*client.Client)
	//requestPath := fmt.Sprintf(client.CollectionPathGet, id)
	//body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	//if err != nil {
	//	d.SetId("")
	//	re := err.(*client.RequestError)
	//	if re.StatusCode == http.StatusNotFound {
	//		return diags
	//	}
	//	return diag.FromErr(err)
	//}
	//retVal := &client.CollectionContainer{}
	//err = json.NewDecoder(body).Decode(retVal)
	//if err != nil {
	//	d.SetId("")
	//	return diag.FromErr(err)
	//}
	//retVal.Child.Info.WorkspaceId = workspaceId
	//fillResourceDataFromCollection(retVal, d)
	return diags
}

func resourceCollectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	//workspaceId, id := client.CollectionDecodeId(d.Id())
	//c := m.(*client.Client)
	//buf := bytes.Buffer{}
	//upCollection := client.CollectionUpdateContainer{}
	//fillCollectionUpdate(&upCollection, d)
	//err := json.NewEncoder(&buf).Encode(upCollection)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//requestPath := fmt.Sprintf(client.CollectionPathGet, id)
	//requestHeaders := http.Header{
	//	headers.ContentType: []string{client.ApplicationJson},
	//}
	//_, err = c.HttpRequest(ctx, http.MethodPatch, requestPath, nil, requestHeaders, &buf)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//// Must re-read the collection to get the full response
	//requestPath = fmt.Sprintf(client.CollectionPathGet, id)
	//body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//retVal := &client.CollectionContainer{}
	//err = json.NewDecoder(body).Decode(retVal)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//retVal.Child.Info.WorkspaceId = workspaceId
	//fillResourceDataFromCollection(retVal, d)
	return diags
}

func resourceCollectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	//_, id := client.CollectionDecodeId(d.Id())
	//c := m.(*client.Client)
	//requestPath := fmt.Sprintf(client.CollectionPathGet, id)
	//_, err := c.HttpRequest(ctx, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//d.SetId("")
	return diags
}
