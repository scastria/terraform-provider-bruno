# Resource: bruno_collection
Represents a collection
## Example usage
```hcl
resource "bruno_collection" "example" {
  name = "My Collection"
  auth = "none"
  pre_request_var {
    key = "url_base"
    value = "https://httpbin.konghq.com"
  }
  pre_request_script = [
    "script line 1",
    "script line 2"
  ]
}
```
## Argument Reference
* `name` - **(Required, String)** The name of the collection.
* `auth` - **(Optional, String)** The default authentication for the collection. Allowed values: `none`, `apikey`, `awsv4`, `basic`, `bearer`, `digest`, `ntlm`, `oauth2`, `wsse`. Default: `none`.
* `pre_request_script` - **(Optional, List of String)** The JS script to run before the request.
* `post_response_script` - **(Optional, List of String)** The JS script to run after the response.
* `pre_request_var` - **(Optional, list{var})** Configuration block for a pre-request variable.  Can be specified multiple times for each var.  Each block supports the fields documented below.
* `post_response_var` - **(Optional, list{var})** Configuration block for a post-response variable.  Can be specified multiple times for each var.  Each block supports the fields documented below.
## var
* `key` - **(Required, String)** The name of the variable.
* `value` - **(Optional, String)** The value of the variable.
* `disabled` - **(Optional, Boolean)** Whether the variable is disabled. Default: `false`
## Attribute Reference
* `id` - **(String)** Relative path to the `collection.bru` file
## Import
Collections can be imported using a proper value of `id` as described above
