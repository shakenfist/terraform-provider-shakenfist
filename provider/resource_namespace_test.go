package provider

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	client "github.com/shakenfist/client-go"
)

// TestAccShakenFistNamespace tests the namespace and key creation.
//
// *** NOTE: Only tested when environment namespace set to 'system'.
//
func TestAccShakenFistNamespace(t *testing.T) {
	var keynames []string

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resType := "shakenfist_namespace."
	resName := "testspace"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				SkipFunc: testAccSystemKey,
				Config:   testAccResourceNamespace(randomName),
				Check: resource.ComposeTestCheckFunc(

					testAccCheckNamespaceExists(resType+resName, &keynames),
					testAccNamespaceKey(&keynames, "key1"),

					resource.TestCheckResourceAttr(
						resType+resName, "name", "testacc-"+randomName),

					testAccNamespaceMetadata(resType+resName, map[string]string{
						"owner":     "cloudy",
						"buildnote": "clouds are awesome",
					}),
				),
			},
		},
	})
}

func testAccResourceNamespace(randomName string) string {
	res := `
	resource "shakenfist_namespace" "testspace" {
		name = "testacc-{name}"
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
	`

	rName := strings.NewReplacer("{name}", randomName)

	return rName.Replace(res)
}

func testAccSystemKey() (bool, error) {
	namespace := os.Getenv("SHAKENFIST_NAMESPACE")

	if namespace != "system" {
		fmt.Printf("Skipping Namespace tests, namespace is not set to 'system'")
		return true, nil
	}

	return false, nil
}

func testAccNamespaceMetadata(
	n string, correctMeta map[string]string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		// Find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// Retrieve the configured instance from the test setup
		apiClient := testAccProvider.Meta().(*client.Client)
		serverMeta, err := apiClient.GetNamespaceMetadata(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Instance (%s) metadata cannot be retrieved: %v",
				rs.Primary.ID, err)
		}

		return compareMetadata(correctMeta, serverMeta)
	}
}

func testAccNamespaceKey(actual *[]string, tfName string) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		// Find the corresponding state object
		n := "shakenfist_key." + tfName
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// Test Terraform parameter against the actual keyname
		tf := rs.Primary.Attributes
		for _, k := range *actual {
			if k == tf["keyname"] {
				// TF keyname exists in the namespace on the server
				return nil
			}
		}

		return fmt.Errorf(
			"TF Keyname (%s) does not exist on server", tf["keyname"])
	}
}

func testAccCheckNamespaceExists(
	n string, keyNames *[]string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		// Find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// Retrieve the configured namespace from the test setup
		apiClient := testAccProvider.Meta().(*client.Client)
		names, err := apiClient.GetNamespaces()
		if err != nil {
			return fmt.Errorf("Namespaces cannot be retrieved: %v", err)
		}

		for _, k := range names {
			if k == rs.Primary.ID {
				// Namespace exists. Retrieve keynames.
				*keyNames, err = apiClient.GetNamespaceKeys(rs.Primary.ID)
				if err != nil {
					return fmt.Errorf("NamespaceExists cannot get keynames: %v",
						err)
				}

				return nil
			}
		}

		return fmt.Errorf("Namespace (%s) does not exist", rs.Primary.ID)
	}
}
