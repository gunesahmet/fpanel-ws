package Defines

import (
	"Crypt"
	"Environments"
	"Gorm"
	"Helpers"
	"Logger"
	"encoding/json"
	"fmt"
)

var Env *Environments.Env

var DbConn Gorm.Connection

var Logging Logger.Interface

type TcpServerMessage struct {
	Command    string            `json:"command,omitempty"`
	Connection string            `json:"connection,omitempty"`
	Proc       string            `json:"proc,omitempty"`
	Data       map[string]string `json:"data,omitempty"`
}

func ConnectDatabase() error {
	Gorm.Env = Env
	err := Gorm.Connect()
	DbConn = Gorm.Conn
	return err
}

func GetIpServer(ipAddr string) (bool, Gorm.ServerData) {
	//fmt.Println(fmt.Sprintf("Searching in servers: %s", ipAddr))
	ipAddr = Helpers.Trim(ipAddr)
	var ipAddress Gorm.IpAddress
	res := DbConn.DB.First(&ipAddress, "ip_address = ?", ipAddr)

	if res.RowsAffected > 0 && len(ipAddress.ServerId) > 0 {
		return GetServer(ipAddress.ServerId)
	} else {
		Helpers.HandleError(nil, 1007, fmt.Sprintf("Ip address not found: %s", ipAddr))
	}
	return false, Gorm.ServerData{}
}

func GetServer(serverId string) (bool, Gorm.ServerData) {
	var serverModel Gorm.Servers
	var serverData Gorm.ServerData
	var secretData Gorm.SecretData
	res := DbConn.DB.First(&serverModel, "id = ?", serverId)
	if res.RowsAffected > 0 && len(serverModel.Name) > 0 {
		Logging.Info(1007, "Server Founded: %s", serverModel.Name)
		crypt, _ := Crypt.NewCrypter([]byte(Env.CliUsername), []byte(Env.CliPassword))
		decodedData, _ := crypt.Decode(serverModel.Data)
		err := json.Unmarshal([]byte(Helpers.Trim(decodedData)), &secretData)
		if err != nil {
			Helpers.HandleError(err, 1007, "SecretData decode failed")
		} else {
			js, _ := json.Marshal(serverModel)
			_ = json.Unmarshal(js, &serverData)
			serverData.SecretData = secretData
		}
		//fmt.Println(fmt.Sprintf("Server Key: %s", serverData.SecretData.Key))
		//fmt.Println(fmt.Sprintf("Server Iv:  %s", serverData.SecretData.Iv))

		return true, serverData
	} else {
		Helpers.HandleError(nil, 1007,
			fmt.Sprintf("Server not found. ServerID: %s", serverId))
	}

	return true, Gorm.ServerData{}
}
