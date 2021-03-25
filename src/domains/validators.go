package domains

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"time"
)

type DomainValidatorData struct {
	Name       string `json:"name" binding:"required"`
	Owner      string `json:"owner" binding:"required"`
	Registrant string `json:"registrant" binding:"required"`
	LoginInfo  string `json:"loginInfo"`
	Package    string `json:"package"`
	Mx         bool   `json:"mx"`
	Ip         string `json:"ip,omitempty" binding:"ip"`
	ServerName string `json:"serverName"`
	Notes      string `json:"notes"`
}
type DomainValidator struct {
	DomainData DomainValidatorData `json:"domain"`
	domain     Domain              `json:"-"`
}

func (self *DomainValidator) fillModelData() {
	self.domain.Name = self.DomainData.Name
	self.domain.Owner = self.DomainData.Owner
	self.domain.Registrant = self.DomainData.Registrant
	self.domain.LoginInfo = self.DomainData.LoginInfo
	self.domain.Package = self.DomainData.Package
	self.domain.Mx = self.DomainData.Mx
	self.domain.Ip = net.ParseIP(self.DomainData.Ip)
	self.domain.ServerName = self.DomainData.ServerName
	self.domain.Notes = self.DomainData.Notes
}

func (self *DomainValidator) Bind(c *gin.Context) error {
	err := c.ShouldBind(&self.DomainData)
	if err != nil {
		zap.S().Debug("Domain Validation Error: ", err)
		return err
	}
	self.fillModelData()
	self.domain.Created = time.Now().Unix()
	self.domain.Updated = time.Now().Unix()

	return nil
}

func (self *DomainValidator) BindUpdate(domain *Domain, c *gin.Context) error {
	err := c.ShouldBind(&self.DomainData)
	if err != nil {
		zap.S().Debug("Domain Validation Error: ", err)
		return err
	}
	self.domain.ID = domain.ID
	self.fillModelData()
	self.domain.Updated = time.Now().Unix()

	return nil
}

// You can put the default value of a Validator here
func NewDomainValidator() DomainValidator {
	domainValidator := DomainValidator{}
	return domainValidator
}
