package Environments

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	"Helpers"
	"Logger"
	"github.com/go-ini/ini"
)

const (
	LINUX   = "linux"
	WINDOWS = "windows"
)

type Env struct {
	SystemType           string `json:"system_type"`
	SystemVersion        string `json:"system_version"`
	PanelDir             string `json:"panel_dir"`
	PanelHttpdDir        string `json:"panel_httpd_dir"`
	PanelCoreDir         string `json:"panel_core_dir"`
	PanelLogDir          string `json:"panel_log_dir"`
	PanelLog             string `json:"panel_log"`
	PhpDir               string `json:"php_dir"`
	PhpIni               string `json:"php_ini"`
	ApacheDir            string `json:"apache_dir"`
	BindDir              string `json:"bind_dir"`
	VhostsDir            string `json:"vhosts_dir"`
	ApacheInit           string `json:"apache_init"`
	ApacheBin            string `json:"apache_bin"`
	RpmOsFile            string `json:"rpm_os_file"`
	DebOsFile            string `json:"deb_os_file"`
	RpmRepoFile          string `json:"rpm_repo_file"`
	DebRepoFile          string `json:"deb_repo_file"`
	PanelConfig          string `json:"panel_config"`
	PanelEtcDir          string `json:"panel_etc_dir"`
	PanelReleaseConf     string `json:"panel_release_conf"`
	PanelHttpConfigIni   string `json:"panel_http_config_ini"`
	Sysconf              string `json:"sysconf"`
	BindZoneConf         string `json:"bind_zone_conf"`
	BindDomainTemplate   string `json:"bind_domain_template"`
	BindZoneDir          string `json:"bind_zone_dir"`
	DomainsConfigDir     string `json:"domains_config_dir"`
	NetworkConf          string `json:"network_conf"`
	HostnameConf         string `json:"hostname_conf"`
	HostsFile            string `json:"hosts_file"`
	LogRotateVhost       string `json:"log_rotate_vhost"`
	LogRotateConf        string `json:"log_rotate_conf"`
	ProftpConf           string `json:"proftp_conf"`
	MyCnf                string `json:"my_cnf"`
	CrontabFile          string `json:"crontab_file"`
	CrontabDir           string `json:"crontab_dir"`
	InitDir              string `json:"init_dir"`
	MysqlDir             string `json:"mysql_dir"`
	MailDir              string `json:"mail_dir"`
	AppsDir              string `json:"apps_dir"`
	RepoServer           string `json:"repo_server"`
	AppsServer           string `json:"apps_server"`
	UpdateServer         string `json:"update_server"`
	Sqlite               string `json:"sqlite"`
	MailDb               string `json:"mail_db"`
	MailDomainDir        string `json:"mail_domain_dir"`
	MimeTypesConfig      string `json:"mime_types_config"`
	WebStatisticsCli     string `json:"web_statistics_cli"`
	WebStatisticsTool    string `json:"web_statistics_tool"`
	FpanelStatisticsTool string `json:"fpanel_statistics_tool"`
	SystemPackageConf    string `json:"system_package_conf"`
	GeneralLogPath       string `json:"general_log_path"`

	UserPasswd           string   `json:"user_passwd"`
	Useradd              string   `json:"useradd"`
	Userdel              string   `json:"userdel"`
	Usermod              string   `json:"usermod"`
	Groupadd             string   `json:"groupadd"`
	Groupdel             string   `json:"groupdel"`
	Groupmod             string   `json:"groupmod"`
	UsertmpmailDir       string   `json:"usertmpmail_dir"`
	Usertmpfile          string   `json:"usertmpfile"`
	Ifconfig             string   `json:"ifconfig"`
	Route                string   `json:"route"`
	Ntpdate              string   `json:"ntpdate"`
	Ntpserver            string   `json:"ntpserver"`
	NoLogin              string   `json:"no_login"`
	BackupExt            string   `json:"backup_ext"`
	SpamStartup          string   `json:"spam_startup"`
	SpamInit             string   `json:"spam_init"`
	SpamLocal            string   `json:"spam_local"`
	PolicyConfig         string   `json:"policy_config"`
	PolicyDb             string   `json:"policy_db"`
	SpamWhiteLocal       string   `json:"spam_white_local"`
	SpamBlackLocal       string   `json:"spam_black_local"`
	SpamClient           string   `json:"spam_client"`
	SpamDaemon           string   `json:"spam_daemon"`
	Vdelivermail         string   `json:"vdelivermail"`
	Vpasswd              string   `json:"vpasswd"`
	VmailUser            string   `json:"vmail_user"`
	VmailGroup           string   `json:"vmail_group"`
	Netcat               string   `json:"netcat"`
	SevenZip             string   `json:"seven_zip"`
	Adminmng             string   `json:"adminmng"`
	Appinstaller         string   `json:"appinstaller"`
	Backupclientmng      string   `json:"backupclientmng"`
	Backupmng            string   `json:"backupmng"`
	Backupsrvmng         string   `json:"backupsrvmng"`
	Compile              string   `json:"compile"`
	Crontabmng           string   `json:"crontabmng"`
	Datemng              string   `json:"datemng"`
	Directorymng         string   `json:"directorymng"`
	Dnsmng               string   `json:"dnsmng"`
	Domaliasmng          string   `json:"domaliasmng"`
	Fixmng               string   `json:"fixmng"`
	Ftpmng               string   `json:"ftpmng"`
	Filemng              string   `json:"filemng"`
	Logmng               string   `json:"logmng"`
	Mailmng              string   `json:"mailmng"`
	Mailqueuemng         string   `json:"mailqueuemng"`
	Cloudmng             string   `json:"cloudmng"`
	Mailsrvmng           string   `json:"mailsrvmng"`
	Migrationmng         string   `json:"migrationmng"`
	Networkmng           string   `json:"networkmng"`
	Packagemng           string   `json:"packagemng"`
	Pagemng              string   `json:"pagemng"`
	Phpmng               string   `json:"phpmng"`
	Protectmng           string   `json:"protectmng"`
	Servicesmng          string   `json:"servicesmng"`
	Sqlmng               string   `json:"sqlmng"`
	Sslmng               string   `json:"sslmng"`
	Subdirectorymng      string   `json:"subdirectorymng"`
	Subusermng           string   `json:"subusermng"`
	Subvhostmng          string   `json:"subvhostmng"`
	Usermng              string   `json:"usermng"`
	Vhostmng             string   `json:"vhostmng"`
	Websrvmng            string   `json:"websrvmng"`
	Webstatistics        string   `json:"webstatistics"`
	Statisticsmng        string   `json:"statisticsmng"`
	Wrapper              string   `json:"wrapper"`
	Fcrypt               string   `json:"fcrypt"`
	PanelTmpDir          string   `json:"panel_tmp_dir"`
	PanelBackupDir       string   `json:"panel_backup_dir"`
	PanelBackupTmpDir    string   `json:"panel_backup_tmp_dir"`
	PhpFpmRunDir         string   `json:"php_fpm_run_dir"`
	Nogroup              string   `json:"nogroup"`
	MailSystemUser       string   `json:"mail_system_user"`
	MailSystemGroup      string   `json:"mail_system_group"`
	NginxDir             string   `json:"nginx_dir"`
	NginxDefaultConf     string   `json:"nginx_default_conf"`
	NginxBin             string   `json:"nginx_bin"`
	PhpFpmDir            string   `json:"php_fpm_dir"`
	PhpFpmWwwConf        string   `json:"php_fpm_www_conf"`
	PhpFpmConf           string   `json:"php_fpm_conf"`
	PhpFpmVhostConf      string   `json:"php_fpm_vhost_conf"`
	CliUsername          string   `json:"cli_username"`
	CliPassword          string   `json:"cli_password"`
	FpanelAdminPasswd    string   `json:"fpanel_admin_passwd"`
	InstallCommand       string   `json:"install_command"`
	ManualInstallCommand string   `json:"manual_install_command"`
	PackageExtension     string   `json:"package_extension"`
	HttpdServer          string   `json:"httpd_server"`
	SmtpServer           string   `json:"smtp_server"`
	ImapServer           string   `json:"imap_server"`
	DatabaseServers      string   `json:"database_servers"`
	SpamService          string   `json:"spam_service"`
	Statistics           string   `json:"statistics"`
	DnsServer            string   `json:"dns_server"`
	FtpServer            string   `json:"ftp_server"`
	Debug                bool     `json:"debug"`
	SystemOS             string   `json:"system_os"`
	DIRECTORYSEPARATOR   string   `json:"directory_separator"`
	StaticKey            string   `json:"static_key"`
	RedisAuth            string   `json:"redis_auth"`
	CloudConfigLocation  string   `json:"cloud_config_location"`
	IsCloud              bool     `json:"is_cloud"`
	IsCloudMain          bool     `json:"is_cloud_main"`
	CloudClients         []string `json:"cloud_clients"`
	CloudMainServer      string   `json:"cloud_main_server"`
	Logger               Logger.Interface
}

