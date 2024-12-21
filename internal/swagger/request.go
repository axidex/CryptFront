package swagger

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"net/http"
)

func RequestOpenApi(client *resty.Client, openApiRoute string) (json.RawMessage, error) {
	resp, err := client.R().Get(openApiRoute)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("unexpected response")
	}

	return resp.Body(), nil
}
