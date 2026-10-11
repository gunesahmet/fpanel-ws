package Command

import (
	"Defines"
	"Environments"
	"Gorm"
	"Helpers"
	"fmt"
	"strings"
)

var ipAddress []Gorm.IpAddress
var servers []Gorm.Servers
var domains []Gorm.Domains

func findWebServerIp(webServerId string) string {
	for _, server := range servers {
		if strings.Compare(webServerId, server.Id) == 0 {
			for _, ip := range ipAddress {
				if strings.Compare(ip.ServerId, server.Id) == 0 {
					return ip.IpAddress
				}
			}
		}
	}
	return ""
}

func StartScreenshots() {
	if !Helpers.DirExists(Env.PanelHttpdDir + "/htdocs/ss/") {
		Helpers.CreateDirectory(Env.PanelHttpdDir+"/htdocs/ss/", true)
		_, _ = Helpers.Exec("/bin/sh", "-c", fmt.Sprintf("chown -R fpanel:fpanel %s/htdocs/ss/", Env.PanelHttpdDir))
	}
	for {
		res := Defines.DbConn.DB.Find(&ipAddress, "type < ?", Environments.DOMAIN)
		if res.RowsAffected <= 0 {
			Helpers.HandleError(nil, 9001, fmt.Sprintf("Ip addresses not found"))
		}
		if Env.IsCloud {
			res = Defines.DbConn.DB.Find(&servers, "type != 'main'")
			if res.RowsAffected <= 0 {
				Helpers.HandleError(nil, 9001, fmt.Sprintf("Servers not found"))
			}
		} else {
			res = Defines.DbConn.DB.Find(&servers, "type = 'main'")
			if res.RowsAffected <= 0 {
				Helpers.HandleError(nil, 9001, fmt.Sprintf("Servers not found"))
			}
		}

		res = Defines.DbConn.DB.Find(&domains, "status = ?", Environments.ACTIVE)
		if res.RowsAffected <= 0 {
			Helpers.HandleError(nil, 9001, fmt.Sprintf("Servers not found"))
		}

		hostFilePath := Env.HostsFile
		origHostData := Helpers.ReadFile(hostFilePath)
		newHostData := make([]string, 0)
		for _, v := range origHostData {
			found := false
			for _, domain := range domains {
				if strings.Contains(v, domain.Name) {
					found = true
				}
			}
			if !found {
				newHostData = append(newHostData, v)
			}
		}

		for _, domain := range domains {
			ipAddr := findWebServerIp(domain.WebServerId)
			if len(ipAddr) > 0 {
				newHostData = append(newHostData, fmt.Sprintf("%s    %s", ipAddr, domain.Name))
			} else {
				Env.Logger.Debug(9999, "Domain WebServer ipAddrId is empty. name:%s web_server_id: %s", domain.Name, domain.WebServerId)
			}
		}

		var command string
		err := Helpers.WriteFile(strings.Join(newHostData, "\n"), hostFilePath)
		if err != nil {
			Helpers.HandleError(err, 9001, fmt.Sprintf("Hosts file write failed"))
		} else {
			var wkhtml string
			if strings.Compare(Env.SystemOS, Environments.WINDOWS) == 0 {
				wkhtml = Env.PanelCoreDir + "/sbin/wkhtmltoimage.exe"
			} else {
				wkhtml = Env.PanelCoreDir + "/sbin/wkhtmltoimage"
			}
			for _, domain := range domains {

				Env.Logger.Debug(9999, "Getting domain: %s screenshot", domain.Name)
				command = fmt.Sprintf("%s --width %s --height %s http://%s/ %s", wkhtml, "1024", "768", domain.Name,
					Env.PanelHttpdDir+"/htdocs/ss/"+domain.Id+".jpg")
				if strings.Compare(Env.SystemOS, Environments.WINDOWS) == 0 {
					_, _ = Helpers.Exec("cmd", "/c", command)
				} else {
					_, _ = Helpers.Exec("/bin/bash", "-c", command)
					_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("chown -R fpanel:fpanel %s", Env.PanelHttpdDir+"/htdocs/ss"))
				}
				Helpers.Sleep(500)
			}
		}
		_ = Helpers.ReadFile(hostFilePath)

		err = Helpers.WriteFile(strings.Join(origHostData, "\n"), hostFilePath)
		if err != nil {
			Helpers.HandleError(err, 9001, fmt.Sprintf("Hosts file write failed"))
		}
		Helpers.Sleep(60 * 60 * 1000)
	}
}
