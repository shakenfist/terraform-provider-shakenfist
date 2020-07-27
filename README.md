Terraform Provider for Shaken Fist
==================================
![Go](https://github.com/shakenfist/terraform-provider-shakenfist/workflows/Go/badge.svg)  ![golangci-lint](https://github.com/shakenfist/terraform-provider-shakenfist/workflows/golangci-lint/badge.svg)

What is this?
-------------
The Shaken Fist Terraform provider is a plugin for Terraform that allows for the full lifecycle management of resources on Shaken Fist.

[Shaken Fist](https://github.com/shakenfist/shakenfist) is a deliberately minimal cloud.

Shaken Fist Resources
---------
This provider supports all Shaken Fist resources:
* Namespaces
* Instances
* Networks
* Floating IP's

Terraform Configuration
-----------------------
Examples of complete configuration files are available in the ```examples``` directory.

### Provider
```
provider "shakenfist" {
    server_url = "http://sf-1:13000"
    namespace = "devtest"
    key = "longsecurekey"
}
```

### Namespaces
* Multiple keys in the same namespace can be set by defining multiple `shakenfist_key` resources.
* Arbitrary metadata can be set on a namespace.

```
resource "shakenfist_namespace" "testspace" {
    name = "testspace"
    metadata = {
        owner = "bob"
        arbitrary = "clouds are awesome"
    }
}

resource "shakenfist_key" "key1" {
    namespace = shakenfist_namespace.testspace.name
    keyname = "key1"
    key = "secret"
}
```

### Instances
* Memory is defined in MB
* Multiple disks can defined
* Disk size is defined in GB
* Arbitrary metadata can be set on a namespace.

```
resource "shakenfist_instance" "jumpbox" {
    name = "jumpbox"
    cpus = 1
    memory = 1024
    disk {
        size = 8
        base = "cirros"
        bus = "ide"
        type = "disk"
    }
    disk {
        size = 3
        bus = "ide"
        type = "disk"
    }
    networks = [
        "uuid=${shakenfist_network.external.id}",
        ]
    metadata = {
        user = "old man"
    }
}
```

### Networks
* Arbitrary metadata can be set on a namespace.

```
resource "shakenfist_network" "external" {
    name = "external"
    netblock = "10.0.1.0/24"
    provide_dhcp = true
    provide_nat = true
    metadata = {
        purpose = "jump-hosts"
    }
}
```

### Floating IP's
```
resource "shakenfist_float" "jump" {
    interface = shakenfist_instance.jumpbox.interfaces[0]
}
```

Testing
-------
Terraform Provider acceptance tests require a Shaken Fist cluster and will modify resources on that cluster.

The provider acceptance tests are run using Make:
```
make testacc
```

The standard Shaken Fist environment variables must be set for the acceptance tests to run successfully.
```
SHAKENFIST_URL="http://sf-1:13000"
SHAKENFIST_NAMESPACE="dev"
SHAKENFIST_KEY="longsecurekey"
```

To run tests on Namespace resources, the Shaken Fist system privilege is required. Therefore the system namespace must be set:
```
SHAKENFIST_NAMESPACE=system
SHAKENFIST_KEY=Ukoh5vie
```