var (
	env    *Env
	envMtx = new(sync.RWMutex)
)

var Logging Logger.Interface

func detectSystemTypeWithLsb() {
	_, lsbRelease := Helpers.Exec("lsb_release", "-i", "-s")
	env.SystemType = strings.ToLower(lsbRelease)
}

func detectSystemType() {
	var systemType string
	if Helpers.FileExists("/etc/os-release") {
		for _, eachline := range Helpers.ReadFile("/etc/os-release") {
			if strings.HasPrefix(eachline, "ID=") {
				systemType = strings.ToLower(Helpers.ReplaceAll(Helpers.RemoveAll(eachline, "ID="), "\"", ""))
				systemType = strings.TrimSpace(systemType)
				systemType = strings.Trim(systemType, "\r")
				systemType = strings.Trim(systemType, "\n")
			}
		}
	}

	if (strings.Compare(systemType, "centos") == 0) || (strings.Compare(systemType, "debian") == 0) || (strings.Compare(systemType, "ubuntu") == 0) {
		env.SystemType = systemType
	} else {
		detectSystemTypeWithLsb()
	}
}

func detectSystemVersion() {
	_, lsbRelease := Helpers.Exec("lsb_release", "-r", "-s")
	lsbReleaseArr := strings.Split(lsbRelease, ".")
	env.SystemVersion = lsbReleaseArr[0]
}

