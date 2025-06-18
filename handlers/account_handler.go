package handlers

import (
	"github.com/gin-gonic/gin"
	apimodels "gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/handlers/models"
	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/services"
	"net/http"
)

type AccountHandler struct {
	accountService *services.AccountService
}

func NewAccountHandler(accountService *services.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var accountRequest apimodels.CreateAccountRequest

	//TODO [vikas] add middleware to extract userID
	userID := c.GetHeader("X-User-ID")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UserID header is missing"})
		c.Abort()
		return
	}

	if err := c.ShouldBindJSON(&accountRequest); err != nil {
		HandleApiErrorResponse(c, err)
		return
	}
	if err := h.accountService.CreateAccount(userID, accountRequest.Email, accountRequest.TopUpThresholdCents, accountRequest.TopUpAmountCents); err != nil {
		HandleApiErrorResponse(c, err)
		return
	}
	c.Status(http.StatusCreated)
}

func (h *AccountHandler) AddBalanceCents(c *gin.Context) {

}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UserID header is missing"})
		c.Abort()
		return
	}
	account, err := h.accountService.GetUserAccount(userID)
	if err != nil {
		HandleApiErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, account)
}
