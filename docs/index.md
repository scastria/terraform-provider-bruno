# Bruno Provider
The Bruno provider is used to control a local filesystem with `*.bru` files.

This provider does NOT cover 100% of Bruno's capabilities.  If there is something missing
that you would like to be added, please submit an Issue in corresponding GitHub repo.

## Example Usage
```hcl
terraform {
  required_providers {
    bruno = {
      source  = "scastria/bruno"
      version = "~> 0.1.0"
    }
  }
}

# Configure the Bruno Provider
provider "bruno" {
  collection_path = "XXXX"
}
```
## Argument Reference
* `collection_path` - **(Required, String)** Local filesystem path to the top level of the Bruno collection.
