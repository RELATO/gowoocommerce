package gowoocommerce

// WooImage contains all the information on product images
type WooImage struct {
	ID int32 `json:"id,omitempty"`
	//DateCreated     string `json:"date_created,omitempty"`
	DateCreatedGMT string `json:"date_created_gmt,omitempty"`
	//DateModified    string `json:"date_modified,omitempty"`
	DateModifiedGMT string `json:"date_modified_gmt,omitempty"`
	SRC             string `json:"src,omitempty"`
	Name            string `json:"name,omitempty"`
	Alt             string `json:"alt,omitempty"`
}

// WooTag interacts witht underlying category tree
type WooTag struct {
	ID   int32  `json:"id,omitempty"`
	Name string `json:"name,omitempty"` // read-only
	Slug string `json:"slug,omitempty"` // read-only
}

// WooDimension stores length, width, and height
type WooDimension struct {
	Length string `json:"length,omitempty"`
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

// WpmlPrice holds custom prices
type WpmlPrice struct {
	RegularPrice string `json:"regular_price,omitempty"`
	SalePrice    string `json:"sale_price,omitempty"`
}

// WooProduct is the struct through which you interface with the WooCommerce backend
type WooProduct struct {
	ID uint32 `json:"id,omitempty"` // read-only!!!
	//Key               uint32                   `json:"-"`
	SKU       string `json:"sku,omitempty"`
	Name      string `json:"name,omitempty"`
	Slug      string `json:"slug,omitempty"`
	Permalink string `json:"permalink,omitempty"` // read-only
	//DateCreated       string                   `json:"date_created,omitempty"`      // read-only
	DateCreatedGmt string `json:"date_created_gmt,omitempty"` // read-only
	//DateModified      string                   `json:"date_modified,omitempty"`     // read-only
	DateModifiedGmt   string `json:"date_modified_gmt,omitempty"` // read-only
	Type              string `json:"type,omitempty"`
	Status            string `json:"status,omitempty"`
	Featured          bool   `json:"featured,omitempty"`
	CatalogVisibility string `json:"catalog_visibility,omitempty"` // Options: visible, catalog, search and hidden. Default is visible.
	Description       string `json:"description,omitempty"`
	ShortDescription  string `json:"short_description,omitempty"`
	//Price             string                   `json:"price,omitempty"`         // read-only
	RegularPrice      string `json:"regular_price,omitempty"`
	SalePrice         string `json:"sale_price,omitempty"`
	DateOnSaleFrom    string `json:"date_on_sale_from,omitempty"`
	DateOnSaleFromGmt string `json:"date_on_sale_from_gmt,omitempty"`
	DateOnSaleTo      string `json:"date_on_sale_to,omitempty"`
	DateOnSaleToGmt   string `json:"date_on_sale_to_gmt,omitempty"`
	//PriceHTML         string                   `json:"price_html,omitempty"`   // read-only
	OnSale bool `json:"on_sale,omitempty"` // read-only
	//Purchasable       bool                     `json:"purchasable,omitempty"`  // read-only
	TotalSales        int32                    `json:"total_sales,omitempty"`  // read-only
	ExternalURL       string                   `json:"external_url,omitempty"` // real outlink
	ButtonText        string                   `json:"button_text,omitempty"`  // external shop to link to
	TaxStatus         string                   `json:"tax_status,omitempty"`   // Options: taxable, shipping and none. Default is taxable
	TaxClass          string                   `json:"tax_class,omitempty"`
	StockQuantity     int32                    `json:"stock_quantity,omitempty"`
	StockStatus       string                   `json:"stock_status,omitempty"`      // Options: instock, outofstock, onbackorder. Default is instock.
	SoldIndividually  bool                     `json:"sold_individually,omitempty"` // Allow one item to be bought in a single order. Default is false
	Weight            string                   `json:"weight,omitempty"`
	Dimensions        WooDimension             `json:"dimensions,omitempty"`
	ShippingRequired  bool                     `json:"shipping_required,omitempty"` // read-only
	ReviewsAllowed    bool                     `json:"reviews_allowed,omitempty"`   // default: true
	AverageRating     string                   `json:"average_rating,omitempty"`    // read-only
	RatingCount       int32                    `json:"rating_count,omitempty"`      // read-only
	RelatedIds        []int32                  `json:"related_ids,omitempty"`       // read_only
	UpsellIds         []int32                  `json:"upsell_ids,omitempty"`
	CrossSellIds      []int32                  `json:"cross_sell_ids,omitempty"`
	ParentID          int32                    `json:"parent_id,omitempty"`
	Categories        []WooCategory            `json:"categories,omitempty"`
	Tags              []WooTag                 `json:"tags,omitempty"`
	Images            []WooImage               `json:"images,omitempty"`
	DefaultAttributes []map[string]interface{} `json:"default_attributes,omitempty"`
	Variations        []string                 `json:"variations,omitempty"`
	GroupedProducts   []int32                  `json:"grouped_products,omitempty"`
	MenuOrder         int32                    `json:"menu_order,omitempty"`
	MetaData          []map[string]interface{} `json:"meta_data,omitempty"`
	Attributes        []WooAttribute           `json:"attributes,omitempty"`
	Brands            []interface{}            `json:"brands,omitempty"`
	Language          string                   `json:"language,omitempty"`
	Lang              string                   `json:"lang,omitempty"`          // relates to the woocommerce multilingual package; otherwise: omit!
	CustomPrices      map[string]WpmlPrice     `json:"custom_prices,omitempty"` // "custom_prices": {"EUR": {"regular_price": 100, "sale_price": 99}}
}

// GetID implements WooItem
func (p WooProduct) GetID() int32 {
	return int32(p.ID)
}

// AddImage adds an image +name ( +text to display when said image not available)
func (p *WooProduct) AddImage(url string, name string, text string) {
	p.Images = append(
		p.Images,
		WooImage{
			SRC:  url,
			Name: name,
			Alt:  text,
		},
	)
}
