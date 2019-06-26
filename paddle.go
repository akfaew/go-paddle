package paddle

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	. "github.com/akfaew/aeutils"
	"github.com/google/go-querystring/query"
)

var defaultBaseURL = "https://vendors.paddle.com/api/2.0/"

type Conf struct {
	VendorID int
	APIKey   string

	// checkout
	SecretKey string
	ProductID int

	// webhook verification
	PublicKey *rsa.PublicKey
}

// Init loads the RSA Public Key from publicKeyPath into Conf.
func (c *Conf) Init(publicKeyPath string) error {
	pubPEM, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return err
	}

	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return fmt.Errorf("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse DER encoded public key: " + err.Error())
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		c.PublicKey = pub
	default:
		return fmt.Errorf("unknown type of public key")
	}

	return nil
}

type ProductService service
type SubscriptionService service

type Client struct {
	client *http.Client
	conf   *Conf
	// Base URL for API requests. baseURL should always be specified with a trailing slash.
	baseURL *url.URL

	// Services used for talking to different parts of the Paddle API.
	Subscription *SubscriptionService
	Product      *ProductService
}

type service struct {
	client *Client
}

func (conf *Conf) NewClient(ctx context.Context, client *http.Client) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		client:  client,
		conf:    conf,
		baseURL: baseURL,
	}
	s := &service{client: c}

	c.Subscription = (*SubscriptionService)(s)
	c.Product = (*ProductService)(s)

	return c
}

// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.baseURL.Path, "/") {
		return nil, fmt.Errorf("baseURL must have a trailing slash, but %q does not", c.baseURL)
	}
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	LogDebugfd(ctx, ">>>>>>>>>>>>>>>>>>>>>>>>>%+v<<<<<<<<<<<<<<", u.String())
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it. If rate limit is exceeded and reset time is in the future,
// Do returns *RateLimitError immediately without making a network API call.
//
// The provided ctx must be non-nil. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// If the error type is *url.Error
		if e, ok := err.(*url.Error); ok {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = url.String()
				return nil, e
			}
		}

		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := checkError(resp, data); err != nil {
		return resp, err
	}

	if v != nil {
		if err := json.Unmarshal(data, v); err != nil {
			return resp, err
		}
	}

	return resp, nil
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse represents an error returned by the Paddle API. Every
// response contains a field called "success". If it's not true, then
// something went wrong.
type ErrorResponse struct {
	response *http.Response // HTTP response that caused this error

	Success    bool  `json:"success"`
	ErrorField Error `json:"error"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: StatusCode: %d Code: %d Message: \"%v\" Success: %v)",
		r.response.Request.Method, r.response.Request.URL,
		r.response.StatusCode, r.ErrorField.Code, r.ErrorField.Message, r.Success)
}

// Every Paddle API response contains a field called "success". If it's not true, then something
// went wrong.
func checkError(r *http.Response, data []byte) error {
	errorResponse := &ErrorResponse{response: r}
	if data != nil {
		if err := json.Unmarshal(data, errorResponse); err != nil {
			return err
		}
	}

	if !errorResponse.Success {
		return errorResponse
	}

	return nil
}
