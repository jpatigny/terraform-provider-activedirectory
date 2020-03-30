package activedirectory

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jpatigny/goPSRemoting"
)

func resourcecomputerObject() *schema.Resource {
	return &schema.Resource{
		Create: resourcecomputerObjectCreate,
		Read:   resourcecomputerObjectRead,
		Delete: resourcecomputerObjectDelete,
		Update: resourcecomputerObjectCreate,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"distname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourcecomputerObjectCreate(d *schema.ResourceData, m interface{}) error {
	//convert the interface so we can use the variables like username, etc
	client := m.(*ActiveDirectoryClient)

	name := d.Get("name").(string)
	distname := d.Get("distname").(string)

	var id string = name
	var psCommand string = "$object = New-ADComputer -Name \\\"" + name + "\\\" -SamAccountName \\\"" + name + "\\\" -Path \\\"" + distname + "\\\" -Confirm:$false"

	_, err := goPSRemoting.RunPowershellCommand(client.username, client.password, client.server, psCommand, client.usessl, client.usessh, client.authentication)
	if err != nil {
		//something bad happened
		return err
	}

	d.SetId(id)

	return nil
}

func resourcecomputerObjectRead(d *schema.ResourceData, m interface{}) error {
	//convert the interface so we can use the variables like username, etc
	client := m.(*ActiveDirectoryClient)

	name := d.Get("name").(string)

	var psCommand string = "$object = Get-ADComputer -Identity \\\"" + name + "\\\" }; if (!$object) { Write-Host 'TERRAFORM_NOT_FOUND' }"
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

	var id string = name
	d.SetId(id)
	return nil
}

func resourcecomputerObjectDelete(d *schema.ResourceData, m interface{}) error {
	//convert the interface so we can use the variables like username, etc
	client := m.(*ActiveDirectoryClient)

	name := d.Get("name").(string)

	var psCommand string = "$object = Remove-ADComputer -Identity \\\"" + name + "\\\" -Confirm:$false"
	_, err := goPSRemoting.RunPowershellCommand(client.username, client.password, client.server, psCommand, client.usessl, client.usessh, client.authentication)
	if err != nil {
		//something bad happened
		return err
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but it is added here for explicitness.
	d.SetId("")

	return nil
}
