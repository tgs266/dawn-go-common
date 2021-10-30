package testing

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func TestGetRequest(app *fiber.App, endpoint string, params *url.Values, token string, response interface{}) {
	req, _ := http.NewRequest("GET", "http://test.com"+endpoint, strings.NewReader(params.Encode()))
	req.Header.Set("Authorization", "token")

	httpResponse, _ := app.Test(req)
	defer httpResponse.Body.Close()
	json.NewDecoder(httpResponse.Body).Decode(&response)
}
