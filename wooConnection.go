package gowoocommerce

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/net/publicsuffix"
)

var loadPageSize = 100

type wooCredentials struct {
	domain string
	key    string
	secret string
}

// WooConnection interfaces with the WooCommerce backend
type WooConnection struct {
	initialized           bool
	credentials           wooCredentials
	jar                   *cookiejar.Jar
	maxRetries            int
	batchStrideSize       int // defines the size of one chunk for the batch upload
	maxConcurrentRequests int // defines how many requests can be sent concurrently
	requestQueue          []WooRequest
}

// Init takes in the credentials before dong any other operation
func (w *WooConnection) Init(domain, key, secret string, productsPerBatch, maxConcurrentRequests, maxRetries int) error {
	w.credentials.domain = domain
	w.credentials.key = key
	w.credentials.secret = secret

	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		return err
	}

	w.jar = jar

	// Just rule of thumb start values to be optimized based on request sizes and time constraints
	w.maxRetries = maxRetries
	w.batchStrideSize = productsPerBatch
	w.maxConcurrentRequests = maxConcurrentRequests

	w.initialized = true

	return nil
}

// Request sends a request: ("GET", "POST"), endpoint, body
func (w *WooConnection) Request(method, endpoint string, body []byte) ([]byte, error) {
	url := w.buildLink(endpoint)

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(w.credentials.key, w.credentials.secret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Jar: w.jar}
	rsp, err := client.Do(req)
	if err == nil {
		defer rsp.Body.Close()
		b, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}

		if rsp.StatusCode != http.StatusOK && rsp.StatusCode != http.StatusCreated {
			return nil, fmt.Errorf("Failed: %s: %s \n %s - %s", method, endpoint, rsp.Status, string(b))
		}
		return b, nil
	}

	return nil, err
}

// GetAllProducts returns all products from the WC backend
func (w *WooConnection) GetAllProducts(verbose bool) ([]WooProduct, error) {
	var currentProducts []WooProduct

	if w.initialized == false {
		return currentProducts, errors.New("Please initialize with your credentials first. WooConnection.Init()")
	}

	totalNumProducts, err := w.getNumItems("/wp-json/wc/v3/products?per_page=1") //get total number of items from product endpoint
	if err != nil {
		return currentProducts, err
	}

	if loadPageSize > totalNumProducts {
		loadPageSize = totalNumProducts
	}

	for offset := 0; offset < totalNumProducts; offset += loadPageSize {
		endpoint := "/wp-json/wc/v3/products?offset=" + strconv.Itoa(offset)
		endpoint += "&per_page=" + strconv.Itoa(loadPageSize)

		w.PushToQueue(
			WooGetRequest{
				Endpoint: endpoint,
			},
		)
	}
	rawResponse, err := w.ExecuteRequestQueue(true, verbose)
	if err != nil {
		return currentProducts, err
	}

	//currentProducts = make([]WooProduct, 0)
	for i := range rawResponse {
		if len(rawResponse[i]) < 1 {
			continue
		}
		var p []WooProduct
		err = json.Unmarshal(rawResponse[i], &p)
		if err != nil {
			continue
		}
		for j := range p {
			currentProducts = append(currentProducts, p[j])
		}
	}

	return currentProducts, nil
}

// PurgeProducts deletes all the products from the woo commerce backend
// Remember: Does not remove the image assets from the server!
func (w *WooConnection) PurgeProducts(verbose bool) error {
	if w.initialized == false {
		return fmt.Errorf("Please initialize with your credentials first. WooConnection.Init()")
	}

	products, err := w.GetAllProducts(true)
	if err != nil {
		return err
	}

	endpoint := "/wp-json/wc/v3/products/batch"

	var r = WooBatchPostRequest{
		Endpoint: endpoint,
	}

	// Preparing queue of product ids to be deleted
	for i := range products {
		r.Delete = append(r.Delete, int(products[i].ID))

		if i%16 == 0 || i >= len(products)-1 {
			w.PushToQueue(r)
			r = WooBatchPostRequest{
				Endpoint: endpoint,
			}
		}
	}

	_, err = w.ExecuteRequestQueue(true, verbose)
	if err != nil {
		return err
	}

	return nil
}

