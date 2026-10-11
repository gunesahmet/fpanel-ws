package Gorm

type Users struct {
	Id                 string           `json:"id,omitempty"`
	Username           string           `json:"username,omitempty"`
	Password           string           `json:"password,omitempty"`
	PasswordLastSet    string           `json:"password_last_set,omitempty"`
	LastLogin          string           `json:"last_login,omitempty"`
	LastLoginIp        string           `json:"last_login_ip,omitempty"`
	Parent             string           `json:"parent,omitempty"`
	Type               int              `json:"type,omitempty"`
	MustChangePassword int              `json:"must_change_password,omitempty"`
	CreatedDate        string           `json:"created_date,omitempty"`
	UpdatedDate        string           `json:"updated_date,omitempty"`
	ExpireDate         string           `json:"expire_date,omitempty"`
	StatusMsg          string           `json:"status_msg,omitempty"`
	Status             int              `json:"status,omitempty"`
	Interface          string           `json:"interface,omitempty"`
	Language           string           `json:"language,omitempty"`
	ApiStatus          int              `json:"api_status,omitempty"`
	ChangedData        string           `json:"changed_data,omitempty"`
	Locked             int              `json:"locked,omitempty"`
	Domains            []Domains        `json:"domains,omitempty" gorm:"foreignKey:ClientId"`
	Parents            *Users           `json:"parents,omitempty" gorm:"foreignKey:Parent"`
	ClientStatistics   ClientStatistics `json:"client_statistics,omitempty" gorm:"foreignKey:ClientId"`
}

func (Users) TableName() string {
	return "users"
}
