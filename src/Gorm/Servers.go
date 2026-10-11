package Gorm

type Servers struct {
	Id            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	IsCore        int    `json:"is_core,omitempty"`
	Os            string `json:"os,omitempty"`
	Type          string `json:"type,omitempty"`
	IpAddress     string `json:"ip_address,omitempty"`
	DomainLimit   string `json:"domain_limit,omitempty"`
	DatabaseLimit string `json:"database_limit,omitempty"`
	MailLimit     string `json:"mail_limit,omitempty"`
	Data          string `json:"data,omitempty"`
	CreatedDate   string `json:"created_date,omitempty"`
}

type ServerData struct {
	Id            string     `json:"id,omitempty"`
	Name          string     `json:"name,omitempty"`
	Description   string     `json:"description,omitempty"`
	IsCore        int        `json:"is_core,omitempty"`
	Os            string     `json:"os,omitempty"`
	Type          string     `json:"type,omitempty"`
	IpAddress     string     `json:"ip_address,omitempty"`
	DomainLimit   string     `json:"domain_limit,omitempty"`
	DatabaseLimit string     `json:"database_limit,omitempty"`
	MailLimit     string     `json:"mail_limit,omitempty"`
	Data          string     `json:"data,omitempty"`
	CreatedDate   string     `json:"created_date,omitempty"`
	SecretData    SecretData `json:"secret_data,omitempty"`
}

type SecretData struct {
	Key       string `json:"key"`
	Iv        string `json:"iv"`
	DefaultIp string `json:"default_ip,omitempty"`
	SqlPass   string `json:"sql_pass,omitempty"`
	SqlSecret string `json:"sql_secret,omitempty"`
}

func (Servers) TableName() string {
	return "servers"
}
