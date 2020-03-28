package activedirectory

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jpatigny/goPSRemoting"
	"errors"
	"strings"
)

func resourceGPOLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceGPOLinkCreate,
		Read:   resourceGPOLinkRead,
		Delete: resourceGPOLinkDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organizational_unit": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enable": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"enforce": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGPOLinkCreate(d *schema.ResourceData, m interface{}) error {
	//convert the interface so we can use the variables like username, etc
	client := m.(*ActiveDirectoryClient)

	name 				:= d.Get("name").(string)
	organizational_unit := d.Get("organizational_unit").(string)
	enable 				:= d.Get("enable").(string)
	enforce 			:= d.Get("enforce").(string)
	var psCommand string
	
	if name == "" {
		return errors.New("Must provide a GPO name.")
	}
	if organizational_unit == "" {
		return errors.New("Must provide an OU where to link GPO to")
	}

	psCommand = "New-GPLink -Name " + name + "-Target " + organizational_unit

    	_, err := goPSRemoting.RunPowershellCommand(client.username, client.password, client.server, psCommand, client.usessl, client.usessh)
	if err != nil {
		//something bad happened
		return err
	}

	var id string = name + "_" + organizational_unit + "_" + enable + "_" + enforce
 	d.SetId(id)
	return nil
}

func resourceGPOLinkRead(d *schema.ResourceData, m interface{}) error {
	//convert the interface so we can use the variables like username, etc
	client := m.(*ActiveDirectoryClient)

	name 				:= d.Get("name").(string)
	organizational_unit := d.Get("organizational_unit").(string)
	enable 				:= d.Get("enable").(string)
	enforce 			:= d.Get("enforce").(string)
	var psCommand string 

	//Get-DnsServerResourceRecord -ZoneName "contoso.com" -Name "Host03" -RRType "A"
	// var psCommand string = "try { $record = Get-DnsServerResourceRecord -ZoneName " + zone_name + " -RRType " + record_type + " -Name " + record_name + " -ErrorAction Stop } catch { $record = '''' }; if ($record) { write-host 'RECORD_FOUND' }"
	psCommand  = `
		$OU = Get-ADOrganizationalUnit -Filter * -Properties gPlink | ? {$_.Name -eq "`+ organizational_unit +`" }
		if($OU.Name) {
			$OUGPLinks = $OU.gPlink.split("][")
			$OUGPLinks =  @($OUGPLinks | ? {$_})
			foreach($GpLink in $OUGPLinks) {
				$GpName = [adsi]$GPlink.split(";")[0] | select -ExpandProperty displayName
				if ($GpName -eq "`+ name + `") {
					Write-host "GPOLINK_FOUND"
					break
				}
			}
		}
	`
		_, err := goPSRemoting.RunPowershellCommand(client.username, client.password, client.server, psCommand, client.usessl, client.usessh)
	if err != nil {
		if !strings.Contains(err.Error(), "ObjectNotFound") {
			//something bad happened
			return err
		} else {
			//not able to find the record - this is an error but ok
			d.SetId("")
			return nil
		}
	}

	var id string = name + "_" + organizational_unit + "_" + enable + "_" + enforce
 	d.SetId(id)

	return nil
}

func resourceGPOLinkDelete(d *schema.ResourceData, m interface{}) error {
	//convert the interface so we can use the variables like username, etc
	client := m.(*ActiveDirectoryClient)

	name 				:= d.Get("name").(string)
	organizational_unit := d.Get("organizational_unit").(string)
	var psCommand string
	
	if name == "" {
		return errors.New("Must provide a GPO name.")
	}
	if organizational_unit == "" {
		return errors.New("Must provide an OU where to link GPO to")
	}

	psCommand = "Remove-GPLink -Name " + name + "-Target " + organizational_unit + " -Confirm:$false -Force"

    	_, err := goPSRemoting.RunPowershellCommand(client.username, client.password, client.server, psCommand, client.usessl, client.usessh)
	if err != nil {
		//something bad happened
		return err
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but it is added here for explicitness.
	d.SetId("")

	return nil
}