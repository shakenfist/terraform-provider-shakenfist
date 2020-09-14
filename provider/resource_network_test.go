package provider

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	client "github.com/shakenfist/client-go"
)

func TestAccShakenFistNetwork(t *testing.T) {
	var network client.Network

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resType := "shakenfist_network."
	resName := "external"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNetwork1(randomName),

				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkExists(resType+resName, &network),
					testAccNetworkValues(&network, resName, randomName),

					resource.TestCheckResourceAttr(
						resType+resName, "name", "testacc-"+randomName+"-external"),
					resource.TestCheckResourceAttrSet(resType+resName, "uuid"),
					resource.TestCheckResourceAttr(
						resType+resName, "netblock", "10.0.1.0/24"),
					resource.TestCheckResourceAttr(
						resType+resName, "provide_dhcp", "true"),
					resource.TestCheckResourceAttr(
						resType+resName, "provide_nat", "false"),

					testAccNetworkMetadata(resType+resName, map[string]string{
						"purpose": "external",
					}),
				),
			},
			{
				// Change the network configuration
				Config: testAccResourceNetwork2(randomName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resType+resName, "name", "testacc-"+randomName+"-external"),
					resource.TestCheckResourceAttrSet(resType+resName, "uuid"),
					resource.TestCheckResourceAttr(
						resType+resName, "netblock", "10.0.99.0/24"),
					resource.TestCheckResourceAttr(
						resType+resName, "provide_dhcp", "true"),
					resource.TestCheckResourceAttr(
						resType+resName, "provide_nat", "true"),
				),
			},
		},
	})
}

func testAccResourceNetwork1(randomName string) string {
	res := `
	resource "shakenfist_network" "external" {
		name = "testacc-{name}-external"
		netblock = "10.0.1.0/24"
		provide_dhcp = true
		provide_nat = false
		metadata = {
			purpose = "external"
		}
	}`

	r := strings.NewReplacer("{name}", randomName)
	return r.Replace(res)
}

// testAccResourceNetwork2 has a different subnet and NAT is now true.
func testAccResourceNetwork2(randomName string) string {
	res := `
	resource "shakenfist_network" "external" {
		name = "testacc-{name}-external"
		netblock = "10.0.99.0/24"
		provide_dhcp = true
		provide_nat = true
		metadata = {
			purpose = "external"
		}
	}`

	r := strings.NewReplacer("{name}", randomName)
	return r.Replace(res)
}

func testAccNetworkValues(actual *client.Network,
	tfName, randomName string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		correctName := "testacc-" + randomName + "-" + tfName
		if actual.Name != correctName {
			return fmt.Errorf("Network name incorrect=%s (should be %s)",
				actual.Name, correctName)
		}

		// Find the corresponding state object
		n := "shakenfist_network." + tfName
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// Test Terraform parameters against the actual network
		tf := rs.Primary.Attributes

		if actual.UUID != tf["uuid"] {
			return fmt.Errorf("Incorrect netblock: %s != %s",
				actual.NetBlock, tf["netblock"])
		}

		if actual.NetBlock != tf["netblock"] {
			return fmt.Errorf("Incorrect netblock: %s != %s",
				actual.NetBlock, tf["netblock"])
		}

		if strconv.FormatBool(actual.ProvideDHCP) != tf["provide_dhcp"] {
			return fmt.Errorf("Incorrect provide_dhcp: %t != %s",
				actual.ProvideDHCP, tf["provide_dhcp"])
		}

		if strconv.FormatBool(actual.ProvideNAT) != tf["provide_nat"] {
			return fmt.Errorf("Incorrect provide_nat: %t != %s",
				actual.ProvideNAT, tf["netblock"])
		}

		return nil
	}
}

func testAccNetworkMetadata(
	n string, correctMeta map[string]string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		// Find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// Retrieve the configured instance from the test setup
		apiClient := testAccProvider.Meta().(*client.Client)
		serverMeta, err := apiClient.GetNetworkMetadata(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Instance (%s) metadata cannot be retrieved: %v",
				rs.Primary.ID, err)
		}

		return compareMetadata(correctMeta, serverMeta)
	}
}

func testAccCheckNetworkExists(n string, net *client.Network) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// Retrieve the configured network from the test setup
		apiClient := testAccProvider.Meta().(*client.Client)
		resp, err := apiClient.GetNetwork(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Network (%s) cannot be retrieved: %v", rs.Primary.ID, err)
		}

		*net = resp

		return nil
	}
}
