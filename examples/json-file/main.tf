terraform {
  required_version = ">= 0.13.0"
  required_providers {
    multiverse = {
      source = "github.com/mobfox/multiverse"
      version = ">=0.0.1"
    }
  }
}
provider "multiverse" {
  executor = "python3"
  script = "json-file.py"
  id_key = "filename"
  environment = {
    api_token = "redacted"
    // example environment
    servername = "api.example.com"
    api_token = "redacted"
  }
  computed = jsonencode([
    "created"])
}

resource "json-file" "h" {
  provider = multiverse
  // because Terraform does not scan local providers for resource types.
  executor = "python3"
  script = "json-file.py"
  id_key = "filename"
  computed = jsonencode([
    "created"])
  config = jsonencode({
    "name": "Don't Step On My Blue Suede Shoes",
    "created-by" : "Elvis Presley",
    "where" : "Gracelands"
    "hit" : "Gold"

  })
}

resource "json-file" "hp" {
  provider = multiverse
  // because Terraform does not scan local providers for resource types.
  config = jsonencode({
    "name": "Another strange resource",
    "main-character" : "Harry Potter",
    "nemesis" : "Tom Riddle",
    "likes" : [
      "Ginny Weasley",
      "Ron Weasley"
    ]
  })
}

resource "json-file" "i" {
  provider = multiverse
  // because Terraform does not scan local providers for resource types.
  config = jsonencode({
    "name": "Fake strange resource"
  })
}

output "hp_name" {
  value = jsondecode(json-file.hp.config)["name"]
}

output "hp_created" {
  value = jsondecode(json-file.hp.dynamic)["created"]
}
