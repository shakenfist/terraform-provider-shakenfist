module github.com/shakenfist/terraform-provider-shakenfist

go 1.14

require (
	github.com/hashicorp/terraform v0.12.26
	github.com/shakenfist/client-go v0.0.0-20200622074543-ea3d1911e584
)

// REMOVE!!
replace github.com/shakenfist/client-go => ../client-go
