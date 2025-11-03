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
}
