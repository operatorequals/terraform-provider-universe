Terraform Provider Multiverse
==================

You can use this provider instead of writing your own Terraform Custom Provider in the Go language. Just write your logic in any language you prefer (python, node, java, shell) and use it with this provider. You can write a script that will be used to create, update or destroy external resources that are not supported by Terraform providers.

Maintainers
-----------

The MobFox DevOps team at [MobFox](https://www.mobfox.com/) maintains this provider. 

This is the birchb1024 fork, maintained by Peter Birch.

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.13
-	[Go](https://golang.org/doc/install) 1.15.2 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/birchb1024/terraform-provider-multiverse`

```sh
$ mkdir -p $GOPATH/src/github.com/mobfox; cd $GOPATH/src/github.com/mobfox
$ git clone git@github.com:birchb1024/terraform-provider-multiverse.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/mobfox/terraform-provider-multiverse
$ make build
```

Using the provider
----------------------

Check the `examples/` directory

Here an example of Spotinst Multai Load Balancer TargetSet

```hcl
provider "multiverse" {}

resource "multiverse_custom_resource" "spotinst_targetset_and_rules" {
  executor = "python3"
  script = "spotinst_mlb_targetset.py"
  id_key = "id"
  data = <<JSON
{
  "name": "test-terraform-test",
  "mlb_id": "lb-123",
  "mlb_deployment_id": "dp-123",
  "mlb_listener_ids": ["ls-123", "ls-456"],
  "test_group_callback_fqdn": "test.fqdn.com",
  "control_group_callback_fqdn": "control.fqdn.com"
}
JSON
}
```

- When you run `terraform apply` the resource will be created / updated
- When you run `terraform destroy` the resource will be destroyed

#### Attributes

##### Input

* `executor (string)` could be anything like python, bash, sh, node, java, awscli ... etc
* `script (string)` the path to your script or program to run, the script must exit with code 0 and return a valid json string
* `id_key (string)` the key of returned result to be used as id by terraform
* `data (string)` must be a valid json string

##### Output
* `resource (map[string])` the output of your script must be a valid json with all keys of type *string* in the form `{"key":"value"}`


#### Referencing in TF template

This an example how to reference the resource and access its attributes

Let's say your script returned the following result
 
```json
{
  "id": "my-123",
  "name": "my-resource",
  "capacity": "20g"
}
```

then the resource in TF will have these attributes

```hcl
id = "my-123"
resource = {
  name = "my-resource"
  capacity = "20g"
}
```

you can access these attributes using variables

```
${multiverse_custom_resource.my_custom_resource.id} # accessing id
${multiverse_custom_resource.my_custom_resource.resource["name"]}
${multiverse_custom_resource.my_custom_resource.resource["capacity"]}
```

#### Why the attribute *data* is stringified JSON?

This will give you flexibility in passing your arguments with mixed types. We couldn't define a with generic mixed types, 
if we used map then all attributes have to be explicitly defined in the schema or all its attributes have the same type.


Writing a script
-------------------

Your script must be able to handle the TF event and the JSON payload *data*

* `event` : will have one of these `create, read, delete, update`
* `data` : is passed via `stdin`

Provider configuration data is passed in the environment variable `multiverse`.

Your script should look something like the hello-world example:

```python
import os
import sys
import json

if __name__ == '__main__':
   event = sys.argv[1]
   config = os.environ.get("multiverse")
   if len(sys.argv) > 2:
       config = sys.argv[2]
   input = sys.stdin.read()
   input_dict = json.loads(input)
   if event == "create":
        input_dict.update({ "id" : "1"})
        result = input_dict
   elif event == "read":
        result =  input_dict
   elif event == "update":
        result =  input_dict
   elif event == "delete":
        result =  {}
   else:
       sys.exit(1)
   print(json.dumps(result))
```

To test your script before using in TF

```bash
echo "{\"key\":\"value\"}" | python3 my_resource.py create
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.15.2+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

 > A good IDE is always beneficial. The kindly folk at [JetBrains](https://www.jetbrains.com/) provide Open Source authors with a free licenses to their excellent [Goland](https://www.jetbrains.com/go/) product, a cross-platform IDE built specially for Go developers   

To compile the provider, run `make build`. This will build the provider and put the provider binary in the workspace directory.

```sh
$ make build
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

To install the provider in the usual places for the `terraform` to use run make install. It will place it the plugin directories:

```
$HOME/.terraform.d/
└── plugins
    ├── github.com
    │   └── mobfox
    │       ├── alpha
    │       │   └── 0.0.1
    │       │       └── linux_amd64
    │       │           └── terraform-provider-alpha
    │       └── multiverse
    │           └── 0.0.1
    │               └── linux_amd64
    │                   ├── terraform-provider-alpha
    │                   └── terraform-provider-multiverse
    ├── terraform-provider-alpha
    └── terraform-provider-multiverse

```

Feel free to contribute!
