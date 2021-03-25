package domains

import (
	"go.mongodb.org/mongo-driver/bson/primitive" // for BSON ObjectID
	"net"
)

// User the user model
type Domain struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `json:"name"`
	Owner      string             `json:"owner"`
	Registrant string             `json:"registrant"`
	LoginInfo  string             `json:"loginInfo"`
	Package    string             `json:"package"`
	Mx         bool               `json:"mx"`
	Ip         net.IP             `json:"ip"`
	ServerName string             `json:"serverName"`
	Notes      string             `json:"notes"`
	Created    int64              `json:"created"`
	Updated    int64              `json:"updated"`
}

func (self *Domain) Save() (bool, error) {
	domainService := new(DomainService) // @TODO factory method
	result, err := domainService.Save(self)
	return result, err
}

func (self *Domain) Delete() (bool, error) {
	domainService := new(DomainService) // @TODO factory method
	result, err := domainService.Delete(self)
	return result, err
}
