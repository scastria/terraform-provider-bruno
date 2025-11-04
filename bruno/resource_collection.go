package bruno

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-bruno/bruno/client"
	"github.com/scastria/terraform-provider-bruno/bruno/client/dsl"
)

const (
	COLLECTION_META_TAG                 = "meta"
	COLLECTION_AUTH_TAG                 = "auth"
	COLLECTION_PRE_REQUEST_VARS_TAG     = "vars:pre-request"
	COLLECTION_POST_RESPONSE_VARS_TAG   = "vars:post-response"
	COLLECTION_PRE_REQUEST_SCRIPT_TAG   = "script:pre-request"
	COLLECTION_POST_RESPONSE_SCRIPT_TAG = "script:post-response"
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
				ForceNew: true,
				Required: true,
			},
			"auth": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "none",
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"apikey",
					"awsv4",
					"basic",
					"bearer",
					"digest",
					"ntlm",
					"oauth2",
					"wsse",
				}, false),
			},
			"pre_request_var": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							ForceNew: true,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"post_response_var": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							ForceNew: true,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"pre_request_script": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"post_response_script": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func createVariableDictBlockFromMap(tag string, variables *schema.Set) *dsl.BruDict {
	retVal := dsl.BruDict{
		Tag:  tag,
		Data: make(map[string]interface{}),
	}
	for _, variable := range variables.List() {
		variableMap := variable.(map[string]interface{})
		varPrefix := ""
		if variableMap["disabled"].(bool) {
			varPrefix = dsl.DISABLED_PREFIX
		}
		retVal.Data[fmt.Sprintf("%s%s", varPrefix, variableMap["key"].(string))] = variableMap["value"].(string)
	}
	return &retVal
}

func createTextBlockFromArray(tag string, arr []interface{}) *dsl.BruText {
	retVal := dsl.BruText{
		Tag: tag,
	}
	lines := convertInterfaceArrayToStringArray(arr)
	retVal.Data = strings.Join(lines, "\n")
	return &retVal
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
	bd := dsl.BruDoc{
		Data: []dsl.BruBlock{},
	}
	// meta
	meta := dsl.BruDict{
		Tag: COLLECTION_META_TAG,
		Data: map[string]interface{}{
			"name": d.Get("name").(string),
		},
	}
	bd.Data = append(bd.Data, &meta)
	// auth
	mode, ok := d.GetOk("auth")
	if ok {
		auth := dsl.BruDict{
			Tag: COLLECTION_AUTH_TAG,
			Data: map[string]interface{}{
				"mode": mode.(string),
			},
		}
		bd.Data = append(bd.Data, &auth)
	}
	// pre_request_var
	preRequestVariables, ok := d.GetOk("pre_request_var")
	if ok {
		preRequestVarDict := createVariableDictBlockFromMap(COLLECTION_PRE_REQUEST_VARS_TAG, preRequestVariables.(*schema.Set))
		bd.Data = append(bd.Data, preRequestVarDict)
	}
	// post_response_var
	postResponseVariables, ok := d.GetOk("post_response_var")
	if ok {
		postResponseVarDict := createVariableDictBlockFromMap(COLLECTION_POST_RESPONSE_VARS_TAG, postResponseVariables.(*schema.Set))
		bd.Data = append(bd.Data, postResponseVarDict)
	}
	// pre_request_script
	preRequestScript, ok := d.GetOk("pre_request_script")
	if ok {
		preRequestScriptText := createTextBlockFromArray(COLLECTION_PRE_REQUEST_SCRIPT_TAG, preRequestScript.([]interface{}))
		bd.Data = append(bd.Data, preRequestScriptText)
	}
	// post_response_script
	postResponseScript, ok := d.GetOk("post_response_script")
	if ok {
		postResponseScriptText := createTextBlockFromArray(COLLECTION_POST_RESPONSE_SCRIPT_TAG, postResponseScript.([]interface{}))
		bd.Data = append(bd.Data, postResponseScriptText)
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

func createMapFromVariableDictBlock(variableDict *dsl.BruDict) []map[string]interface{} {
	var retVal []map[string]interface{}
	for k, v := range variableDict.Data {
		variableMap := map[string]interface{}{}
		variableMap["key"] = strings.TrimPrefix(k, dsl.DISABLED_PREFIX)
		variableMap["value"] = v
		variableMap["disabled"] = strings.HasPrefix(k, dsl.DISABLED_PREFIX)
		retVal = append(retVal, variableMap)
	}
	return retVal
}

func createArrayFromTextBlock(textBlock *dsl.BruText) []string {
	return strings.Split(textBlock.Data, "\n")
}

func resourceCollectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	collectionSchema := map[string]string{
		COLLECTION_META_TAG:                 dsl.DICT_TAG,
		COLLECTION_AUTH_TAG:                 dsl.DICT_TAG,
		COLLECTION_PRE_REQUEST_VARS_TAG:     dsl.DICT_TAG,
		COLLECTION_POST_RESPONSE_VARS_TAG:   dsl.DICT_TAG,
		COLLECTION_PRE_REQUEST_SCRIPT_TAG:   dsl.TEXT_TAG,
		COLLECTION_POST_RESPONSE_SCRIPT_TAG: dsl.TEXT_TAG,
	}
	doc, err := dsl.ImportDoc(c.GetAbsolutePath(d.Id()), collectionSchema)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	// meta
	metaBlock, err := doc.GetBlock(COLLECTION_META_TAG)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	metaDict := metaBlock.(*dsl.BruDict)
	name := metaDict.Data["name"].(string)
	d.Set("name", name)
	// auth
	authBlock, err := doc.GetBlock(COLLECTION_AUTH_TAG)
	if err == nil {
		authDict := authBlock.(*dsl.BruDict)
		auth := authDict.Data["mode"].(string)
		d.Set("auth", auth)
	}
	// pre_request_var
	preRequestVarBlock, err := doc.GetBlock(COLLECTION_PRE_REQUEST_VARS_TAG)
	if err == nil {
		preRequestVarDict := preRequestVarBlock.(*dsl.BruDict)
		preRequestVariableMap := createMapFromVariableDictBlock(preRequestVarDict)
		d.Set("pre_request_var", preRequestVariableMap)
	}
	// post_response_var
	postResponseVarBlock, err := doc.GetBlock(COLLECTION_POST_RESPONSE_VARS_TAG)
	if err == nil {
		postResponseVarDict := postResponseVarBlock.(*dsl.BruDict)
		postResponseVariableMap := createMapFromVariableDictBlock(postResponseVarDict)
		d.Set("post_response_var", postResponseVariableMap)
	}
	// pre_request_script
	preRequestScriptBlock, err := doc.GetBlock(COLLECTION_PRE_REQUEST_SCRIPT_TAG)
	if err == nil {
		preRequestScriptText := preRequestScriptBlock.(*dsl.BruText)
		preRequestScriptArr := createArrayFromTextBlock(preRequestScriptText)
		d.Set("pre_request_script", preRequestScriptArr)
	}
	// post_response_script
	postResponseScriptBlock, err := doc.GetBlock(COLLECTION_POST_RESPONSE_SCRIPT_TAG)
	if err == nil {
		postResponseScriptText := postResponseScriptBlock.(*dsl.BruText)
		postResponseScriptArr := createArrayFromTextBlock(postResponseScriptText)
		d.Set("post_response_script", postResponseScriptArr)
	}
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
