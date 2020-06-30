provider "shakenfist" {
    address = "http://localhost"
    port = 13000
    namespace = "system"
    key = "Ukoh5vie"
}

resource "shakenfist_namespace" "testname" {
    name = "testname"
    keyname = "testkeyname"
    key = "shhsecrect"
}
