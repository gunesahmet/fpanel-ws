package Gorm

type FtpUsers struct {
	Id   string `json:"id,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	CreatedDate string `json:"created_date,omitempty"`
	UpdatedDate string `json:"updated_date,omitempty"`
}

func (FtpUsers) TableName() string {
	return "ftp_users"
}
