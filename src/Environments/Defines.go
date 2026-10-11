package Environments

const (
	ROOT     int = 0
	ADMIN    int = 1
	RESELLER int = 2
	CLIENT   int = 3
	DOMAIN   int = 4
	EMAIL    int = 5
	SUPPORT  int = 6

	TRAFFIC_LIMIT_EXCEEDED  int = 2
	DISK_LIMIT_EXCEEDED     int = 3
	MANUAL_PASSIVE_BY_ADMIN int = 4
	MANUAL_PASSIVE_BY_OWNER int = 5
	ERROR_STATUS            int = 6
	EXPIRED            int = 7

	UNLIMITED string = "unlimited"

	ACTIVE  int = 1
	PASSIVE int = 0

	MYSQL   int = 1
	MARIADB int = 2
	PGSQL   int = 3
	MSSQL   int = 4
)
