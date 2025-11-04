terraform {
  required_providers {
    bruno = {
      source = "github.com/scastria/bruno"
    }
  }
}

provider "bruno" {
  collection_path = "/Users/shawncastrianni/GIT/bruno/terraform"
}

resource "bruno_collection" "Collection" {
  name = "test"
  # auth = "apikey"
  pre_request_var {
    key = "hi"
    value = "there"
    disabled = true
  }
  pre_request_var {
    key = "hi2"
    value = "there2"
  }
  # post_response_var {
  #   key = "hip"
  # }
}
