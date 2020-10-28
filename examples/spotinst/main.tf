terraform {
  required_version = ">= 0.13.0"
  required_providers {
    multiverse = {
      source = "github.com/mobfox/multiverse"
      version = ">=0.0.3"
    }
  }
}


resource "multiverse" "spotinst_targetset_and_rules" {
  executor = "python3"
  script = "spotinst_mlb_targetset.py"
  id_key = "id"
  config = jsonencode({
    "name": "test-terraform-test",
    "mlb_id": "lb-123",
    "mlb_deployment_id": "dp-123",
    "mlb_listener_ids": [
      "ls-123",
      "ls-456"],
    "test_group_callback_fqdn": "test.fqdn.com",
    "control_group_callback_fqdn": "control.fqdn.com"
    "testTargetSet": null
    "controlTargetSet": null
  })
}

output "resources" {
  value = "${multiverse.spotinst_targetset_and_rules.id}"
}

output "test_data" {
  value = "${multiverse.spotinst_targetset_and_rules.config}"
}

output "test_targetset_id" {
  value = "${jsondecode(multiverse.spotinst_targetset_and_rules.config)["testTargetSet"]}"
}

//output "control_targetset_id" {
//  value = "${multiverse.spotinst_targetset_and_rules.data["controlTargetSet"]}"
//}