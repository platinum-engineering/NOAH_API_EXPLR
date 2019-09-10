package delegation

import (
	"github.com/noah-blockchain/noah-explorer-api/helpers"
	"github.com/noah-blockchain/noah-explorer-api/resource"
	"github.com/noah-blockchain/noah-explorer-tools/models"
)

type Resource struct {
	Coin     string `json:"coin"`
	PubKey   string `json:"pub_key"`
	Value    string `json:"value"`
	BipValue string `json:"bip_value"`
}

func (resource Resource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	stake := model.(models.Stake)

	return Resource{
		Coin:     stake.Coin.Symbol,
		PubKey:   stake.Validator.GetPublicKey(),
		Value:    helpers.QNoahStr2Noah(stake.Value),
		BipValue: helpers.QNoahStr2Noah(stake.NoahValue),
	}
}
