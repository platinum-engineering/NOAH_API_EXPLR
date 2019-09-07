package slash

import (
	"github.com/noah-blockchain/noah-explorer-api/helpers"
	"github.com/noah-blockchain/noah-explorer-api/resource"
	"github.com/noah-blockchain/noah-explorer-tools/models"
	"time"
)

type Resource struct {
	BlockID   uint64 `json:"block"`
	Coin      string `json:"coin"`
	Amount    string `json:"amount"`
	Address   string `json:"address"`
	Validator string `json:"validator"`
	Timestamp string `json:"timestamp"`
}

func (Resource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	slash := model.(models.Slash)

	return Resource{
		BlockID:   slash.BlockID,
		Coin:      slash.Coin.Symbol,
		Amount:    helpers.QNoahStr2Noah(slash.Amount),
		Address:   slash.Address.GetAddress(),
		Validator: slash.Validator.GetPublicKey(),
		Timestamp: slash.Block.CreatedAt.Format(time.RFC3339),
	}
}
