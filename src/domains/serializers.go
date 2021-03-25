package domains

type domainSerializer struct{}

type DomainData struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Owner      string `json:"owner"`
	Registrant string `json:"registrant"`
	LoginInfo  string `json:"loginInfo"`
	Package    string `json:"package"`
	Mx         bool   `json:"mx"`
	Ip         string `json:"ip"`
	ServerName string `json:"serverName"`
	Notes      string `json:"notes"`
	Created    int64  `json:"created"`
	Updated    int64  `json:"updated"`
}

func NewDomainSerializer() *domainSerializer {
	return &domainSerializer{}
}

func (self *domainSerializer) Serialize(domain *Domain) DomainData {
	domainData := DomainData{
		ID:         domain.ID.Hex(),
		Name:       domain.Name,
		Owner:      domain.Owner,
		Registrant: domain.Registrant,
		LoginInfo:  domain.LoginInfo,
		Package:    domain.Package,
		Mx:         domain.Mx,
		Ip:         domain.Ip.String(),
		ServerName: domain.ServerName,
		Notes:      domain.Notes,
		Created:    domain.Created,
		Updated:    domain.Updated,
	}
	return domainData
}

func (self *domainSerializer) SerializeMany(domains *[]Domain) []DomainData {
	var res []DomainData
	res = make([]DomainData, 0)
	for _, domain := range *domains {
		res = append(res, self.Serialize(&domain))
	}
	return res
}
