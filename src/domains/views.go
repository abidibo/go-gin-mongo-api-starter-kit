package domains

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"systems-management-api/auth"
	"systems-management-api/core/utils"
)

// Returns all domains, admin or superadmin roles required
// @Summary Domains list
// @Description Retrieves all domains
// @Security BearerAuth
// @Tags domains
// @Accept  json
// @Produce  json
// @Success 200 {array} []DomainData
// @Failure 403 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /domains/ [get]
func domainListView(c *gin.Context) {
	domainService := new(DomainService)
	domains, err := domainService.all()

	if err != nil {
		zap.S().Error("Error while getting all domains, Reason: ", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Message: "Cannot fetch domains",
		})
	} else {
		serializer := NewDomainSerializer()
		c.JSON(http.StatusOK, serializer.SerializeMany(domains))
	}
}

var DomainListView = auth.RoleRequired([]string{"admin", "superadmin"}, domainListView)

// Returns domain given its id
// @Summary Domain detail
// @Description Retrieves one domain given its id
// @Security BearerAuth
// @Tags domains
// @Accept  json
// @Produce  json
// @Param id path string true "Domain ID"
// @Success 200 {object} DomainData
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /domain/{id} [get]
func domainDetailView(c *gin.Context) {
	domainService := new(DomainService)
	domain, err := domainService.GetById(c.Param("id"))

	if err != nil {
		zap.S().Errorw("Error while getting domain, Reason: ", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusNotFound, utils.ErrorResponse{Message: "Domain not found"})
	} else {
		serializer := NewDomainSerializer()
		c.JSON(http.StatusOK, serializer.Serialize(domain))
	}
}

var DomainDetailView = auth.RoleRequired([]string{"admin", "superadmin"}, domainDetailView)

// Creates a domain
// @Summary Create domain
// @Description Creates a domain
// @Security BearerAuth
// @Tags domains
// @Accept  json
// @Produce  json
// @Param domain body DomainValidatorData true "Domain data"
// @Success 201 {object} DomainData
// @Failure 403 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Router /domain/ [post]
func createDomainView(c *gin.Context) {
	domainValidator := NewDomainValidator()
	if err := domainValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse{Message: err.Error()})
		return
	}

	if _, err := domainValidator.domain.Save(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse{Message: fmt.Sprintf("Cannot insert user: %v", err)})
		return
	}
	serializer := NewDomainSerializer()
	c.JSON(http.StatusCreated, serializer.Serialize(&domainValidator.domain))
}

var CreateDomainView = auth.RoleRequired([]string{"admin", "superadmin"}, createDomainView)

// Updates a domain
// @Summary Update domain
// @Description Updates a domain
// @Security BearerAuth
// @Tags domains
// @Accept  json
// @Produce  json
// @Param id path string true "Domain ID"
// @Param user body DomainValidatorData true "Domain data"
// @Success 200 {object} DomainData
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Router /domain/{id} [put]
func updateDomainView(c *gin.Context) {
	domainService := new(DomainService)
	domain, err := domainService.GetById(c.Param("id"))

	if err != nil {
		zap.S().Errorw("Error while getting domain, Reason: ", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusNotFound, utils.ErrorResponse{Message: "Domain not found"})
		return
	}

	domainValidator := NewDomainValidator()
	if err := domainValidator.BindUpdate(domain, c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse{Message: err.Error()})
		return
	}

	if _, err := domainValidator.domain.Save(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.ErrorResponse{Message: fmt.Sprintf("Cannot update domain: %v", err)})
		return
	}
	serializer := NewDomainSerializer()
	c.JSON(http.StatusOK, serializer.Serialize(&domainValidator.domain))
}

var UpdateDomainView = auth.RoleRequired([]string{"admin", "superadmin"}, updateDomainView)

// Deletes a domain
// @Summary Delete domain
// @Description Deletes a domain
// @Security BearerAuth
// @Tags domains
// @Accept  json
// @Produce  json
// @Param id path string true "Domain ID"
// @Success 204
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /domain/{id} [delete]
func deleteDomainView(c *gin.Context) {
	domainService := new(DomainService)
	domain, err := domainService.GetById(c.Param("id"))

	if err != nil {
		zap.S().Errorw("Error while getting domain, Reason: ", "id", c.Param("id"), "error", err)
		c.JSON(http.StatusNotFound, utils.ErrorResponse{Message: "Domain not found"})
	} else {
		if _, err := domain.Delete(); err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse{Message: err.Error()})
		}
		c.JSON(http.StatusNoContent, gin.H{})
	}
}

var DeleteDomainView = auth.RoleRequired([]string{"admin", "superadmin"}, deleteDomainView)
