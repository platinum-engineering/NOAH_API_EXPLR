package validators

import (
	"github.com/noah-blockchain/noah-explorer-api/core"
	"github.com/noah-blockchain/noah-explorer-api/errors"
	"github.com/noah-blockchain/noah-explorer-api/helpers"
	"github.com/noah-blockchain/noah-explorer-api/resource"
	"github.com/noah-blockchain/noah-explorer-api/tools"
	"github.com/noah-blockchain/noah-explorer-api/transaction"
	"github.com/noah-blockchain/noah-explorer-api/validator"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetValidatorRequest struct {
	PublicKey string `uri:"publicKey"    binding:"required,noahPubKey"`
}

// TODO: replace string to int
type GetValidatorTransactionsRequest struct {
	Page       string  `form:"page"        binding:"omitempty,numeric"`
	StartBlock *string `form:"startblock"  binding:"omitempty,numeric"`
	EndBlock   *string `form:"endblock"    binding:"omitempty,numeric"`
}

// Get list of transaction by validator public key
func GetValidatorTransactions(c *gin.Context) {
	var validatorRequest GetValidatorRequest
	var request GetValidatorTransactionsRequest

	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	err := c.ShouldBindUri(&validatorRequest)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// validate request query
	err = c.ShouldBindQuery(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	publicKey := helpers.RemoveNoahPrefix(validatorRequest.PublicKey)
	pagination := tools.NewPagination(c.Request)
	txs := explorer.TransactionRepository.GetPaginatedTxsByFilter(transaction.ValidatorFilter{
		ValidatorPubKey: publicKey,
		StartBlock:      request.StartBlock,
		EndBlock:        request.EndBlock,
	}, &pagination)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(txs, transaction.Resource{}, pagination))
}

// Get validator detail by public key
func GetValidator(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetValidatorRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	data := explorer.ValidatorRepository.GetByPublicKey(helpers.RemoveNoahPrefix(request.PublicKey))

	// check validator to existing
	if data == nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Validator not found.", c)
		return
	}

	activeValidatorIds := explorer.ValidatorRepository.GetActiveValidatorIds()
	totalStake := explorer.ValidatorRepository.GetTotalStakeByActiveValidators(activeValidatorIds)

	c.JSON(http.StatusOK, gin.H{
		"data": validator.Resource{}.Transform(data, activeValidatorIds, totalStake),
	})
}
