package hostctl

import (
	"context"
	"os"
	"runtime"

	"github.com/guumaster/hostctl/pkg/file"
	"github.com/guumaster/hostctl/pkg/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceProfileRead,
		CreateContext: resourceProfileCreate,
		DeleteContext: resourceProfileDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"profile": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"domains": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func resourceProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	resourceProfileRead(ctx, d, m)
	return nil
}

func resourceProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	profile := d.Get("profile").(string)
	// domains := d.Get("domains").([]interface{})
	ip := "127.0.0.1"

	if !d.HasChange("profile") && !d.HasChange("domains") {
		return nil
	}

	// for _, domain := range domains {
	// 	i := domain.(map[string]interface{})
	// }

	hostnames := []string{"my.web.com"}

	h, err := file.NewFile(getDefaultHostFile())
	if err != nil {
		return diag.FromErr(err)
	}

	err = h.AddRoute(profile, types.NewRoute(ip, hostnames...))
	if err != nil {
		return diag.FromErr(err)
	}

	err = h.Flush()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(profile)

	return diags
}

func getDefaultHostFile() string {
	if runtime.GOOS == "linux" {
		return "/etc/hosts" //nolint: goconst
	}

	envHostFile := os.Getenv("HOSTCTL_FILE")
	if envHostFile != "" {
		return envHostFile
	}

	if runtime.GOOS == "windows" {
		return `C:/Windows/System32/Drivers/etc/hosts`
	}

	return "/etc/hosts"
}
