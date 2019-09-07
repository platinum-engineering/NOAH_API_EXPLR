package reward

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/noah-blockchain/noah-explorer-api/events"
	"github.com/noah-blockchain/noah-explorer-api/helpers"
	"github.com/noah-blockchain/noah-explorer-api/tools"
	"github.com/noah-blockchain/noah-explorer-tools/models"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Get filtered list of rewards by Noah address
func (repository Repository) GetPaginatedByAddress(filter events.SelectFilter, pagination *tools.Pagination) []models.Reward {
	var rewards []models.Reward
	var err error

	// get count of rewards
	pagination.Total, err = repository.db.Model(&rewards).
		Column("Address.address").
		Apply(filter.Filter).
		Count()
	helpers.CheckErr(err)

	if pagination.Total == 0 {
		return nil
	}

	// get rewards
	err = repository.db.Model(&rewards).
		Column("Address.address", "Validator.public_key", "Block.created_at").
		Apply(filter.Filter).
		Apply(pagination.Filter).
		Order("block.id DESC").
		Order("reward.amount").
		Select()
	helpers.CheckErr(err)

	return rewards
}

type ChartData struct {
	Time   time.Time
	Amount string
}

// Get filtered chart data by Noah address
func (repository Repository) GetChartData(address string, filter tools.Filter) []ChartData {
	var rewards models.Reward
	var chartData []ChartData

	err := repository.db.Model(&rewards).
		Column("Address._").
		ColumnExpr("SUM(amount) as amount").
		Where("address.address = ?", address).
		Apply(filter.Filter).
		Select(&chartData)

	helpers.CheckErr(err)

	return chartData
}
