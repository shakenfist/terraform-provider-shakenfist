provider "shakenfist" {
    hostname = "sf-1"
    port = 13000
    namespace = "system"
    key = "Ukoh5vie"
}

resource "shakenfist_namespace" "testspace" {
    name = "testspace"
    metadata = {
        owner = "cloudy"
        buildnote = "clouds are awesome"
    }
}

resource "shakenfist_key" "key1" {
    namespace = shakenfist_namespace.testspace.name
    keyname = "testkey1"
    key = "secret"
}

resource "shakenfist_key" "key2" {
    namespace = shakenfist_namespace.testspace.name
    keyname = "testkey2"
    key = "ENeXqQb3QFvhbMFnby3UN6SsLw6dP8hDuGyZAt"
}
