package commercetools

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/labd/commercetools-go-sdk/platform"
)

func resourceProduct() *schema.Resource {
	return &schema.Resource{
		Description:   "Product are used to describe common characteristics, most importantly common custom",
		CreateContext: resourceProductCreate,
		ReadContext:   resourceProductRead,
		UpdateContext: resourceProductUpdate,
		DeleteContext: resourceProductDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description:      "[LocalizedString](https://docs.commercetools.com/api/types#localizedstring)",
				Type:             TypeLocalizedString,
				ValidateDiagFunc: validateLocalizedStringKey,
				Required:         true,
			},
			"slug": {
				Description:      "[LocalizedString](https://docs.commercetools.com/api/types#localizedstring)",
				Type:             TypeLocalizedString,
				ValidateDiagFunc: validateLocalizedStringKey,
				Required:         true,
			},
			"product_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of Product Type",
			},
		},
	}
}

func resourceProductCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := getClient(m)

	draft := platform.ProductDraft{
		Name:        expandLocalizedString(d.Get("name")),
		Slug: expandLocalizedString(d.Get("slug")),
	}

	if d.Get("product_type").(string) != "" {
		productType := platform.ProductTypeResourceIdentifier{}
		productType.ID = stringRef(d.Get("product_type"))
		draft.ProductType = productType
	}

	var ctType *platform.Product
	err := resource.RetryContext(ctx, 20*time.Second, func() *resource.RetryError {
		var err error

		ctType, err = client.Products().Post(draft).Execute(ctx)
		return processRemoteError(err)
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ctType.ID)
	d.Set("version", ctType.Version)

	return resourceProductRead(ctx, d, m)
}

func resourceProductRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := getClient(m)

	ctType, err := client.Products().WithId(d.Id()).Get().Execute(ctx)
	if err != nil {
		if IsResourceNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if ctType == nil {
		d.SetId("")
	} else {
		d.Set("version", ctType.Version)
		d.Set("key", ctType.Key)

	}
	return nil
}

func resourceProductUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	return nil
}

func resourceProductDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	return nil
}

