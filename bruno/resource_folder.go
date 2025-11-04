package bruno

import (
	"context"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-bruno/bruno/client"
	"github.com/scastria/terraform-provider-bruno/bruno/client/dsl"
)

const (
	FOLDER_META_TAG  = "meta"
	FOLDER_AUTH_TAG  = "auth"
	FOLDER_TESTS_TAG = "tests"
)

func resourceFolder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFolderCreate,
		ReadContext:   resourceFolderRead,
		UpdateContext: resourceFolderUpdate,
		DeleteContext: resourceFolderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_folder_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "",
			},
			"auth": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "none",
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"apikey",
					"awsv4",
					"basic",
					"bearer",
					"digest",
					"inherit",
					"ntlm",
					"oauth2",
					"wsse",
				}, false),
			},
			"tests": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceFolderCreateOrUpdate(ctx context.Context, d *schema.ResourceData, m interface{}, isCreate bool) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	parentFolderId := d.Get("parent_folder_id").(string)
	name := d.Get("name").(string)
	bd := dsl.BruDoc{
		Data: []dsl.BruBlock{},
	}
	// meta
	meta := dsl.BruDict{
		Tag: FOLDER_META_TAG,
		Data: map[string]interface{}{
			"name": name,
		},
	}
	bd.Data = append(bd.Data, &meta)
	// auth
	mode, ok := d.GetOk("auth")
	if ok {
		auth := dsl.BruDict{
			Tag: FOLDER_AUTH_TAG,
			Data: map[string]interface{}{
				"mode": mode.(string),
			},
		}
		bd.Data = append(bd.Data, &auth)
	}
	// tests
	tests, ok := d.GetOk("tests")
	if ok {
		testsText := createTextBlockFromArray(FOLDER_TESTS_TAG, tests.([]interface{}))
		bd.Data = append(bd.Data, testsText)
	}
	absPath := c.GetAbsolutePath(path.Dir(parentFolderId), name, "folder.bru")
	err := bd.ExportDoc(absPath)
	if err != nil {
		if isCreate {
			d.SetId("")
		}
		return diag.FromErr(err)
	}
	relativePath, err := c.GetRelativePath(absPath)
	if err != nil {
		if isCreate {
			d.SetId("")
		}
		return diag.FromErr(err)
	}
	d.SetId(relativePath)
	return diags
}

func resourceFolderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceFolderCreateOrUpdate(ctx, d, m, true)
}

func resourceFolderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	folderSchema := map[string]string{
		FOLDER_META_TAG:  dsl.DICT_TAG,
		FOLDER_AUTH_TAG:  dsl.DICT_TAG,
		FOLDER_TESTS_TAG: dsl.TEXT_TAG,
	}
	doc, err := dsl.ImportDoc(c.GetAbsolutePath(d.Id()), folderSchema)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	// meta
	metaBlock, err := doc.GetBlock(FOLDER_META_TAG)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	metaDict := metaBlock.(*dsl.BruDict)
	name := metaDict.Data["name"].(string)
	d.Set("name", name)
	// auth
	authBlock, err := doc.GetBlock(FOLDER_AUTH_TAG)
	if err == nil {
		authDict := authBlock.(*dsl.BruDict)
		auth := authDict.Data["mode"].(string)
		d.Set("auth", auth)
	}
	// tests
	testsBlock, err := doc.GetBlock(FOLDER_TESTS_TAG)
	if err == nil {
		testsText := testsBlock.(*dsl.BruText)
		testsArr := createArrayFromTextBlock(testsText)
		d.Set("tests", testsArr)
	}
	return diags
}

func resourceFolderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceFolderCreateOrUpdate(ctx, d, m, false)
}

func resourceFolderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	err := os.Remove(c.GetAbsolutePath(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
