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

Here an example of Spot.io Multi Load Balancer TargetSet

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

Provider configuration data is passed in the environment variables:

* `id` - if not `create` this is the ID of the resource as returned to Terraform in the create
* `script` -
* `id_key` - 
* `executor` -
* `script` -

Plus any attributes present in the `environment` section in the provider block.

Your script could look something like the hello-world example below. This script maintains files in the file system 
containing JSON data in the HCL. 

```python
import os
import sys
import json
import tempfile

if __name__ == '__main__':
   event = sys.argv[1]
   id = os.environ.get("id")
   script = os.environ.get("script")
   input = sys.stdin.read()
   input_dict = json.loads(input)
   result = None
   if event == "create":
        ff = tempfile.NamedTemporaryFile(mode = 'w+',  prefix=script, delete=False)
        ff.write(json.dumps(input_dict))
        ff.close()
        input_dict.update({ "id" : ff.name})
        result = input_dict
   elif event == "read":
        fr=open(id, mode='r+')
        data = fr.read()
        fr.close()
        if len(data) > 0:
            result = json.loads(data)
        else:
            result = {}
   elif event == "update":
       fu=open(id,mode='w+')
       fu.write(json.dumps(input_dict))
       fu.close()
       result = input_dict
   elif event == "delete":
       os.remove(id)
       result =  {}
   print(json.dumps(result))


```

To test your script before using in TF

```bash
echo "{\"key\":\"value\"}" | python3 my_resource.py create
```
## Renaming the Resource Type

In your Terraform source code you may not want to see the resource type `multiverse`. You might a 
better name, reflecting the actual resouce type you're managing. So you might want this instead:

```hcl-terraform
resource "spot_io_elastic_instance" "myapp" {
  executor = "python3"
  script = "spotinst_mlb_targetset.py"
  id_key = "id"
  data = <<JSON
        {
         . . .
        }
JSON
}
```

This can be achieved by copying or linking to the provider binary file with a name inclusive of the new resource type name:

```shell script
 # Move to the plugins directory wherein lies the provider
cd ~/.terraform.d/plugins/github.com/mobfox/alpha/0.0.1/linux_amd64
# Copy the original file
cp terraform-provider-multiverse  terraform-provider-spot_io_elastic_instance
# or maybe link it
ln -s terraform-provider-multiverse  terraform-provider-spot_io_elastic_instance
```

Then you need to configure the provider in your TF file:

```hcl-terraform
terraform {
  required_version = ">= 0.13.0"
  required_providers {
    multiverse = {
      source = "github.com/mobfox/multiverse"
      version = ">=0.0.1"
    }
    spot_io_elastic_instance = {
      source = "github.com/mobfox/spot_io_elastic_instance"
      version = ">=0.0.1"
    }
  }
}
```
How des this work? The provider tells Terraform that it supports a single resource type. It extracts the name of the resource type 
from it's own executable in the plugins directory.

You could potentially use this to 'fake out' a normal provider to investigate its behaviour or 
emulate a defunct provider.


## Configuring the Provider

Terraform allows [configuration of providers](https://www.terraform.io/docs/configuration/providers.html#provider-configuration-1), 
in a `'provider` clause. The multiverse provider also has configuration where you specify the default executor, script and id fields.  
An additional field `environment` contains a map of environment variables which are passed to the script when it is executed. 

This means you don't need to repeat the `executor` nad `script` each time the provider is used.  You can 
override the defaults in the resource block as below.

```hcl-terraform
provider "alpha" {
  environment = {
    servername = "api.example.com"
    api_token = "redacted"
  }
  executor = "python3"
  script = "hello_world.py"
  id_key = "id"
}

resource "alpha" "h1" {
  data = <<JSON
    {
      "name": "test-terraform-test-1",
    }
JSON
}

resource "alpha" "h2" {
  script = "hello_world_v2.py"
  data = <<JSON
    {
      "name": "test-terraform-test-2",
    }
JSON
}

```
 

## Developing the Provider


If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.15.2+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

 > A good IDE is always beneficial. The kindly folk at [JetBrains](https://www.jetbrains.com/) provide Open Source authors with a free licenses to their excellent [Goland](https://www.jetbrains.com/go/) product, a cross-platform IDE built specially for Go developers   

To compile the provider, run `make build`. This will build the provider and put the provider binary in the workspace directory.

```sh script
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
