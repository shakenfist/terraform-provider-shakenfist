// With many thanks to the example code from
// https://github.com/spaceapegames/terraform-provider-example
package provider

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	client "github.com/shakenfist/client-go"
)

func resourceInstance() *schema.Resource {
	fmt.Print()
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the instance",
				ForceNew:    true,
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the instance",
			},
			"cpus": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The number of CPUs for the instance",
				ForceNew:    true,
			},
			"memory": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The amount of RAM for the instance in GB",
				ForceNew:    true,
			},
			"disks": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"networks": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"interfaces": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ssh_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ssh key to embed into the instance via config drive",
			},
			"user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "User data to pass to the instance via config drive, encoded as base64",
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString},
			},
		},
		Create: resourceCreateInstance,
		Read:   resourceReadInstance,
		Delete: resourceDeleteInstance,
		Exists: resourceExistsInstance,
		Update: resourceUpdateInstance,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateInstance(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	var disks []client.DiskSpec
	var err error

	for _, disk := range d.Get("disks").([]interface{}) {
		var diskSpec client.DiskSpec
		for _, diskElem := range strings.Split(disk.(string), ",") {
			e := strings.Split(diskElem, "=")
			if e[0] == "base" {
				diskSpec.Base = e[1]
			} else if e[0] == "size" {
				diskSpec.Size, err = strconv.Atoi(e[1])
				if err != nil {
					return fmt.Errorf("Disk size is not numeric: %v", err)
				}
			} else if e[0] == "bus" {
				diskSpec.Bus = e[1]
			} else if e[0] == "type" {
				diskSpec.Type = e[1]
			} else {
				return fmt.Errorf("Incorrect disk spec, should be one of "+
					"\"base\", \"size\", \"bus\", \"type\": %s", e[0])
			}
		}
		disks = append(disks, diskSpec)
	}

	var networks []client.NetworkSpec
	for _, net := range d.Get("networks").([]interface{}) {
		var netSpec client.NetworkSpec
		for _, netElem := range strings.Split(net.(string), ",") {
			e := strings.Split(netElem, "=")
			if e[0] == "uuid" {
				netSpec.NetworkUUID = e[1]
			} else if e[0] == "address" {
				netSpec.Address = e[1]
			} else if e[0] == "macaddress" {
				netSpec.MACAddress = e[1]
			} else if e[0] == "model" {
				netSpec.Model = e[1]
			} else {
				return fmt.Errorf("Incorrect network type, should be one of "+
					"\"uuid\", \"address\", \"macaddress\", \"model\": %s", e[0])
			}
		}
		networks = append(networks, netSpec)
	}

	inst, err := apiClient.CreateInstance(d.Get("name").(string), d.Get("cpus").(int),
		d.Get("memory").(int), networks, disks, d.Get("ssh_key").(string),
		d.Get("user_data").(string))
	if err != nil {
		return fmt.Errorf("Unable to create instance: %v", err)
	}

	// If Shaken Fist has a server error, it can return a blank UUID
	if inst.UUID == "" {
		return fmt.Errorf("Shaken Fist has returned a null instance UUID")
	}

	// Signal to TF that the instance exists by setting the ID.
	d.Set("uuid", inst.UUID)
	d.SetId(inst.UUID)

	// Set metadata on the instance
	for k, v := range d.Get("metadata").(map[string]interface{}) {
		val, ok := v.(string)
		if ok != true {
			return fmt.Errorf("Tag value is not a string")
		}

		err := apiClient.SetMetadata(client.TypeInstance, inst.UUID, k, val)
		if err != nil {
			return fmt.Errorf("CreateInstance cannot store metadata: %v", err)
		}
	}

	return resourceReadInstance(d, m)
}

func resourceReadInstance(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	inst, err := apiClient.GetInstance(d.Id())
	if err != nil {
		return fmt.Errorf("Unable to retrieve instance: %v", err)
	}

	d.Set("uuid", inst.UUID)
	d.Set("name", inst.Name)
	d.Set("cpus", inst.CPUs)
	d.Set("memory", inst.Memory)
	d.Set("disks", inst.DiskSpecs)
	d.Set("ssh_key", inst.SSHKey)
	d.Set("node", inst.Node)
	d.Set("console_port", inst.ConsolePort)
	d.Set("vdi_port", inst.VDIPort)
	d.Set("user_data", inst.UserData)

	// Retrieve Interface UUID's
	uuid, err := getInterfaceUUIDs(apiClient, d.Id())
	if err != nil {
		return fmt.Errorf("ReadInstance error: %v", err)
	}
	d.Set("interfaces", uuid)

	// Retrieve metadata
	metadata, err := apiClient.GetMetadata(client.TypeInstance, inst.UUID)
	if err != nil {
		return fmt.Errorf("ReadInstance unable to retrieve metadata: %v", err)
	}
	d.Set("metadata", metadata)

	return nil
}

func resourceDeleteInstance(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	err := apiClient.DeleteInstance(d.Id())
	if err != nil {
		return fmt.Errorf("Unable to retrieve network: %v", err)
	}
	d.SetId("")
	return nil
}

func resourceExistsInstance(d *schema.ResourceData, m interface{}) (bool, error) {
	apiClient := m.(*client.Client)

	i, err := apiClient.GetInstance(d.Id())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		} else {
			return false, fmt.Errorf("Unable to check instance existence: %v", err)
		}
	}

	if i.State == "deleted" {
		return false, nil
	}

	return true, nil
}

func resourceUpdateInstance(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	if d.HasChange("metadata") {
		if err := updateMetadata(client.TypeInstance, d, apiClient); err != nil {
			return fmt.Errorf("UpdateInstance error: %v", err)
		}
	}

	return nil
}

func getInterfaceUUIDs(apiClient *client.Client, instanceUUID string) ([]string, error) {
	interfaces, err := apiClient.GetInstanceInterfaces(instanceUUID)
	if err != nil {
		return []string{}, fmt.Errorf(
			"unable to retrieve instance interfaces: %v", err)
	}

	var uuid []string
	for _, i := range interfaces {
		uuid = append(uuid, i.UUID)
	}

	return uuid, nil
}
