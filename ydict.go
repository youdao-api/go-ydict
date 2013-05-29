/*
	ydict is the SDK for youdao dictionary/fanyi open api.
	http://fanyi.youdao.com/openapi
*/
package ydict

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

/*
	Result is the data structure for query result
*/
type Result struct {
	// Error-code
	ErrorCode int
	// Query string. Could be different from the request.
	Query string
	// Translations
	Translation []string
	// Basic dictionary result. Could be <nil>
	Basic *struct {
		// Phonetic data. Could be ""
		Phonetic string
		// Explains
		Explains []string
	}
	// Web mining dictionary result. Could be a zero-length slice
	Web []struct {
		// Entry key
		Key string
		// Explains
		Value []string
	}
}

/*
	Client is the data struct for a ydict client
*/
type Client struct {
	// Base URL
	BaseURL string
	// Keyfrom
	Keyfrom string
	// Key
	Key string
}

var OnlineBaseURL string = "http://fanyi.youdao.com/"

/*
	NewClient returns a *Client with baseURL, keyfrom and key set
*/
func NewClient(baseURL, keyfrom, key string) *Client {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}

	return &Client{
		BaseURL: baseURL,
		Keyfrom: keyfrom,
		Key:     key,
	}
}

/*
	NewClient returns a *Client with keyfrom and key set using BaseURL
*/
func NewOnlineClient(keyfrom, key string) *Client {
	return NewClient(OnlineBaseURL, keyfrom, key)
}

type result struct {
	ErrorCode   int      `json:"errorCode"`
	Query       string   `json:"query"`
	Translation []string `json:"translation"`
	Basic       *struct {
		Phonetic string   `json:"phonetic"`
		Explains []string `json:"explains"`
	} `json:"basic"`
	Web []struct {
		Key   string   `json:"key"`
		Value []string `json:"value"`
	} `json:"web"`
}

func (r *result) asResult() *Result {
	res := &Result{
		ErrorCode:   r.ErrorCode,
		Query:       r.Query,
		Translation: r.Translation,
	}
	// copy Basic field
	if r.Basic != nil {
		res.Basic = &struct {
			Phonetic string
			Explains []string
		}{
			Phonetic: r.Basic.Phonetic,
			Explains: r.Basic.Explains,
		}
	}

	// copy Web field
	res.Web = make([]struct {
		Key   string
		Value []string
	}, len(r.Web))
	for i := range r.Web {
		res.Web[i].Key = r.Web[i].Key
		res.Web[i].Value = r.Web[i].Value
	}

	return res
}

/*
	Query returns the result of a query
*/
func (c *Client) Query(q string) (*Result, error) {
	return c.QueryHttp(http.DefaultClient, q)
}

/*
	Query returns the result of a query with customized *http.Client
*/
func (c *Client) QueryHttp(httpClient *http.Client, q string) (*Result, error) {
	requestURL := fmt.Sprintf(
		"%sopenapi.do?keyfrom=%s&key=%s&type=data&doctype=json&version=1.1&q=%s",
		c.BaseURL, template.URLQueryEscaper(c.Keyfrom),
		template.URLQueryEscaper(c.Key), template.URLQueryEscaper(q))
	fmt.Println(requestURL)

	resp, err := httpClient.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	var res result
	err = dec.Decode(&res)
	if err != nil {
		return nil, err
	}

	return res.asResult(), nil
}
