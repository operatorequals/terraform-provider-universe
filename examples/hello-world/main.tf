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
  configuration = "\"Hello Enola\""
}

provider "alpha" {
  configuration = "\"Hello Lucy\""
}

resource "multiverse" "h" {
  executor = "python3"
  script = "hello_world.py"
  id_key = "id"
  data = <<JSON
{
  "name": "test-terraform-test",
  "created-by" : "Elvis Presley",
  "where" : "gracelands"
}
JSON
}


resource "alpha" "hp" {
  executor = "python3"
  script = "hello_world.py"
  id_key = "id"
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
  executor = "python3"
  script = "hello_world.py"
  id_key = "id"
  data = <<JSON
{
  "name": "Fake strange alpha resource"
}
JSON
}

output "data" {
  value = "${alpha.hp.data}"
}
