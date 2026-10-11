package Gorm

type Domains struct {
	Id               string           `json:"id,omitempty"`
	Name             string           `json:"name,omitempty"`
	ClientId         string           `json:"client_id,omitempty"`
	SslCertId        string           `json:"ssl_cert_id,omitempty"`
	FtpUserId        string           `json:"ftp_user_id,omitempty"`
	IpaddrId         string           `json:"ipaddr_id,omitempty"`
	Status           int              `json:"status,omitempty"`
	StatusMsg        string           `json:"status_msg,omitempty"`
	CreatedDate      string           `json:"created_date,omitempty"`
	ExpireDate       string           `json:"expire_date,omitempty"`
	UpdatedDate      string           `json:"updated_date,omitempty"`
	WebServerId      string           `json:"web_server_id,omitempty"`
	MailServerId     string           `json:"mail_server_id,omitempty"`
	DnsServerId      string           `json:"dns_server_id,omitempty"`
	Params           string           `json:"params,omitempty"`
	Prefix           string           `json:"prefix,omitempty"`
	DomainStatistics DomainStatistics `json:"domain_statistics,omitempty" gorm:"foreignKey:DomainId"`
	FtpUsers         FtpUsers         `json:"ftp_users,omitempty" gorm:"foreignKey:FtpUserId"`
	Users            Users            `json:"users,omitempty" gorm:"foreignKey:Username;references:Name"`
	WebServer        Servers          `json:"web_server,omitempty" gorm:"foreignKey:Id;references:WebServerId"`
	MailServer       Servers          `json:"mail_server,omitempty" gorm:"foreignKey:Id;references:MailServerId"`
	DnsServer        Servers          `json:"dns_server,omitempty" gorm:"foreignKey:Id;references:DnsServerId"`
}

func (Domains) TableName() string {
	return "domains"
}
