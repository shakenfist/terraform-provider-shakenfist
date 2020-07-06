package provider

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/shakenfist/client-go"
)

func TestAccShakenFistInstance(t *testing.T) {
	var instance client.Instance

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resType := "shakenfist_instance."
	resName1 := "jump"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceInstance1(randomName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resType+resName1, &instance),
					testAccInstanceValues(&instance, resName1, randomName),

					resource.TestCheckResourceAttrSet(resType+resName1, "uuid"),
					resource.TestCheckResourceAttr(resType+resName1, "name",
						"testacc-"+randomName+"-jump"),
					resource.TestCheckResourceAttr(
						resType+resName1, "cpus", "1"),
					resource.TestCheckResourceAttr(
						resType+resName1, "memory", "1024"),
					resource.TestCheckResourceAttrSet(
						resType+resName1, "node"),
					resource.TestCheckResourceAttrSet(
						resType+resName1, "console_port"),
					resource.TestCheckResourceAttrSet(
						resType+resName1, "vdi_port"),
					resource.TestCheckResourceAttrSet(
						resType+resName1, "state"),

					testAccInstanceMetadata(resType+resName1, map[string]string{
						"person": "old man",
						"action": "shakes fist",
					}),
				),
			},
			{
				Config: testAccResourceInstance2(randomName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resType+resName1, &instance),
					testAccInstanceValues(&instance, resName1, randomName),

					resource.TestCheckResourceAttr(
						resType+resName1, "cpus", "2"),
					resource.TestCheckResourceAttr(
						resType+resName1, "memory", "2048"),

					testAccInstanceMetadata(resType+resName1, map[string]string{
						"person": "old man",
						"action": "screams into rock",
					}),
				),
			},
		},
	})
}

func testAccResourceInstance1(randomName string) string {
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
		networks = [
			"uuid=${shakenfist_network.external.id}",
			]
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
	}`

	r := strings.NewReplacer("{name}", randomName)
	return r.Replace(res)
}

func testAccResourceInstance2(randomName string) string {
	res := `
	resource "shakenfist_instance" "jump" {
		name = "testacc-{name}-jump"
		cpus = 2
		memory = 2048
		disk {
			size = 8
			base = "cirros"
			bus = "ide"
			type = "disk"
		}
		disk {
			size = 1
			bus = "ide"
			type = "disk"
		}
		networks = [
			"uuid=${shakenfist_network.external.id}",
			]
		metadata = {
			person = "old man"
			action = "screams into rock"
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
	}`

	r := strings.NewReplacer("{name}", randomName)
	return r.Replace(res)
}

func testAccInstanceValues(actual *client.Instance,
	tfName, randomName string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		correctName := "testacc-" + randomName + "-" + tfName
		if actual.Name != correctName {
			return fmt.Errorf("Instance name incorrect=%s (should be %s)",
				actual.Name, correctName)
		}

		// Find the corresponding state object
		n := "shakenfist_instance." + tfName
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// Test Terraform parameters against the actual instance
		tf := rs.Primary.Attributes

		if strconv.Itoa(actual.CPUs) != tf["cpus"] {
			return fmt.Errorf("Incorrect number of CPU's: %d != %s",
				actual.CPUs, tf["cpus"])
		}

		if strconv.Itoa(actual.Memory) != tf["memory"] {
			return fmt.Errorf("Incorrect memory: %d != %s",
				actual.Memory, tf["memory"])
		}

		return nil
	}
}

func testAccInstanceMetadata(
	n string, correctMeta map[string]string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		// Find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// Retrieve the configured instance from the test setup
		apiClient := testAccProvider.Meta().(*client.Client)
		serverMeta, err := apiClient.GetInstanceMetadata(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Instance (%s) metadata cannot be retrieved: %v",
				rs.Primary.ID, err)
		}

		return compareMetadata(correctMeta, serverMeta)
	}
}

func testAccCheckInstanceExists(n string, instance *client.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// Retrieve the configured instance from the test setup
		apiClient := testAccProvider.Meta().(*client.Client)
		resp, err := apiClient.GetInstance(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Instance (%s) cannot be retrieved: %v", rs.Primary.ID, err)
		}

		*instance = resp

		return nil
	}
}
