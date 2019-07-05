package gowoocommerce

// WooItem is the object used to hold requests and responses towards the woocommerce api
// Examples: WooProduct, WooAttribute, WooCategory
type WooItem interface {
	GetID() int32
}
