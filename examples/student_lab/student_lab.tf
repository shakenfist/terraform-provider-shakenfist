// Create a two instance lab in a new namespace.
//
// Two instances connected via internal network. Access via
// a second network allocated a floating IP.
//


// Create new namespace using Shaken Fist "system" privilege.
provider "shakenfist" {
    alias = "system"
    hostname = "sf-1"
    port = 13000
    namespace = "system"
    key = "Ukoh5vie"
}

resource "shakenfist_namespace" "lab123" {
    provider = shakenfist.system
    name = "lab123"
    metadata = {
        owner = "cloudy"
        buildnote = "Pre build student lab"
        student = "no allocation yet"
    }
}

resource "shakenfist_key" "key1" {
    provider = shakenfist.system
    namespace = shakenfist_namespace.lab123.name
    keyname = "student"
    key = "secretsadf32jkhsdf234dsf"
}

//
// Build lab resources in the new namespace.
//
provider "shakenfist" {
    hostname = "sf-1"
    port = 13000
    namespace = shakenfist_namespace.lab123.name
    key = shakenfist_key.key1.key
}

// Jump host
resource "shakenfist_instance" "jump" {
    name = "jump"
    cpus = 1
    memory = 1024
    disks = [
        "size=8,base=cirros,bus=ide,type=disk",
        ]
    networks = [
        "uuid=${shakenfist_network.external.id}",
        ]
    metadata = {
        person = "old man"
        action = "shakes fist"
    }
}

resource "shakenfist_float" "external" {
    interface = shakenfist_instance.jump.interfaces[0]
}

resource "shakenfist_network" "external" {
    name = "external"
    netblock = "10.0.1.0/24"
    provide_dhcp = true
    provide_nat = true
    metadata = {
        purpose = "external"
    }
}

// Target host
resource "shakenfist_instance" "target" {
    name = "target"
    cpus = 1
    memory = 1024
    disks = [
        "size=8,base=cirros,bus=ide,type=disk",
        ]
    networks = [
        "uuid=${shakenfist_network.internal.id}",
        ]
}

resource "shakenfist_network" "internal" {
    name = "internal"
    netblock = "10.0.2.0/24"
    provide_dhcp = true
    provide_nat = false
    metadata = {
        purpose = "internal"
    }
}
