package stake

import (
	"github.com/noah-blockchain/noah-explorer-api/helpers"
	"github.com/noah-blockchain/noah-explorer-api/resource"
	"github.com/noah-blockchain/noah-explorer-tools/models"
)

type Resource struct {
	Coin     string `json:"coin"`
	Address  string `json:"address"`
	Value    string `json:"value"`
	BipValue string `json:"bip_value"`
}

func (Resource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	stake := model.(models.Stake)

	return Resource{
		Coin:     stake.Coin.Symbol,
		Address:  stake.OwnerAddress.GetAddress(),
		Value:    helpers.QNoahStr2Noah(stake.Value),
		BipValue: helpers.QNoahStr2Noah(stake.NoahValue),
	}
}
