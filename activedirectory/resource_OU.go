package activedirectory

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jpatigny/goPSRemoting"
)

func resourceOU() *schema.Resource {
	return &schema.Resource{
		Create: resourceOUMappingCreate,
		Read:   resourceOUMappingRead,
		Delete: resourceOUMappingDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protected": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOUCreate(d *schema.ResourceData, m interface{}) error {
	//convert the interface so we can use the variables like username, etc
	client := m.(*ActiveDirectoryClient)

	name := d.Get("name").(string)
	path := d.Get("path").(string)
	protected := d.Get("protected").(bool)
	var id string = name + "_" + path + "_"
	var psCommand string

	fmt.Println("[Create] starting...")
	fmt.Println("[Create] check status of protected var")
	if protected {
		fmt.Println("[Create] protected var is set to true")
		psCommand = "New-ADOrganizationalUnit -Name \\\"" + name + `" -Path "` + path + "\\\" -ProtectedFromAccidentalDeletion $True"
		fmt.Println("%v, psCommand")
	} else {
		fmt.Println("[Create] protected var is set to true")
		psCommand = "New-ADOrganizationalUnit -Name \\\"" + name + `" -Path "` + path + "\\\" -ProtectedFromAccidentalDeletion $False"
		fmt.Println("%v, psCommand")
	}
	_, err := goPSRemoting.RunPowershellCommand(client.username, client.password, client.server, psCommand, client.usessl, client.usessh, client.authentication)
	if err != nil {
		//something bad happened
		return err
	}

	fmt.Println("OU successfully created")

	d.SetId(id)

	return nil
}

func resourceOURead(d *schema.ResourceData, m interface{}) error {
	//convert the interface so we can use the variables like username, etc
	client := m.(*ActiveDirectoryClient)

	name := d.Get("name").(string)
	path := d.Get("path").(string)
	//protected   := d.Get("protected").(string)

	//var psCommand string = "$object = Get-ADObject -SearchBase \\\"" + target_path + "\\\" -Filter {(name -eq \\\"" + object_name + "\\\") -AND (ObjectClass -eq \\\"" + object_class + "\\\")}; if (!$object) { Write-Host 'TERRAFORM_NOT_FOUND' }"
	var psCommand string = "$ou =  Get-ADOrganizationalUnit -Filter 'Name -like \\\"" + name + "\\\"; if (!$ou) { Write-host 'TERRAFORM_NOT_FOUND' }"

	stdout, err := goPSRemoting.RunPowershellCommand(client.username, client.password, client.server, psCommand, client.usessl, client.usessh, client.authentication)
	if err != nil {
		//something bad happened
		return err
	}

	if strings.Contains(stdout, "TERRAFORM_NOT_FOUND") {
		//not able to find the record - this is an error but ok
		d.SetId("")
		return nil
	}

	var id string = name + "_" + path
	d.Set("address", id)
	return nil
}

func resourceOUDelete(d *schema.ResourceData, m interface{}) error {
	//convert the interface so we can use the variables like username, etc
	client := m.(*ActiveDirectoryClient)

	name := d.Get("name").(string)
	path := d.Get("path").(string)

	var psCommand string = `Remove-ADOrganizationalUnit -Identity "OU=` + name + `,` + path + `" -Recursive	-Confirm:$False`
	_, err := goPSRemoting.RunPowershellCommand(client.username, client.password, client.server, psCommand, client.usessl, client.usessh, client.authentication)
	if err != nil {
		//something bad happened
		return err
	}

	fmt.Println("OU successfully created")
	// d.SetId("") is automatically called assuming delete returns no errors, but it is added here for explicitness.
	d.SetId("")

	return nil
}
