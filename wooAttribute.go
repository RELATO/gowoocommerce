package gowoocommerce

// WooAttribute provides additional general fields for the products
type WooAttribute struct {
	ID      int32    `json:"id,omitempty"`
	Name    string   `json:"name,omitempty"`
	Option  string   `json:"option,omitempty"`  // "term"
	Options []string `json:"options,omitempty"` // "terms"
	Slug    string   `json:"slug,omitempty"`
	Visible bool     `json:"visible,omitempty"`
	Type    string   `json:"type,omitempty"` // "select" by default
}

// GetID implements WooItem
func (a WooAttribute) GetID() int32 {
	return a.ID
}
