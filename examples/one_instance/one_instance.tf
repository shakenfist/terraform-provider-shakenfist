terraform {
    required_providers {
        shakenfist = {
            source = "shakenfist/shakenfist"
            versions = ">=0.3"
        }
    }
}

provider "shakenfist" {
    server_url = "http://sf-1:13000"
    namespace = "testspace"
    key = "secret"
}

resource "shakenfist_instance" "sftest" {
    name = "sftest"
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
    video {
        model = "cirrus"
        memory = 16384
    }
    network {
        network_uuid = shakenfist_network.sf-net-1.id
    }
    network {
        network_uuid = shakenfist_network.sf-net-1.id
        ipv4 = "10.0.1.17"
        model = "e1000"
    }
    network {
        network_uuid = shakenfist_network.sf-net-1.id
        mac = "12:34:56:78:9a:Bc"
    }
    metadata = {
        person = "old man"
        action = "shakes fist"
    }
}

resource "shakenfist_network" "sf-net-1" {
    name = "sf-net-1"
    netblock = "10.0.1.0/24"
    provide_dhcp = true
    provide_nat = true
    metadata = {
        purpose = "external"
    }
}

resource "shakenfist_float" "sf-float-1" {
    interface = shakenfist_instance.sftest.network[0].interface_uuid
}
