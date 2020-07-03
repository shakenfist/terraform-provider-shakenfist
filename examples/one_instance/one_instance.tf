provider "shakenfist" {
    hostname = "sf-1"
    port = 13000
    namespace = "testspace"
    key = "secret"
}

resource "shakenfist_instance" "sftest" {
    name = "sftest"
    cpus = 1
    memory = 1024
    disks = [
        "size=8,base=cirros,bus=ide,type=disk",
        ]
    networks = [
        "uuid=${shakenfist_network.sf-net-1.id}",
        ]
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
    interface = shakenfist_instance.sftest.interfaces[0]
}
