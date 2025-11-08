package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/scastria/terraform-provider-bruno/bruno"
)

//const (
//	//COLLECTION_META_TAG                 = "meta"
//	//COLLECTION_AUTH_TAG                 = "auth"
//	//COLLECTION_PRE_REQUEST_VARS_TAG     = "vars:pre-request"
//	//COLLECTION_POST_RESPONSE_VARS_TAG   = "vars:post-response"
//	//COLLECTION_PRE_REQUEST_SCRIPT_TAG   = "script:pre-request"
//	//COLLECTION_POST_RESPONSE_SCRIPT_TAG = "script:post-response"
//	REQUEST_META_TAG         = "meta"
//	REQUEST_TESTS_TAG        = "tests"
//	REQUEST_CONNECT_TAG      = "connect"
//	REQUEST_DELETE_TAG       = "delete"
//	REQUEST_GET_TAG          = "get"
//	REQUEST_HEAD_TAG         = "head"
//	REQUEST_OPTIONS_TAG      = "options"
//	REQUEST_PATCH_TAG        = "patch"
//	REQUEST_POST_TAG         = "post"
//	REQUEST_PUT_TAG          = "put"
//	REQUEST_TRACE_TAG        = "trace"
//	REQUEST_QUERY_PARAMS_TAG = "params:query"
//	REQUEST_HEADERS_TAG      = "headers"
//	REQUEST_JSON_BODY_TAG    = "body:json"
//)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: bruno.Provider,
	})
	//schema := map[string]string{
	//	//COLLECTION_META_TAG:                 dsl.DICT_TAG,
	//	//COLLECTION_AUTH_TAG:                 dsl.DICT_TAG,
	//	//COLLECTION_PRE_REQUEST_VARS_TAG:     dsl.DICT_TAG,
	//	//COLLECTION_POST_RESPONSE_VARS_TAG:   dsl.DICT_TAG,
	//	//COLLECTION_PRE_REQUEST_SCRIPT_TAG:   dsl.TEXT_TAG,
	//	//COLLECTION_POST_RESPONSE_SCRIPT_TAG: dsl.TEXT_TAG,
	//	REQUEST_META_TAG:         dsl.DICT_TAG,
	//	REQUEST_TESTS_TAG:        dsl.TEXT_TAG,
	//	REQUEST_CONNECT_TAG:      dsl.DICT_TAG,
	//	REQUEST_DELETE_TAG:       dsl.DICT_TAG,
	//	REQUEST_GET_TAG:          dsl.DICT_TAG,
	//	REQUEST_HEAD_TAG:         dsl.DICT_TAG,
	//	REQUEST_OPTIONS_TAG:      dsl.DICT_TAG,
	//	REQUEST_PATCH_TAG:        dsl.DICT_TAG,
	//	REQUEST_POST_TAG:         dsl.DICT_TAG,
	//	REQUEST_PUT_TAG:          dsl.DICT_TAG,
	//	REQUEST_TRACE_TAG:        dsl.DICT_TAG,
	//	REQUEST_QUERY_PARAMS_TAG: dsl.DICT_TAG,
	//	REQUEST_JSON_BODY_TAG:    dsl.TEXT_TAG,
	//	REQUEST_HEADERS_TAG:      dsl.DICT_TAG,
	//}
	//doc, err := dsl.ImportDoc("/Users/shawncastrianni/GIT/bruno/data_api_one-development/Sector APIs/sectors-periodicals/get.bru", schema)
	//if err != nil {
	//	panic(err)
	//}
	//err = doc.ExportDoc("test.bru")
	//println(err)
}
