package address

import (
	"github.com/noah-blockchain/noah-explorer-api/balance"
	"github.com/noah-blockchain/noah-explorer-api/resource"
	"github.com/noah-blockchain/noah-explorer-tools/models"
)

type Resource struct {
	Address  string               `json:"address"`
	Balances []resource.Interface `json:"balances"`
}

func (Resource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	address := model.(models.Address)

	return Resource{
		Address:  address.GetAddress(),
		Balances: resource.TransformCollection(address.Balances, balance.Resource{}),
	}
}
