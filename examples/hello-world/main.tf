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
    api_token = "redacted" // example environment
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
  data = jsonencode({
      "name": "test-terraform-test-43",
      "created-by" : "Elvis Presley",
      "where" : "gracelands"
    })
}


resource "alpha" "hp" {
  data = jsonencode({

      "name": "Another strange alpha resource",
      "main-character" : "Harry Potter",
      "nemesis" : "Tom Riddle",
      "likes" : [
        "Ginny Weasley",
        "Ron Weasley"
      ]
    })
}

resource "alpha" "i" {

  data = jsonencode({
    "name": "Fake strange alpha resource"
  })
}

output "hp_name" {
  value = "${jsondecode(alpha.hp.data)["name"]}"
}
