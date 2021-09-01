//
// This example needs environment variables to specify resource types:
//
//   export TERRAFORM_UNIVERSE_RESOURCETYPES='json_file'
//   export TERRAFORM_LINUX_RESOURCETYPES='json_file'
//
terraform {
  required_version = ">= 0.13.0"
  required_providers {
    jsonfile = {
      source = "github.com/birchb1024/jsonfile"
      version = ">=0.0.5"
    }
    jsonfile2 = {
      source = "github.com/birchb1024/jsonfile2"
      version = ">=0.0.5"
    }
  }
}
provider "jsonfile" {
  executor = "python3"
  script = "json_file.py"
  id_key = "filename"
}

provider "jsonfile2" {
  executor = "python3"
  script = "json_file2.py"
  id_key = "filename"
}

resource "jsonfile" "f" {
  config = jsonencode({
    "name": "Don't Step On My Blue Suede Shoes",
    "created-by" : "Elvis Presley",
    "where" : "Gracelands"
    "hit" : "Gold"
    "@created": null
  })
}

resource "jsonfile2" "f2" {
  config = jsonencode({
    "name": "Billie Jean",
    "created-by" : "Michael Jackson",
    "where" : "US"
    "hit" : "Gold"
    "@created": null
  })
}

output "f_filename" {
    value = jsonfile.f.id
}

output "f2_filename" {
    value = jsonfile2.f2.id
}
