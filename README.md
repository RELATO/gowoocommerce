# gowoocommerce

Go package to interface with the WooCommerce Product API. Based on a closed source project for [stillgrove](https://stillgrove.com). 
Currently only supports Creating, Reading, Updateing, Deleting Products. Other features like manipulating Category Trees or Attributes are missing: Happy about every contribution :)

## Installing
```
go get "github.com/michael-stiller/gowoocommerce"
```

## Examples:

### Initialize
```
import gwc "github.com/michael-stiller/gowoocommerce"

var w gwc.WooConnection

domain := "www.domain.com"
key := "yourkey"
secret := "yoursecret"

err := w.Init(domain, key, secret)
if err != nil {
    panic(err)
}

```
### Query existing products:
```
products, _ := w.GetAllProducts()
fmt.Println(products)
```

### Query Existing Categories:
https://woocommerce.github.io/woocommerce-rest-api-docs/#list-all-product-categories

```
// query parameters can be appended via query string
cats, _ := w.QueryCategories("")
```

### Batch Create/Update/Delete products
```
// Define a new product
var newProduct = gwc.WooProduct{
    Name:             "test_product",
    ShortDescription: "Very good product",
}

// Append Product to the Create/Update/Delete array in the request struct
var req = gwc.WooBatchRequest{
    Endpoint: "/wp-json/wc/v3/products/batch",
    Create:   []gwc.WooProduct{newProduct},
}

// Push request to queue
w.PushToQueue(req)

// Execute queue
// (returns an array of raw json []byte in case you want to further use the response object)
_, err = w.ExecuteRequestQueue()
if err != nil {
    panic(err)
}
```

## Authors
* **Michael Stiller** - *Initial work* - [michael-stiller](https://github.com/michael-stiller)

See also the list of [contributors](https://github.com/michael-stiller/gotradedoubler/contributors) who participated in this project.

## License

This project is licensed under the GNUgpl3 License - see the [LICENSE.md](LICENSE.md) file for details
