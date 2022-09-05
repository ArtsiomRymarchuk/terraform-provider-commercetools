resource "commercetools_product" "new_product" {
  name = {
    en-US = "<Product_name>"
  }
  slug = {
    en-US = "<Product_slug>"
  }
  product_type = "<product_type_id>"
}
