terraform {
  required_version = ">= 0.13.0"
  required_providers {
    multiverse = {
      source = "github.com/mobfox/multiverse"
      version = ">=0.0.1"
    }
    alpha = {
      source = "github.com/mobfox/alpha"
      version = ">=0.0.1"
    }
  }
}
provider "multiverse" {
  environment = {
    api_token = "redacted"
  }
}

provider "alpha" {
  environment = {
    servername = "api.example.com"
    api_token = "redacted"
  }
  executor = "python3"
  script = "hello_world.py"
  id_key = "id"
}

resource "multiverse" "h" {
  executor = "python3"
  script = "hello_world.py"
  id_key = "id"
  data = <<JSON
    {
      "name": "test-terraform-test-43",
      "created-by" : "Elvis Presley",
      "where" : "gracelands"
    }
JSON
}


resource "alpha" "hp" {
  data = <<JSON
    {
      "name": "Another strange alpha resource",
      "created-by" : "Harry Potter",
      "nemesis" : "Tom Riddle",
      "likes" : [
        "Ginny Weasley",
        "Ron Weasley"
      ]
    }
JSON
}

resource "alpha" "i" {

  data = <<JSON
{
  "name": "Fake strange alpha resource"
}
JSON
}

output "data" {
  value = "${alpha.hp.data}"
}
