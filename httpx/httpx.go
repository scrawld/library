package httpx

import (
	"net/http"
)

/**
 * Get issues a GET to the specified URL.
 * The caller can configure the new client by passing configuration options
 * Note: When resp is nil, you have to do it manually `response.Body.Close()`
 *
 * Example:
 *
 * var (
 * 	respBody = map[string]interface{}{}
 * )
 * _, err = Get("http://127.0.0.1:8080", &respBody)
 * or
 * _, err = Get("http://127.0.0.1:8080", respBody,
 * 	SetHeader(map[string]string{"Content-Type": "application/json"}),
 * 	SetTimeout(time.Second*10),
 * )
 */
func Get(url string, respBody interface{}, options ...ClientOptionFunc) (r *http.Response, err error) {
	return Call(url, respBody, options...)
}

/**
 * Post sends an HTTP request and returns an HTTP response, following
 * The caller can configure the new client by passing configuration options
 * Note: When resp is nil, you have to do it manually `response.Body.Close()`
 *
 * Example:
 *
 * var (
 *	req      = []byte(`{"name":"xxx","amount":1}`)
 * 	respBody = map[string]interface{}{}
 * )
 * _, err = PostJson("http://127.0.0.1:8080", req, &respBody)
 * or
 * _, err = Post("http://127.0.0.1:8080", req, respBody,
 * 	SetHeader(map[string]string{"Content-Type": "application/json"}),
 * 	SetTimeout(time.Second*10),
 * )
 */
func Post(url string, body interface{}, respBody interface{}, options ...ClientOptionFunc) (r *http.Response, err error) {
	options = append(options, SetMethod("POST"), SetBody(body))
	return Call(url, respBody, options...)
}
