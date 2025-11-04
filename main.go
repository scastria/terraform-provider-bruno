package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/scastria/terraform-provider-bruno/bruno"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: bruno.Provider,
	})
	//schema := map[string]string{
	//	"meta": "dict",
	//	"arr":  "array",
	//}
	//doc, err := dsl.ImportDoc("/Users/shawncastrianni/GIT/bruno/terraform/collection.bru", schema)
	//if err != nil {
	//	panic(err)
	//}
	//exported := doc.Export()
	//println(exported)
}
