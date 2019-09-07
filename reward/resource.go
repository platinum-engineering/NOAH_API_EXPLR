package reward

import (
	"github.com/noah-blockchain/noah-explorer-api/helpers"
	"github.com/noah-blockchain/noah-explorer-api/resource"
	"github.com/noah-blockchain/noah-explorer-tools/models"
	"time"
)

type Resource struct {
	BlockID   uint64 `json:"block"`
	Role      string `json:"role"`
	Amount    string `json:"amount"`
	Address   string `json:"address"`
	Validator string `json:"validator"`
	Timestamp string `json:"timestamp"`
}

func (Resource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	reward := model.(models.Reward)

	return Resource{
		BlockID:   reward.BlockID,
		Role:      reward.Role,
		Amount:    helpers.QNoahStr2Noah(reward.Amount),
		Address:   reward.Address.GetAddress(),
		Validator: reward.Validator.GetPublicKey(),
		Timestamp: reward.Block.CreatedAt.Format(time.RFC3339),
	}
}
