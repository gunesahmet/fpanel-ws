package Gorm

type ClientStatistics struct {
	Id          int64  `json:"id,omitempty"`
	AccountName string `json:"account_name,omitempty"`
	ClientId    string `json:"client_id,omitempty"`
	ParentId    string `json:"parent_id,omitempty"`
	Type        string `json:"type,omitempty"`
	Mail        int64  `json:"mail,omitempty"`
	MailMax     string  `json:"mail_max,omitempty"`
	Db          int64  `json:"db,omitempty"`
	DbMax       string  `json:"db_max,omitempty"`
	Httpdoc     int64  `json:"httpdoc,omitempty"`
	Disk        int64  `json:"disk,omitempty"`
	DiskMax     string  `json:"disk_max,omitempty"`
	Traffic     int64  `json:"traffic,omitempty"`
	TrafficMax  string  `json:"traffic_max,omitempty"`
	Logs        int64  `json:"logs,omitempty"`
	Status      int    `json:"status,omitempty"`
	Statistic   int    `json:"statistic,omitempty"`
	OverUse     int    `json:"over_use,omitempty"`
	CreatedDate string `json:"created_date,omitempty"`
	ExpireDate  string `json:"expire_date,omitempty"`
	LastInit    string `json:"last_init,omitempty"`
	UpdatedDate string `json:"updated_date,omitempty"`
}

func (ClientStatistics) TableName() string {
	return "client_statistics"
}
