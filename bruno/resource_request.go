package bruno

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-bruno/bruno/client"
	"github.com/scastria/terraform-provider-bruno/bruno/client/dsl"
)

const (
	REQUEST_META_TAG         = "meta"
	REQUEST_TESTS_TAG        = "tests"
	REQUEST_CONNECT_TAG      = "connect"
	REQUEST_DELETE_TAG       = "delete"
	REQUEST_GET_TAG          = "get"
	REQUEST_HEAD_TAG         = "head"
	REQUEST_OPTIONS_TAG      = "options"
	REQUEST_PATCH_TAG        = "patch"
	REQUEST_POST_TAG         = "post"
	REQUEST_PUT_TAG          = "put"
	REQUEST_TRACE_TAG        = "trace"
	REQUEST_QUERY_PARAMS_TAG = "params:query"
	REQUEST_HEADERS_TAG      = "headers"
	REQUEST_JSON_BODY_TAG    = "body:json"
)

func resourceRequest() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRequestCreate,
		ReadContext:   resourceRequestRead,
		UpdateContext: resourceRequestUpdate,
		DeleteContext: resourceRequestDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "http",
				ValidateFunc: validation.StringInSlice([]string{
					"graphql",
					"http",
				}, false),
			},
			"folder_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "",
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  REQUEST_GET_TAG,
				ValidateFunc: validation.StringInSlice([]string{
					REQUEST_CONNECT_TAG,
					REQUEST_DELETE_TAG,
					REQUEST_GET_TAG,
					REQUEST_HEAD_TAG,
					REQUEST_OPTIONS_TAG,
					REQUEST_PATCH_TAG,
					REQUEST_POST_TAG,
					REQUEST_PUT_TAG,
					REQUEST_TRACE_TAG,
				}, false),
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
			"body": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"base_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query_param": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"header": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
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

func resourceRequestCreateOrUpdate(ctx context.Context, d *schema.ResourceData, m interface{}, isCreate bool) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	folderId := d.Get("folder_id").(string)
	name := d.Get("name").(string)
	requestType := d.Get("type").(string)
	method := d.Get("method").(string)
	baseUrl := d.Get("base_url").(string)
	bd := dsl.BruDoc{
		Data: []dsl.BruBlock{},
	}
	// meta
	meta := dsl.BruDict{
		Tag: REQUEST_META_TAG,
		Data: map[string]interface{}{
			"name": name,
			"type": requestType,
		},
	}
	bd.Data = append(bd.Data, &meta)
	// method
	methodDict := dsl.BruDict{
		Tag: method,
		Data: map[string]interface{}{
			//"url":  baseUrl,
			"body": "none",
		},
	}
	bd.Data = append(bd.Data, &methodDict)
	// auth
	auth, ok := d.GetOk("auth")
	if ok {
		methodDict.Data["auth"] = auth.(string)
	}
	// query_param
	queryParams, ok := d.GetOk("query_param")
	requestQuery := url.Values{}
	if ok {
		queryParamDict := createVariableDictBlockFromMap(REQUEST_QUERY_PARAMS_TAG, queryParams.(*schema.Set))
		for k, v := range queryParamDict.Data {
			if strings.HasPrefix(k, dsl.DISABLED_PREFIX) {
				continue
			}
			requestQuery.Add(k, fmt.Sprintf("%v", v))
		}
		bd.Data = append(bd.Data, queryParamDict)
	}
	encodedQuery := requestQuery.Encode()
	// Manually append query params to prevent URL parsing from escaping variables in the host and path
	if len(encodedQuery) > 0 {
		methodDict.Data["url"] = baseUrl + "?" + encodedQuery
	} else {
		methodDict.Data["url"] = baseUrl
	}
	// header
	headers, ok := d.GetOk("header")
	if ok {
		headerDict := createVariableDictBlockFromMap(REQUEST_HEADERS_TAG, headers.(*schema.Set))
		bd.Data = append(bd.Data, headerDict)
	}
	// body
	body, ok := d.GetOk("body")
	if ok {
		methodDict.Data["body"] = "json"
		bodyText := createTextBlockFromArray(REQUEST_JSON_BODY_TAG, body.([]interface{}))
		bd.Data = append(bd.Data, bodyText)
	}
	// tests
	tests, ok := d.GetOk("tests")
	if ok {
		testsText := createTextBlockFromArray(REQUEST_TESTS_TAG, tests.([]interface{}))
		bd.Data = append(bd.Data, testsText)
	}
	absPath := c.GetAbsolutePath(path.Dir(folderId), fmt.Sprintf("%s.bru", prepareFolderName(name)))
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

func resourceRequestCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceRequestCreateOrUpdate(ctx, d, m, true)
}

func resourceRequestRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestSchema := map[string]string{
		REQUEST_META_TAG:         dsl.DICT_TAG,
		REQUEST_TESTS_TAG:        dsl.TEXT_TAG,
		REQUEST_CONNECT_TAG:      dsl.DICT_TAG,
		REQUEST_DELETE_TAG:       dsl.DICT_TAG,
		REQUEST_GET_TAG:          dsl.DICT_TAG,
		REQUEST_HEAD_TAG:         dsl.DICT_TAG,
		REQUEST_OPTIONS_TAG:      dsl.DICT_TAG,
		REQUEST_PATCH_TAG:        dsl.DICT_TAG,
		REQUEST_POST_TAG:         dsl.DICT_TAG,
		REQUEST_PUT_TAG:          dsl.DICT_TAG,
		REQUEST_TRACE_TAG:        dsl.DICT_TAG,
		REQUEST_QUERY_PARAMS_TAG: dsl.DICT_TAG,
		REQUEST_JSON_BODY_TAG:    dsl.TEXT_TAG,
		REQUEST_HEADERS_TAG:      dsl.DICT_TAG,
	}
	doc, err := dsl.ImportDoc(c.GetAbsolutePath(d.Id()), requestSchema)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	// meta
	metaBlock, err := doc.GetBlock(REQUEST_META_TAG)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	metaDict := metaBlock.(*dsl.BruDict)
	name := metaDict.Data["name"].(string)
	d.Set("name", name)
	requestType := metaDict.Data["type"].(string)
	d.Set("type", requestType)
	// method
	// Have to check all possible methods
	methods := []string{
		REQUEST_CONNECT_TAG,
		REQUEST_DELETE_TAG,
		REQUEST_GET_TAG,
		REQUEST_HEAD_TAG,
		REQUEST_OPTIONS_TAG,
		REQUEST_PATCH_TAG,
		REQUEST_POST_TAG,
		REQUEST_PUT_TAG,
		REQUEST_TRACE_TAG,
	}
	var methodBlock dsl.BruBlock
	methodBlock = nil
	for _, method := range methods {
		methodBlock, err = doc.GetBlock(method)
		if err == nil {
			d.Set("method", method)
			break
		}
	}
	if methodBlock != nil {
		methodDict := methodBlock.(*dsl.BruDict)
		// auth
		auth, exists := methodDict.Data["auth"]
		if exists {
			d.Set("auth", auth)
		}
		// base_url
		url, exists := methodDict.Data["url"]
		if exists {
			// Manually strip query params to prevent URL parsing from escaping variables in the host and path
			d.Set("base_url", strings.Split(url.(string), "?")[0])
		}
	}
	// query_param
	var queryParamMap []map[string]interface{}
	queryParamMap = nil
	queryParamBlock, err := doc.GetBlock(REQUEST_QUERY_PARAMS_TAG)
	if err == nil {
		queryParamDict := queryParamBlock.(*dsl.BruDict)
		queryParamMap = createMapFromVariableDictBlock(queryParamDict)
	}
	d.Set("query_param", queryParamMap)
	// header
	var headerMap []map[string]interface{}
	headerMap = nil
	headerBlock, err := doc.GetBlock(REQUEST_HEADERS_TAG)
	if err == nil {
		headerDict := headerBlock.(*dsl.BruDict)
		headerMap = createMapFromVariableDictBlock(headerDict)
	}
	d.Set("header", headerMap)
	// body
	var bodyArr []string
	bodyArr = nil
	bodyBlock, err := doc.GetBlock(REQUEST_JSON_BODY_TAG)
	if err == nil {
		bodyText := bodyBlock.(*dsl.BruText)
		bodyArr = createArrayFromTextBlock(bodyText)
	}
	d.Set("body", bodyArr)
	// tests
	var testsArr []string
	testsArr = nil
	testsBlock, err := doc.GetBlock(REQUEST_TESTS_TAG)
	if err == nil {
		testsText := testsBlock.(*dsl.BruText)
		testsArr = createArrayFromTextBlock(testsText)
	}
	d.Set("tests", testsArr)
	// folder_id
	folderDir := path.Dir(d.Id())
	if folderDir != "." {
		folderId := path.Join(folderDir, "folder.bru")
		d.Set("folder_id", folderId)
	}
	return diags
}

func resourceRequestUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceRequestCreateOrUpdate(ctx, d, m, false)
}

func resourceRequestDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	err := os.Remove(c.GetAbsolutePath(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
