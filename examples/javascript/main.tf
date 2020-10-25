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
  id_key = "id"
  computed = jsonencode([
    "created"])
  javascript = "echo.js"
}

resource "multiverse" "j" {
  config = jsonencode({
    "name": "Don't Step On My Blue Suede Shoes",
    "created-by" : "Elvis Presley",
    "where" : "Gracelands"
    "hit" : "Gold"
  })
}
output "j_name" {
  value = jsondecode(multiverse.j.config)["name"]
}
