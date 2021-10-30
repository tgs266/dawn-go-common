package testing

import (
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

func TestGetRequestParams(app *fiber.App, endpoint string, params url.Values, token string, response interface{}) int {
	req, _ := http.NewRequest("GET", "http://test.com"+endpoint, strings.NewReader(params.Encode()))
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
