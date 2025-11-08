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
  pre_request_script = split("\n",
  <<-EOT
    console.log("This is a pre-request script");
  EOT
  )
}

resource "bruno_folder" "Folder1" {
  name = "Folder1"
  auth = "inherit"
}

resource "bruno_folder" "ChildFolder" {
  parent_folder_id = bruno_folder.Folder1.id
  name = "Child Folder"
  auth = "inherit"
  tests = split("\n",
  <<-EOT
    test("Collection test: Request status code is 200", function () {
      expect(res.getStatus()).to.equal(200);
    });

    test("Endpoint test: Data array is not empty", function () {
      expect(res.getBody().length).to.not.equal(0);
    });
  EOT
  )
}

resource "bruno_folder" "ProblemFolder" {
  parent_folder_id = bruno_folder.Folder1.id
  name = "/ids/countries/{id}"
}

resource "bruno_request" "Request" {
  folder_id = bruno_folder.ChildFolder.id
  name = "get"
  base_url = "https://httpbin.konghq.com/anything"
  query_param {
    key = "sample"
    value = "value2"
  }
  query_param {
    key = "sample"
    value = "value"
  }
  query_param {
    key = "sample"
    value = "value2"
    disabled = true
  }
}

