provider "shakenfist" {
    address = "http://localhost"
    port = 3000
    namespace = "and"
    keyname = "testkey"
    key = "test"
}

resource "sf_instance" "sftest" {
    name = "sftest"
    cpus = 1
    memory = 1
    disks = [
        {
            size=8,
            base="cirros",
            bus="ide",
            type="disk"
            }

    ]
    networks = ["${sf_network.sf-net-1.uuid}"]
}

resource "sf_network" "sf-net-1" {
    name = "sf-net-1"
    netblock = "10.0.1.0/24"
    provide_dhcp = true
    provide_nat = true
}
