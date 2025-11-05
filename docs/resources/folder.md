# Resource: bruno_folder
Represents a folder
## Example usage
```hcl
resource "bruno_folder" "example" {
  name = "My Folder"
  tests = [
    "tests line 1",
    "tests line 2"
  ]
}
```
## Argument Reference
* `name` - **(Required, String)** The name of the folder.
* `parent_folder_id` - **(Optional, ForceNew, String)** The parent folder id.
* `auth` - **(Optional, String)** The default authentication for the folder. Allowed values: `none`, `apikey`, `awsv4`, `basic`, `bearer`, `digest`, `inherit`, `ntlm`, `oauth2`, `wsse`. Default: `none`.
* `pre_request_var` - **(Optional, list{var})** Configuration block for a pre-request variable.  Can be specified multiple times for each var.  Each block supports the fields documented below.
* `post_response_var` - **(Optional, list{var})** Configuration block for a post-response variable.  Can be specified multiple times for each var.  Each block supports the fields documented below.
* `tests` - **(Optional, List of String)** The default tests to perform after the response.
## var
* `key` - **(Required, String)** The name of the variable.
* `value` - **(Optional, String)** The value of the variable.
* `disabled` - **(Optional, Boolean)** Whether the variable is disabled. Default: `false`
## Attribute Reference
* `id` - **(String)** Relative path to the `folder.bru` file
## Import
Folders can be imported using a proper value of `id` as described above
