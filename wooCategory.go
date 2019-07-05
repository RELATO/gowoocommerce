package gowoocommerce

// WooCategoryLink can be either self, collection, or up
// e.g.: "https://example.com/wp-json/wc/v3/products/categories/15"
type WooCategoryLink struct {
	Href string `json:"href,omitempty"`
}

// WooCategoryLinks are generated from the category ids:
// e.g.: "href": "https://example.com/wp-json/wc/v3/products/categories/15"
type WooCategoryLinks struct {
	Self       []WooCategoryLink `json:"self,omitempty"`
	Collection []WooCategoryLink `json:"collection,omitempty"`
	Up         []WooCategoryLink `json:"up,omitempty"`
}

// WooCategory convers objects relating to the WC Category tree
type WooCategory struct {
	ID          int32            `json:"id,omitempty"`
	Name        string           `json:"name"`
	Alt         string           `json:"alt,omitempty"`
	Slug        string           `json:"slug,omitempty"`
	Parent      int32            `json:"parent,omitempty"`
	Description string           `json:"description,omitempty"`
	Image       WooImage         `json:"image,omitempty"`
	MenuOrder   int32            `json:"menu_order,omitempty"`
	Count       int32            `json:"count,omitempty"`
	Links       WooCategoryLinks `json:"_links,omitempty"` // read-only
}

// GetID implements WooItem
func (c WooCategory) GetID() int32 {
	return c.ID
}
