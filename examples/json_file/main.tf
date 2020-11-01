//
// This example needs environment variables to specify resource types:
//
//   export TERRAFORM_UNIVERSE_RESOURCETYPES='json_file'
//   export TERRAFORM_LINUX_RESOURCETYPES='json_file'
//
terraform {
  required_version = ">= 0.13.0"
  required_providers {
    universe = {
      source = "github.com/birchb1024/universe"
      version = ">=0.0.4"
    }
    linux = {
      source = "github.com/birchb1024/linux"
      version = ">=0.0.4"
    }
  }
}
provider "universe" {
  executor = "python3"
  script = "json_file.py"
  id_key = "filename"
  environment = {
    api_token = "redacted"
    servername = "api.example.com"
    api_token = "redacted"
  }
}

resource "universe_json_file" "h" {
  config = jsonencode({
    "name": "Don't Step On My Blue Suede Shoes",
    "created-by" : "Elvis Presley",
    "where" : "Gracelands"
    "hit" : "Gold"
    "@created": null
  })
}

resource "universe_json_file" "hp" {
  config = jsonencode({
    "name": "Another strange resource",
    "main-character" : "Harry Potter",
    "nemesis" : "Tom Riddle",
    "likes" : [
      "Ginny Weasley",
      "Ron Weasley"
    ],
    "@created": 23
  })
}

resource "linux_json_file" "i" {
  executor = "python3"
  script = "json_file.py"
  id_key = "filename"
  config = jsonencode({
    "name": "Fake strange resource"
  })
}

output "hp_name" {
  value = jsondecode(universe_json_file.hp.config)["name"]
}

output "hp_created" {
  value = jsondecode(universe_json_file.hp.config)["@created"]
}