func Init() {
	systemOs := strings.ToLower(runtime.GOOS)
	systemOs = strings.Trim(systemOs, "\n")
	systemOs = strings.Trim(systemOs, "\r")
	systemOs = strings.TrimSpace(systemOs)
	env = new(Env)
	if strings.Compare(systemOs, "linux") == 0 {
		detectSystemType()
		detectSystemVersion()
	}
	env.SystemOS = systemOs
	env.DIRECTORYSEPARATOR = "/"
	if strings.Compare(systemOs, "linux") == 0 {
		env.PanelDir = "/usr/local/fpanel"
		env.PanelCoreDir = env.PanelDir + "/core"
		env.PanelLogDir = "/var/log/fpanel"
		env.PanelLog = "/var/log/fpanel/*.log"
		env.VhostsDir = "/var/www"
		env.HostsFile = "/etc/hosts"
		env.GeneralLogPath = "/var/log"

		if Helpers.FileExists(env.PanelCoreDir + "/bin/debug") {
			env.Debug = true
		}

		if os.Geteuid() != 0 {
			fmt.Println("Tool is running without root privileges.")
			os.Exit(1001)
		}
	} else {
		env.PanelDir = os.Getenv("SystemDrive") + "/fpanel"
		env.PanelCoreDir = env.PanelDir + "/core"
		env.PanelLogDir = env.PanelDir + "/var/log/fpanel"
		env.PanelLog = env.PanelDir + "/var/log/fpanel/*.log"
		env.VhostsDir = os.Getenv("SystemDrive") + "/vhosts"
		env.HostsFile = os.Getenv("windir") + "/system32/drivers/etc/hosts"
		env.HostnameConf = env.PanelDir + "/etc/hostname"
		env.GeneralLogPath = os.Getenv("windir") + "/Logs"
		if Helpers.FileExists(env.PanelCoreDir + "/bin/debug") {
			env.Debug = true
		}
	}

	if !Helpers.DirExists(env.PanelLogDir + "/websocket") {
		Helpers.CreateDirectory(env.PanelLogDir+"/websocket", true)
	}

	Logging = Logger.New(log.New(os.Stdin, "\r\n", log.LstdFlags), Logger.Config{
		LogLevel:                  Logger.Silent,
		IgnoreRecordNotFoundError: true,
		Colorful:                  true,
		PanelLogDir:               env.PanelLogDir + "/websocket/",
		Debugger:                  env.Debug,
	})

	if Helpers.FileExists(env.PanelDir + "/redis/redis.conf") {
		lines := Helpers.ReadFile(env.PanelDir + "/redis/redis.conf")
		for _, line := range lines {
			if strings.HasPrefix(line, "requirepass") {
				env.RedisAuth = Helpers.Trim(Helpers.RemoveAll(line, "requirepass"))
			}
		}
	}
	env.PanelTmpDir = env.PanelDir + "/tmp"
	env.PanelHttpdDir = env.PanelDir + "/httpd"

	env.PanelEtcDir = env.PanelDir + "/etc"
	env.PanelConfig = env.PanelEtcDir + "/panel.conf"
	env.DomainsConfigDir = env.PanelDir + "/configs/domains/"
	env.PanelBackupDir = env.PanelDir + "/backup"
	env.PanelBackupTmpDir = env.PanelDir + "/backup/.tmp"
	env.PanelReleaseConf = env.PanelDir + "/etc/release.conf"                // ???
	env.PanelHttpConfigIni = env.PanelHttpdDir + "/htdocs/config/config.ini" // ???
	env.Sysconf = env.PanelDir + "/etc/sys.conf"
	env.NetworkConf = env.PanelDir + "/etc/network.conf"
	env.WebStatisticsCli = env.PanelHttpdDir + "/htdocs/cli.php"
	env.WebStatisticsTool = env.PanelCoreDir + "/sbin/webstatistics"
	env.FpanelStatisticsTool = env.PanelCoreDir + "/sbin/statistics"
	env.SystemPackageConf = env.PanelDir + "/etc/.sys.packages"

	env.CloudConfigLocation = env.PanelDir + "/etc/cloud.config.ini"
	env.IsCloudMain = false
	env.IsCloud = false
	if Helpers.FileExists(env.CloudConfigLocation) {
		cloudConfigIni, err := ini.ShadowLoad(env.CloudConfigLocation)
		if err == nil {
			env.CloudMainServer = cloudConfigIni.Section("socket.server").Key("address").String()
		}
		env.IsCloud = true
	} else {
		httpConfigIni, err := ini.ShadowLoad(env.PanelHttpConfigIni)
		if err == nil {
			env.IsCloudMain = true
			env.CloudClients = httpConfigIni.Section("clients").Key("address[]").ValueWithShadows()
		}
		env.CloudMainServer = "127.0.0.1"
	}

	cfg, err := ini.Load(env.SystemPackageConf)
	if err != nil {
		Logging.Error(1004, "Fpanel system packages config doesn't exists.")
		os.Exit(1004)
	}
	env.HttpdServer = cfg.Section("").Key("HTTPD_SERVER").String()
	env.FtpServer = cfg.Section("").Key("FTP_SERVER").String()
	env.DnsServer = cfg.Section("").Key("DNS_SERVER").String()
	env.SmtpServer = cfg.Section("").Key("SMTP_SERVER").String()
	env.ImapServer = cfg.Section("").Key("IMAP_SERVER").String()
	env.DatabaseServers = cfg.Section("").Key("DATABASE_SERVERS").String()
	env.Statistics = cfg.Section("").Key("STATISTICS").String()
	env.SpamService = cfg.Section("").Key("SPAM_SERVICE").String()
	env.MailDir = env.PanelDir + "/mail"

	if strings.Compare(systemOs, "linux") == 0 {
		if (strings.Compare(env.SystemType, "debian") == 0) || (strings.Compare(env.SystemType, "ubuntu") == 0) {
			env.MyCnf = "/etc/mysql/mysql.conf.d/mysqld.cnf" // If debian
			env.ManualInstallCommand = "dpkg -i --force-overwrite "
			env.InstallCommand = "apt-get install -y -qq "
			env.PackageExtension = ".deb"
		} else if strings.Compare(env.SystemType, "centos") == 0 {
			env.MyCnf = "/etc/my.cnf" // If centos
			env.InstallCommand = "yum install -y "
			env.ManualInstallCommand = "rpm -ivh --force --nodeps "
			env.PackageExtension = ".rpm"
		} else {
			Logging.Error(1002, "Undetected os type.")
			os.Exit(1002)
		}
		env.HostnameConf = "/etc/hostname"
		env.LogRotateVhost = "/etc/logrotate.d/vhosts"
		env.LogRotateConf = "/etc/logrotate.conf"

		env.ApacheDir = "/etc/apache2"
		env.BindDir = "/etc/bind/"
		env.BindZoneConf = "/etc/bind/domains.conf"
		env.BindDomainTemplate = "/etc/bind/domains.template"
		env.BindZoneDir = "/etc/bind/domains/"
		env.ApacheInit = "/etc/init.d/apache2"
		env.ApacheBin = "/usr/sbin/apache2ctl"
		env.ProftpConf = "/etc/proftpd/proftpd.conf"
		env.SystemType = Helpers.RemoveAll(env.SystemType, "\"")
		env.CrontabFile = "/etc/crontab"
		env.CrontabDir = "/var/spool/cron"
		env.InitDir = "/etc/init.d"
		env.MysqlDir = "/var/lib/mysql"
		env.MimeTypesConfig = "/etc/mime.types"
		env.Sqlite = "sqlite3"
		env.MailDb = env.MailDir + "/etc/mail.db"
		env.MailDomainDir = env.MailDir + "/domains/"

		env.UserPasswd = "/etc/passwd"
		env.Useradd = "/usr/sbin/useradd"
		env.Userdel = "/usr/sbin/userdel"
		env.Usermod = "/usr/sbin/usermod"
		env.Groupadd = "/usr/sbin/groupadd"
		env.Groupdel = "/usr/sbin/groupdel"
		env.Groupmod = "/usr/sbin/groupmod"
		env.UsertmpmailDir = "/var/mail"
		env.Usertmpfile = "/var/spool/mail"
		env.Ifconfig = "/sbin/ifconfig"
		env.Route = "/sbin/route"
		env.Ntpdate = "/usr/sbin/ntpdate"
		env.Ntpserver = "0.asia.pool.ntp.org"
		env.NoLogin = "/sbin/nologin"

		env.SpamStartup = "/etc/sysconfig/spamassassin"
		env.SpamInit = env.InitDir + "/spamassassin"
		env.SpamLocal = "/etc/spamassassin/local.cf"
		env.SpamWhiteLocal = "/etc/spamassassin/whitelist.cf"
		env.SpamBlackLocal = "/etc/spamassassin/blacklist.cf"
		env.SpamClient = "/usr/bin/spamc"
		env.SpamDaemon = "/usr/bin/spamd"

		env.Vdelivermail = env.PanelDir + "/bin/vdelivermail"
		env.Vpasswd = env.PanelDir + "/bin/vpasswd"
		env.VmailUser = "vmail"
		env.VmailGroup = "vmail"

		env.PhpFpmRunDir = "/var/spool/phpfpm"
		env.Nogroup = "nogroup"
		env.MailSystemUser = "vmail"
		env.MailSystemGroup = "vmail"
		env.NginxDir = "/etc/nginx"
		env.NginxDefaultConf = "/etc/nginx/nginx.conf"
		env.NginxBin = "/usr/sbin/nginx"
		env.PhpFpmDir = "/etc/php5/fpm"
		env.PhpFpmWwwConf = env.PhpFpmDir + "/pool.d/www.conf"
		env.PhpFpmConf = env.PhpFpmDir + "/php-fpm.conf"
		env.PhpFpmVhostConf = env.PhpFpmDir + "/pool.d/zzz999_httpd.conf"

		if strings.Compare(env.SystemType, "centos") == 0 {
			env.ApacheDir = "/etc/httpd"
			env.ApacheInit = "/etc/init.d/httpd"
			env.ApacheBin = "/usr/sbin/httpd"
			env.Nogroup = "nobody"

			if Helpers.FileExists("/usr/bin/sqlite3") {
				env.Sqlite = "/usr/bin/sqlite3"
			} else {
				env.Sqlite = "/usr/bin/sqlite"
			}
		}
	} else {
		env.PackageExtension = ".msi"
		env.ManualInstallCommand = "c:\\windows\\system32\\msiexec.exe /i /qn "
		if strings.Compare(env.SmtpServer, "mail_enable") == 0 {
			env.MailDomainDir = env.MailDir + "/Postoffices/"
		} else if strings.Compare(env.SmtpServer, "hmail") == 0 {
			env.MailDomainDir = env.MailDir + "/Data/"
		}
	}
	env.AppsDir = env.PanelDir + "/apps"
	env.RepoServer = "repo.fpanel.net"
	env.AppsServer = "apps.fpanel.net"
	env.UpdateServer = "update.fpanel.net"
	env.PolicyConfig = env.MailDir + "/etc/policyd.conf"
	env.PolicyDb = env.MailDir + "/etc/policyd.db"

	env.BackupExt = ".zip"

	env.Adminmng = env.PanelCoreDir + "/bin/admninmng"
	env.Appinstaller = env.PanelCoreDir + "/bin/appinstaller"
	env.Backupmng = env.PanelCoreDir + "/bin/backupmng"
	env.Compile = env.PanelCoreDir + "/bin/compile"
	env.Crontabmng = env.PanelCoreDir + "/bin/crontabmng"
	env.Datemng = env.PanelCoreDir + "/bin/datemng"
	env.Directorymng = env.PanelCoreDir + "/bin/directorymng"
	env.Dnsmng = env.PanelCoreDir + "/bin/dnsmng"
	env.Domaliasmng = env.PanelCoreDir + "/bin/domaliasmng"
	env.Fixmng = env.PanelCoreDir + "/bin/fixmng"
	env.Ftpmng = env.PanelCoreDir + "/bin/ftpmng"
	env.Filemng = env.PanelCoreDir + "/bin/filemng"
	env.Logmng = env.PanelCoreDir + "/bin/logmng"
	env.Mailmng = env.PanelCoreDir + "/bin/mailmng"
	env.Mailqueuemng = env.PanelCoreDir + "/bin/mailqueuemng"
	env.Mailsrvmng = env.PanelCoreDir + "/bin/mailsrvmng"
	env.Cloudmng = env.PanelCoreDir + "/bin/cloudmng"
	env.Migrationmng = env.PanelCoreDir + "/bin/migrationmng"
	env.Networkmng = env.PanelCoreDir + "/bin/networkmng"
	env.Packagemng = env.PanelCoreDir + "/bin/packagemng"
	env.Pagemng = env.PanelCoreDir + "/bin/pagemng"
	env.Phpmng = env.PanelCoreDir + "/bin/phpmng"
	env.Protectmng = env.PanelCoreDir + "/bin/protectmng"
	env.Servicesmng = env.PanelCoreDir + "/bin/servicesmng"
	env.Sqlmng = env.PanelCoreDir + "/bin/sqlmng"
	env.Sslmng = env.PanelCoreDir + "/bin/sslmng"
	env.Subdirectorymng = env.PanelCoreDir + "/bin/subdirectorymng"
	env.Subusermng = env.PanelCoreDir + "/bin/subusermng"
	env.Subvhostmng = env.PanelCoreDir + "/bin/subvhostmng"
	env.Usermng = env.PanelCoreDir + "/bin/usermng"
	env.Vhostmng = env.PanelCoreDir + "/bin/vhostmng"
	env.Websrvmng = env.PanelCoreDir + "/bin/websrvmng"
	env.Webstatistics = env.PanelCoreDir + "/bin/webstatistics"
	env.Statisticsmng = env.PanelCoreDir + "/bin/statisticsmng"
	env.Wrapper = env.PanelCoreDir + "/bin/wrapper"
	env.Fcrypt = env.PanelCoreDir + "/sbin/fcrypt"
	env.StaticKey = "87f5e1c8c1dd33cfd9dbf2990b60d595"
	env.Logger = Logging

	if Helpers.FileExists(env.PanelConfig) {
		configStatus := false
		for _, line := range Helpers.ReadFile(env.PanelConfig) {
			if len(line) > 0 {
				keys := strings.Fields(line)
				env.CliUsername = keys[0]
				env.CliPassword = keys[1][0:16]
				env.FpanelAdminPasswd = keys[2]
				configStatus = true
				break
			}
		}

		if !configStatus {
			Logging.Error(1003, "Fpanel config doesn't exists.")
			os.Exit(1003)
		}
	} else {
		Logging.Error(1003, "Fpanel config doesn't exists.")
		os.Exit(1003)
	}

}

func GetConfig() *Env {
	envMtx.RLock()
	defer envMtx.RUnlock()
	return env
}
