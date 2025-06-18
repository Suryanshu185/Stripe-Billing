package handlers

import (
	"github.com/gin-gonic/gin"
	apimodels "gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/handlers/models"
	"gitlab.com/amcop-saas-platform/vcs/vcs/billing-microservice/services"
	"net/http"
)

type ResourceHandler struct {
	resourceService *services.ActiveResourceService
}

func NewResourceHandler(resourceService *services.ActiveResourceService) *ResourceHandler {
	return &ResourceHandler{
		resourceService: resourceService,
	}
}

func (s *ResourceHandler) GetUserResources(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UserID header is missing"})
		return
	}

	resources, err := s.resourceService.GetUserResource(userID)
	if err != nil {
		HandleApiErrorResponse(c, err)
		return
	}

	var apiResponse []*apimodels.ResourcesResponse

	for _, resource := range resources {
		apiResponse = append(apiResponse, &apimodels.ResourcesResponse{
			Provider:        resource.Provider,
			InstanceID:      resource.InstanceID,
			StartTime:       resource.StartTime,
			LastChargedTime: resource.LastChargedTime,
			HourlyRateCents: resource.HourlyRateCents,
		})
	}

	c.JSON(http.StatusOK, apiResponse)
}

func (s *ResourceHandler) AddResource(c *gin.Context) {
	var addResourceRequest apimodels.AddResourceRequest
	if err := c.ShouldBindJSON(&addResourceRequest); err != nil {
		HandleApiErrorResponse(c, err)
	}

	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UserID header is missing"})
		return
	}

	err := s.resourceService.AddResource(userID, addResourceRequest.InstanceID, addResourceRequest.Provider, addResourceRequest.HourlyRateCents)
	if err != nil {
		HandleApiErrorResponse(c, err)
		return
	}
	c.Status(http.StatusCreated)
}
