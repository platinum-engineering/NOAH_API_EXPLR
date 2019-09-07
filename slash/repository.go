package slash

import (
	"github.com/noah-blockchain/noah-explorer-api/events"
	"github.com/noah-blockchain/noah-explorer-api/helpers"
	"github.com/noah-blockchain/noah-explorer-api/tools"
	"github.com/noah-blockchain/noah-explorer-tools/models"
	"github.com/go-pg/pg"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repository Repository) GetPaginatedByAddress(filter events.SelectFilter, pagination *tools.Pagination) []models.Slash {
	var slashes []models.Slash
	var err error

	pagination.Total, err = repository.db.Model(&slashes).
		Column("Coin.symbol", "Address.address", "Validator.public_key", "Block.created_at").
		Apply(filter.Filter).
		Apply(pagination.Filter).
		Order("block_id DESC").
		SelectAndCount()

	helpers.CheckErr(err)

	return slashes
}
