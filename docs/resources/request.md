# Resource: bruno_request
Represents a request
## Example usage
```hcl
resource "bruno_folder" "Folder" {
  name = "My Folder"
}
resource "bruno_request" "example" {
  folder_id = bruno_folder.Folder.id
  name = "My Request"
  method = "get"
  base_url = "https://httpbin.konghq.com/anything"
  query_param {
    key = "p1"
    value = "v1"
  }
  header {
    key = "h1"
    value = "v1"
  }
}
```
## Argument Reference
* `name` - **(Required, String)** The name of the request.
* `type` - **(Optional, String)** The type of the request. Allowed values: `graphql`, `http`. Default: `http`.
* `folder_id` - **(Optional, ForceNew, String)** The parent folder id.
* `method` - **(Optional, String)** The method of the request. Allowed values: `connect`, `delete`, `get`, `head`, `options`, `patch`, `post`, `put`, `trace`. Default: `get`.
* `auth` - **(Optional, String)** The authentication for the request. Allowed values: `none`, `apikey`, `awsv4`, `basic`, `bearer`, `digest`, `inherit`, `ntlm`, `oauth2`, `wsse`. Default: `none`.
* `base_url` - **(Required, String)** The base url of the request (excluding query params).
* `body` - **(Optional, List of String)** The raw JSON body of the request.
* `query_param` - **(Optional, list{query_param})** Configuration block for a query_param.  Can be specified multiple times for each query_param (allows same key for a multi-value parameter).  Each block supports the fields documented below.
* `header` - **(Optional, list{header})** Configuration block for a header.  Can be specified multiple times for each header (allows same key for a multi-value header).  Each block supports the fields documented below.
* `tests` - **(Optional, List of String)** The tests to perform after the response.
## query_param
* `key` - **(Required, String)** The name of the query param.
* `value` - **(Optional, String)** The value of the query param.
* `disabled` - **(Optional, Boolean)** Whether the query param is enabled. Default: `false`
## header
* `key` - **(Required, String)** The name of the header.
* `value` - **(Optional, String)** The value of the header.
* `disabled` - **(Optional, Boolean)** Whether the header is enabled. Default: `false`
## Attribute Reference
* `id` - **(String)** Relative path to the `<name>.bru` file
## Import
Requests can be imported using a proper value of `id` as described above
