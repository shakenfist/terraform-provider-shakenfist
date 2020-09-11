package provider

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccShakenFistInstanceNet(t *testing.T) {
	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resFloat := "shakenfist_float.jump"

	ipAddr := regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$`)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFloat(randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resFloat, "interface"),
					// Check the address is properly formated IPv4
					resource.TestMatchResourceAttr(resFloat, "ipv4", ipAddr),
				),
			},
		},
	})
}

func testAccResourceFloat(randomName string) string {
	res := `
	resource "shakenfist_instance" "jump" {
		name = "testacc-{name}-jump"
		cpus = 1
		memory = 1024
		disk {
			size = 8
			base = "cirros"
			bus = "ide"
			type = "disk"
		}
		video {
			model = "cirrus"
			memory = 16384
		}
		network {
			network_uuid = shakenfist_network.external.id
		}
		metadata = {
			person = "old man"
			action = "shakes fist"
		}
	}

	resource "shakenfist_network" "external" {
		name = "testacc-{name}-external"
		netblock = "10.0.1.0/24"
		provide_dhcp = true
		provide_nat = false
		metadata = {
			purpose = "external"
		}
	}

	resource "shakenfist_float" "jump" {
		interface = shakenfist_instance.jump.network[0].interface_uuid
	}`

	r := strings.NewReplacer("{name}", randomName)
	return r.Replace(res)
}
