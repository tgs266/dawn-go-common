package testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func TestGetRequest(app *fiber.App, endpoint string, token string, response interface{}) int {
	req, _ := http.NewRequest("GET", "http://test.com"+endpoint, nil)
	req.Header.Set("Authorization", "token")

	httpResponse, _ := app.Test(req)
	defer httpResponse.Body.Close()
	json.NewDecoder(httpResponse.Body).Decode(&response)
	return httpResponse.StatusCode
}

func TestGetRequestParams(app *fiber.App, endpoint string, params map[string]string, token string, response interface{}) int {
	var args []string
	for k := range params {
		args = append(args, k+"="+params[k])
	}

	req, _ := http.NewRequest("GET", "http://test.com"+endpoint+"?"+strings.Join(args, "&"), nil)
	req.Header.Set("Authorization", "token")

	httpResponse, _ := app.Test(req)
	defer httpResponse.Body.Close()
	json.NewDecoder(httpResponse.Body).Decode(&response)
	return httpResponse.StatusCode
}

func TestPostRequest(app *fiber.App, endpoint string, params url.Values, token string, response interface{}) int {
	req, _ := http.NewRequest("POST", "http://test.com"+endpoint, strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "token")

	httpResponse, _ := app.Test(req)
	defer httpResponse.Body.Close()
	json.NewDecoder(httpResponse.Body).Decode(&response)
	return httpResponse.StatusCode
}

func TestPostRequestJson(app *fiber.App, endpoint string, params []byte, token string, response interface{}) int {
	req, _ := http.NewRequest("POST", "http://test.com"+endpoint, bytes.NewBuffer(params))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "token")

	httpResponse, _ := app.Test(req)
	defer httpResponse.Body.Close()
	json.NewDecoder(httpResponse.Body).Decode(&response)
	return httpResponse.StatusCode
}

func TestPutRequest(app *fiber.App, endpoint string, params url.Values, token string, response interface{}) int {
	req, _ := http.NewRequest("PUT", "http://test.com"+endpoint, strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "token")

	httpResponse, _ := app.Test(req)
	defer httpResponse.Body.Close()
	json.NewDecoder(httpResponse.Body).Decode(&response)
	return httpResponse.StatusCode
}
