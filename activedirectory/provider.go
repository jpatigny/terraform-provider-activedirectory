package activedirectory

import (
	"fmt"
	"log"
	
	logger "github.com/wayfair/terraform-provider-utils/log"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider allows making changes to Microsoft AD
// Utilises Powershell to connect to domain controller
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("USERNAME", nil),
				Description: "Username to connect to AD.",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PASSWORD", nil),
				Description: "The password to connect to AD.",
			},
			"server": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVER", nil),
				Description: "The AD server to connect to.",
			},
			"usessl": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("USESSL", false),
				Description: "Whether or not to use HTTPS to connect to WinRM",
			},
			"usessh": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("USESSH", false),
				Description: "Whether or not to use SSH to connect to WinRM",
			},
			"authentication": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTHENTICATION", nil),
				Description: "Authentication protocol for WinRM",
			},
			"default_computer_container": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DEFAULT_COMPUTER_CONTAINER", nil),
				Description: "The default computer container to move objects to on a delete - Defaults to '(Get-ADDomain | select computerscont*)'",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"activedirectory_ouMapping":       resourceOUMapping(),
			"activedirectory_groupMembership": resourcegroupMembership(),
			"activedirectory_ou":              resourceOU(),
			"activedirectory_computerObject":  resourcecomputerObject(),
		},

		ConfigureFunc: providerConfigure,
	}
}
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	username := d.Get("username").(string)
	if username == "" {
		return nil, fmt.Errorf("The 'username' property was not specified.")
	}

	usessh := d.Get("usessh").(string)

	password := d.Get("password").(string)
	if password == "" && usessh == "0" {
		return nil, fmt.Errorf("The 'password' property was not specified and usessh was false.")
	}

	server := d.Get("server").(string)
	if server == "" {
		return nil, fmt.Errorf("The 'server' property was not specified.")
	}

	usessl := d.Get("usessl").(string)
	authentication := d.Get("authentication").(string)
	default_computer_container := d.Get("default_computer_container").(string)

	client := ActiveDirectoryClient{
		username:                   username,
		password:                   password,
		server:                     server,
		usessl:                     usessl,
		usessh:                     usessh,
		authentication:             authentication,
		default_computer_container: default_computer_container,
	}

	return &client, nil
}

type ActiveDirectoryClient struct {
	username                   string
	password                   string
	server                     string
	usessl                     string
	usessh                     string
	authentication             string
	default_computer_container string
}

// Configuration options for the provider logging
type LoggingConfig struct {
	// The log level to use
	LogLevel logger.LogLevel
	// The path to the log file
	LogFile string
}
