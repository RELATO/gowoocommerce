package gowoocommerce

import (
	"encoding/json"
	"errors"
	"fmt"
)

// WooRequest is implemented for Batch/Post and Get
type WooRequest interface {
	Send(w *WooConnection) ([]byte, error)
}

// WooPostRequest can be used for synchronous requests to the products, attributes, or categories endpoint
type WooPostRequest struct {
	Endpoint string
	Payload  WooItem
}

// Send implements the WooRequest interface
func (p WooPostRequest) Send(w *WooConnection) ([]byte, error) {
	var err error
	if w.initialized == false {
		return nil, fmt.Errorf("Please initialize with your credentials first. WooConnection.Init()")
	}

	body, err := json.Marshal(p.Payload)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for i := 0; i < w.maxRetries; i++ {
		resp, err := w.Request("POST", p.Endpoint, body)
		if err == nil {
			return resp, nil
		}
		fmt.Println(err)
	}
	return nil, fmt.Errorf("Error sending request - %v", err)
}

// WooBatchPostRequest sends a payload of batch creations, updates and/or deletions
type WooBatchPostRequest struct {
	Endpoint string    `json:"-"`
	Create   []WooItem `json:"create,omitempty"` // Create requests must not have IDs -the WC backend will generate them
	Update   []WooItem `json:"update,omitempty"` // Update requests must have IDs
	Delete   []int     `json:"delete,omitempty"` // Delete requests can only be IDs
}

// Send implements the WooRequest Interface
func (b WooBatchPostRequest) Send(w *WooConnection) ([]byte, error) {
	if w.initialized == false {
		return nil, errors.New("Please initialize with your credentials first. WooConnection.Init()")
	}

	vars := struct {
		url  string
		body []byte
		err  error
		resp []byte
	}{
		url: w.credentials.domain + b.Endpoint,
	}

	vars.body, vars.err = json.Marshal(b)
	if vars.err != nil {
		return nil, vars.err
	}

	for i := 0; i < w.maxRetries; i++ {
		vars.resp, vars.err = w.Request("POST", b.Endpoint, vars.body) //w.Post(b.Endpoint, vars.body)
		if vars.err == nil {
			return vars.resp, nil
		}
		fmt.Println(vars.err)
	}
	return nil, fmt.Errorf("Error sending request - %v", vars.err)
}

// WooGetRequest implements GET request via a WooConnection
type WooGetRequest struct {
	Endpoint string
}

// Send implementes the WooRequest interface
func (g WooGetRequest) Send(w *WooConnection) ([]byte, error) {
	if w.initialized == false {
		return nil, errors.New("Please initialize with your credentials first. WooConnection.Init()")
	}

	var err error
	for i := 0; i < w.maxRetries; i++ {
		resp, err := w.Request("GET", g.Endpoint, nil)
		if err == nil {
			return resp, nil
		}
		fmt.Println(err)
	}
	return nil, fmt.Errorf("Error sending request - %v", err)
}
