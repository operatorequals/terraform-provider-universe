terraform {
  required_version = ">= 0.13.0"
  required_providers {
    multiverse = {
      source = "github.com/mobfox/multiverse"
      version = "0.0.1"
    }
    alpha = {
      source = "github.com/mobfox/alpha"
      version = "0.0.1"
    }
  }
}

resource "multiverse" "h" {
  executor = "python3"
  script = "hello_world.py"
  id_key = "id"
  data = <<JSON
{
  "name": "test-terraform-test",
  "url" : "https://foo.bar.quux/",
  "created-by" : "Elvis Presley"
}
JSON
}

resource "alpha" "h" {
  executor = "python3"
  script = "hello_world.py"
  id_key = "id"
  data = <<JSON
{
  "name": "A strange alpha resource",
  "url" : "https://alpha.com/",
  "created-by" : "Winston Churchill",
  "subject" : "Painting"
}
JSON
}

output "resources" {
  value = "${multiverse.h.id}"
}
