package bipdev

import (
	"github.com/noah-blockchain/noah-explorer-api/bipdev/responses"
	"github.com/noah-blockchain/noah-explorer-api/tools"
)

type Api struct {
	client *tools.HttpClient
}

func NewApi(host string) *Api {
	return &Api{
		client: tools.NewHttpClient(host),
	}
}

func (api *Api) GetCurrentPrice() (*responses.CurrentPriceResponse, error) {
	response := &responses.CurrentPriceResponse{}
	err := api.client.Get("/api/price", response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
