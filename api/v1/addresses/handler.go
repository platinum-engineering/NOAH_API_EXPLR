package addresses

import (
	"github.com/noah-blockchain/noah-explorer-api/address"
	"github.com/noah-blockchain/noah-explorer-api/chart"
	"github.com/noah-blockchain/noah-explorer-api/core"
	"github.com/noah-blockchain/noah-explorer-api/core/config"
	"github.com/noah-blockchain/noah-explorer-api/delegation"
	"github.com/noah-blockchain/noah-explorer-api/errors"
	"github.com/noah-blockchain/noah-explorer-api/events"
	"github.com/noah-blockchain/noah-explorer-api/helpers"
	"github.com/noah-blockchain/noah-explorer-api/resource"
	"github.com/noah-blockchain/noah-explorer-api/reward"
	"github.com/noah-blockchain/noah-explorer-api/slash"
	"github.com/noah-blockchain/noah-explorer-api/tools"
	"github.com/noah-blockchain/noah-explorer-api/transaction"
	"github.com/noah-blockchain/noah-explorer-tools/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetAddressRequest struct {
	Address string `uri:"address" binding:"noahAddress"`
}

type GetAddressesRequest struct {
	Addresses []string `form:"addresses[]" binding:"required,noahAddress,max=50"`
}

// TODO: replace string to int
type FilterQueryRequest struct {
	StartBlock *string `form:"startblock" binding:"omitempty,numeric"`
	EndBlock   *string `form:"endblock"   binding:"omitempty,numeric"`
	Page       *string `form:"page"       binding:"omitempty,numeric"`
}

type StatisticsQueryRequest struct {
	Scale     *string `form:"scale"     binding:"omitempty,eq=minute|eq=hour|eq=day"`
	StartTime *string `form:"startTime" binding:"omitempty,timestamp"`
	EndTime   *string `form:"endTime"   binding:"omitempty,timestamp"`
}

// Get list of addresses
func GetAddresses(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetAddressesRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// remove Noah wallet prefix from each address
	noahAddresses := make([]string, len(request.Addresses))
	for key, addr := range request.Addresses {
		noahAddresses[key] = helpers.RemoveNoahPrefix(addr)
	}

	// fetch addresses
	addresses := explorer.AddressRepository.GetByAddresses(noahAddresses)

	// extend the model array with empty model if not exists
	if len(addresses) != len(noahAddresses) {
		for _, item := range noahAddresses {
			if isModelsContainAddress(item, addresses) {
				continue
			}

			addresses = append(addresses, *makeEmptyAddressModel(item, explorer.Environment.BaseCoin))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resource.TransformCollection(addresses, address.Resource{}),
	})
}

// Get address detail
func GetAddress(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	noahAddress, err := getAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch address
	model := explorer.AddressRepository.GetByAddress(*noahAddress)

	// if no models found
	if model == nil {
		model = makeEmptyAddressModel(*noahAddress, explorer.Environment.BaseCoin)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": new(address.Resource).Transform(*model),
	})
}

// Get list of transactions by Noah address
func GetTransactions(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	noahAddress, err := getAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// validate request query
	var requestQuery FilterQueryRequest
	err = c.ShouldBindQuery(&requestQuery)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	pagination := tools.NewPagination(c.Request)
	txs := explorer.TransactionRepository.GetPaginatedTxsByAddresses(
		[]string{*noahAddress},
		transaction.BlocksRangeSelectFilter{
			StartBlock: requestQuery.StartBlock,
			EndBlock:   requestQuery.EndBlock,
		}, &pagination)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(txs, transaction.Resource{}, pagination))
}

// Get list of rewards by Noah address
func GetRewards(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	filter, pagination, err := prepareEventsRequest(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	rewards := explorer.RewardRepository.GetPaginatedByAddress(*filter, pagination)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(rewards, reward.Resource{}, *pagination))
}

// Get list of slashes by Noah address
func GetSlashes(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	filter, pagination, err := prepareEventsRequest(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	slashes := explorer.SlashRepository.GetPaginatedByAddress(*filter, pagination)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(slashes, slash.Resource{}, *pagination))
}

// Get list of delegations by Noah address
func GetDelegations(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	noahAddress, err := getAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	pagination := tools.NewPagination(c.Request)
	delegations := explorer.StakeRepository.GetByAddress(*noahAddress, &pagination)

	// fetch total delegated sum
	totalDelegated, err := explorer.StakeRepository.GetSumInBipValueByAddress(*noahAddress)
	helpers.CheckErr(err)

	// create additional field
	additionalFields := map[string]interface{}{"total_delegated_bip_value": helpers.QNoahStr2Noah(totalDelegated)}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollectionWithAdditionalFields(
		delegations, delegation.Resource{}, pagination, additionalFields))
}

// Get rewards statistics by noah address
func GetRewardsStatistics(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	noahAddress, err := getAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	var requestQuery StatisticsQueryRequest
	if err := c.ShouldBindQuery(&requestQuery); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// set scale instead of default if exists
	scale := config.DefaultStatisticsScale
	if requestQuery.Scale != nil {
		scale = *requestQuery.Scale
	}

	// fetch data
	chartData := explorer.RewardRepository.GetChartData(*noahAddress, chart.SelectFilter{
		Scale:     scale,
		StartTime: requestQuery.StartTime,
		EndTime:   requestQuery.EndTime,
	})

	c.JSON(http.StatusOK, gin.H{
		"data": resource.TransformCollection(chartData, chart.RewardResource{}),
	})
}

func prepareEventsRequest(c *gin.Context) (*events.SelectFilter, *tools.Pagination, error) {
	noahAddress, err := getAddressFromRequestUri(c)
	if err != nil {
		return nil, nil, err
	}

	var requestQuery FilterQueryRequest
	if err := c.ShouldBindQuery(&requestQuery); err != nil {
		return nil, nil, err
	}

	pagination := tools.NewPagination(c.Request)

	return &events.SelectFilter{
		Address:    *noahAddress,
		StartBlock: requestQuery.StartBlock,
		EndBlock:   requestQuery.EndBlock,
	}, &pagination, nil
}

// Get noah address from current request uri
func getAddressFromRequestUri(c *gin.Context) (*string, error) {
	var request GetAddressRequest
	if err := c.ShouldBindUri(&request); err != nil {
		return nil, err
	}

	noahAddress := helpers.RemoveNoahPrefix(request.Address)
	return &noahAddress, nil
}

// Return model address with zero base coin
func makeEmptyAddressModel(noahAddress string, baseCoin string) *models.Address {
	return &models.Address{
		Address: noahAddress,
		Balances: []*models.Balance{{
			Coin: &models.Coin{
				Symbol: baseCoin,
			},
			Value: "0",
		}},
	}
}

// Check that array of address models contain exact noah address
func isModelsContainAddress(noahAddress string, models []models.Address) bool {
	for _, item := range models {
		if item.Address == noahAddress {
			return true
		}
	}

	return false
}