// QueryCategories returns all categories from the WC backend
// https://woocommerce.github.io/woocommerce-rest-api-docs/#list-all-product-categories
func (w *WooConnection) QueryCategories(searchString string) ([]WooCategory, error) {
	var categories []WooCategory
	pageSize := 10

	if w.initialized == false {
		return categories, errors.New("Please initialize with your credentials first. WooConnection.Init()")
	}

	endpoint := "/wp-json/wc/v3/products/categories"

	totalNumCats, err := w.getNumItems(fmt.Sprintf("%s?per_page=1%s", endpoint, searchString)) //get total number of items from category endpoint
	if err != nil {
		return categories, err
	}

	endpoint += fmt.Sprintf("?per_page=%d%s", pageSize, searchString)

	var r WooGetRequest

	nPage := 1
	for i := 0; i < totalNumCats; i += pageSize {
		r = WooGetRequest{
			Endpoint: fmt.Sprintf("%s&page=%d", endpoint, nPage),
		}
		w.PushToQueue(r)
		nPage++
	}
	rawResonse, err := w.ExecuteRequestQueue(true, false)
	if err != nil {
		return categories, err
	}

	categories = make([]WooCategory, totalNumCats)
	for i := range rawResonse {
		var c []WooCategory
		err = json.Unmarshal(rawResonse[i], &c)
		if err != nil {
			return categories, err
		}
		for j := range c {
			idx := (pageSize * i) + j
			categories[idx] = c[j]
		}
	}

	return categories, nil
}

// PushToQueue appends a WooREquest to the queue to later be executed
func (w *WooConnection) PushToQueue(r WooRequest) {
	w.requestQueue = append(w.requestQueue, r)
}

// ExecuteRequestQueue executes all the request that were pushed before and returns an array of the raw responses as bytes
// if strict: returns on any error; else: finishes regardless of errors
func (w *WooConnection) ExecuteRequestQueue(strict, verbose bool) ([][]byte, error) {
	var rawResponse [][]byte

	if len(w.requestQueue) == 0 {
		return rawResponse, nil
	}
	var wg sync.WaitGroup

	input := make(chan WooRequest, len(w.requestQueue))
	output := make(chan []byte, len(w.requestQueue))
	errChan := make(chan error, len(w.requestQueue))

	// Increment waitgroup counter and create go routines
	for i := 0; i < w.maxConcurrentRequests; i++ {
		wg.Add(1)
		go func(input chan WooRequest, output chan []byte) {
			defer wg.Done()

			for req := range input {
				resp, err := req.Send(w)
				errChan <- err
				output <- resp
			}
		}(input, output)
	}

	// Producer: load up input channel with jobs
	for _, job := range w.requestQueue {
		input <- job
	}
	fmt.Printf("%d scheduled \n", len(w.requestQueue))

	close(input)

	rawResponse = make([][]byte, len(w.requestQueue))
	var i int
	for res := range output {
		rawResponse[i] = res
		i++
		if verbose == true {
			progressBar(i, len(w.requestQueue))
		}

		if i >= len(w.requestQueue) {
			close(output)
			close(errChan)
			break
		}
	}

	for err := range errChan {
		if err != nil {
			fmt.Println(err)
			if strict == true {
				return rawResponse, err
			}
		}
	}

	wg.Wait()

	w.requestQueue = nil

	return rawResponse, nil
}

// ViewRequestQueue returns the marshalled requests as they will be sent by ExecuteRequestQueue
func (w *WooConnection) ViewRequestQueue() ([][]byte, error) {
	output := make([][]byte, len(w.requestQueue))

	var err error
	var errorCounter uint32

	var wg sync.WaitGroup
	for i := range w.requestQueue {
		wg.Add(1)
		go func(it int) {
			defer wg.Done()
			output[it], err = json.Marshal(w.requestQueue[it])
			if err != nil {
				fmt.Println(err)
				atomic.AddUint32(&errorCounter, 1)
			}

		}(i)
		if i%8 == 0 {
			wg.Wait()
		}
		if errorCounter > 0 {
			err = fmt.Errorf("Encountered %d error in %d requests", errorCounter, len(w.requestQueue))
			break
		}
	}
	wg.Wait()

	return output, err
}

// getNumItems returns the total number of items (products, categories) from the repsonse header of a given endpoint
func (w *WooConnection) getNumItems(endpoint string) (int, error) {
	if w.initialized == false {
		return 0, errors.New("Please initialize with your credentials first. WooConnection.Init()")
	}

	url := w.buildLink(endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.SetBasicAuth(w.credentials.key, w.credentials.secret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Jar: w.jar}
	rsp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	if rsp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Failed: %s\n %s", endpoint, rsp.Status)
	}

	numP := rsp.Header.Get("X-WP-Total")
	np, err := strconv.ParseInt(numP, 0, 32)
	if err != nil {
		return 0, fmt.Errorf("Unable to gather total number of products - %v", err)
	}

	return int(np), nil
}

func (w *WooConnection) buildLink(endpoint string) string {
	url := w.credentials.domain + endpoint
	if strings.HasPrefix(url, "https") == true {
		if strings.Index(url, "?") >= 0 {
			url += "&"
		} else {
			url += "?"
		}
		url += fmt.Sprintf("consumer_key=%s&consumer_secret=%s", w.credentials.key, w.credentials.secret)
	}
	return url
}
