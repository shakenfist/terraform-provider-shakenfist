package provider

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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

					testAccInstanceVideo(resType+resName1, client.VideoSpec{
						Model:  "cirrus",
						Memory: 16384,
					}),

					testAccInstanceNetwork(resType+resName1,
						[]client.NetworkSpec{
							{
								MACAddress: "12:34:56:78:9a:bc",
								Model:      "e1000",
								Address:    "10.0.2.100",
							},
						}),

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

					testAccInterfaceOrder(resType+resName1, []string{
						"external",
						"second",
						"third",
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
		video {
			model = "cirrus"
			memory = 16384
		}
		network {
			network_uuid = shakenfist_network.external.id
			mac = "12:34:56:78:9a:bc"
			model = "e1000"
			ipv4 = "10.0.2.100"
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
		video {
			model = "cirrus"
			memory = 16384
		}
		network {
			network_uuid = shakenfist_network.external.id
		}
		network {
			network_uuid = shakenfist_network.second.id
		}
		network {
			network_uuid = shakenfist_network.third.id
		}
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
	}

	resource "shakenfist_network" "second" {
		name = "testacc-{name}-second"
		netblock = "10.0.2.0/24"
		provide_dhcp = true
		provide_nat = false
		metadata = {
			purpose = "second"
		}
	}

	resource "shakenfist_network" "third" {
		name = "testacc-{name}-third"
		netblock = "10.0.3.0/24"
		provide_dhcp = true
		provide_nat = false
		metadata = {
			purpose = "third net"
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

func testAccInstanceVideo(
	n string, correctVideo client.VideoSpec) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		// Find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// Retrieve the configured instance from the test setup
		apiClient := testAccProvider.Meta().(*client.Client)
		inst, err := apiClient.GetInstance(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Instance (%s) cannot be retrieved: %v",
				rs.Primary.ID, err)
		}

		if inst.Video != correctVideo {
			return fmt.Errorf("Instance video is %v but should be %v",
				inst.Video, correctVideo)
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

func testAccInterfaceOrder(n string,
	netOrder []string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		// Find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// Retrieve the instance interfaces from the test setup
		apiClient := testAccProvider.Meta().(*client.Client)
		interfaces, err := apiClient.GetInstanceInterfaces(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf(
				"Instance (%s) interfaces cannot be retrieved: %v",
				rs.Primary.ID, err)
		}

		// Build interface lookup table
		iOrder := map[int]string{}
		for _, i := range interfaces {
			iOrder[i.Order] = i.NetworkUUID
		}

		// Test the expected Terraform network UUID's
		for count, net := range netOrder {

			// Get the Terraform network object
			rs, ok := s.RootModule().Resources["shakenfist_network."+net]
			if !ok {
				return fmt.Errorf("Network not found: %s", net)
			}

			// Check the network UUID is set on the interface
			if rs.Primary.ID != iOrder[count] {
				return fmt.Errorf(
					"Network (%s) is not on correct interface %d, have %s",
					net, count, iOrder[count])
			}
		}

		return nil
	}
}

func testAccInstanceNetwork(n string,
	correctNets []client.NetworkSpec) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		// Find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// Retrieve the configured instance from the test setup
		apiClient := testAccProvider.Meta().(*client.Client)
		interfaces, err := apiClient.GetInstanceInterfaces(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Instance (%s) cannot be retrieved: %v",
				rs.Primary.ID, err)
		}

		countActual := len(interfaces)
		countCorrect := len(correctNets)
		if countActual != countCorrect {
			return fmt.Errorf("Got %d interfaces, expected %d interfaces",
				countActual, countCorrect)
		}

		for i, actual := range interfaces {
			if actual.IPv4 != correctNets[i].Address {
				return fmt.Errorf("IP address is %s but should be %s",
					actual.IPv4, correctNets[i].Address)
			}
			if actual.MACAddress != correctNets[i].MACAddress {
				return fmt.Errorf("MAC address is %s but should be %s",
					actual.MACAddress, correctNets[i].MACAddress)
			}
			if actual.Model != correctNets[i].Model {
				return fmt.Errorf("Model is %s but should be %s",
					actual.Model, correctNets[i].Model)
			}
		}

		return nil
	}
}
