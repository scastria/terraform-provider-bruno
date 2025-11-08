package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/scastria/terraform-provider-bruno/bruno"
)

//const (
//	COLLECTION_META_TAG                 = "meta"
//	COLLECTION_AUTH_TAG                 = "auth"
//	COLLECTION_PRE_REQUEST_VARS_TAG     = "vars:pre-request"
//	COLLECTION_POST_RESPONSE_VARS_TAG   = "vars:post-response"
//	COLLECTION_PRE_REQUEST_SCRIPT_TAG   = "script:pre-request"
//	COLLECTION_POST_RESPONSE_SCRIPT_TAG = "script:post-response"
//)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: bruno.Provider,
	})
	//schema := map[string]string{
	//	COLLECTION_META_TAG:                 dsl.DICT_TAG,
	//	COLLECTION_AUTH_TAG:                 dsl.DICT_TAG,
	//	COLLECTION_PRE_REQUEST_VARS_TAG:     dsl.DICT_TAG,
	//	COLLECTION_POST_RESPONSE_VARS_TAG:   dsl.DICT_TAG,
	//	COLLECTION_PRE_REQUEST_SCRIPT_TAG:   dsl.TEXT_TAG,
	//	COLLECTION_POST_RESPONSE_SCRIPT_TAG: dsl.TEXT_TAG,
	//}
	//doc, err := dsl.ImportDoc("/Users/shawncastrianni/GIT/bruno/data_api_one-development/collection.bru", schema)
	//if err != nil {
	//	panic(err)
	//}
	//err = doc.ExportDoc("test.bru")
	//println(err)
}
