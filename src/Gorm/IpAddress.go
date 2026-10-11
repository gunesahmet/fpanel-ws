package Gorm

type IpAddress struct {
	Id          string `json:"id,omitempty" gorm:"primaryKey"`
	Index       string `json:"index,omitempty"`
	IpAddress   string `json:"ip_address,omitempty"`
	OldData     string `json:"old_data,omitempty"`
	SslStatus   string `json:"ssl_status,omitempty"`
	Netmask     string `json:"netmask,omitempty"`
	Gateway     string `json:"gateway,omitempty"`
	Interface   string `json:"interface,omitempty"`
	Type        string `json:"type,omitempty"`
	Format      string `json:"format,omitempty"`
	IsDefault   string `json:"is_default,omitempty"`
	ServerId    string `json:"server_id,omitempty"`
	Init        string `json:"init,omitempty"`
	Process     string `json:"process,omitempty"`
	CreatedDate string `json:"created_date,omitempty"`
}

func (IpAddress) TableName() string {
	return "ip_address"
}
