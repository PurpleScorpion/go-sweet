package system

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"shared/constants"
	"shared/logger"
	"shared/utils"

	"github.com/Azure/azure-kusto-go/kusto"
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"github.com/PurpleScorpion/go-sweet-orm/v3/mapper"
	_ "github.com/go-sql-driver/mysql"
)

var (
	kcsb *kusto.ConnectionStringBuilder
)

func initMySQL() {
	host := keqing.ValueString("${sweet.mysql.host}")
	if keqing.IsEmpty(host) {
		return
	}
	logger.Info("Init MySQL....")
	username := keqing.ValueString("${sweet.mysql.user}")
	if keqing.IsEmpty(username) {
		panic("mysql username is empty")
	}
	pwd := keqing.ValueObject("${sweet.mysql.password}")
	password := ""
	switch pwd.(type) {
	case int:
		password = fmt.Sprintf("%d", pwd.(int))
	case string:
		password = pwd.(string)
	}

	if keqing.IsEmpty(password) {
		panic("mysql password is empty")
	}
	port := keqing.ValueInt("${sweet.mysql.port}")
	if port == 0 {
		port = 3306
	}
	if port <= 0 || port > 65535 {
		panic("mysql port is invalid")
	}
	dbName := keqing.ValueString("${sweet.mysql.dbName}")
	if keqing.IsEmpty(dbName) {
		panic("mysql dbName is empty")
	}

	maxIdleConns := keqing.ValueInt("${sweet.mysql.maxIdleConns}")
	maxOpenConns := keqing.ValueInt("${sweet.mysql.maxOpenConns}")

	useSSL := keqing.ValueBool("${sweet.mysql.useSSL}")

	conf := mapper.MySQLConf{
		UserName:    username,
		Password:    password,
		DbName:      dbName,
		Port:        port,
		Host:        host,
		MaxIdleConn: maxIdleConns,
		MaxOpenConn: maxOpenConns,
	}

	if useSSL {
		profilesActive := os.Getenv(constants.PROFILES_ACTIVE)
		pemName := keqing.ValueString("${sweet.mysql.pem}")
		if keqing.IsEmpty(pemName) {
			panic("mysql ssl pem is empty")
		}

		pemPath := ""
		if keqing.IsEmpty(profilesActive) {
			// 如果docker环境变量为空, 则认为是本地环境
			pemPath = "src/main/resources/" + pemName
		} else {
			// 否则是docker环境
			pemPath = "system/" + pemName
		}

		rootCertPool := x509.NewCertPool()
		pem, _ := ioutil.ReadFile(pemPath)
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			panic("Failed to append PEM.")
		}
		conf.TlsCertPool = rootCertPool
	}

	mapper.SetMySqlConf(conf)
	mapper.RegisterMySql()
}

func initAdx() {
	host := keqing.ValueString("${sweet.adx.host}")
	if keqing.IsEmpty(host) {
		return
	}
	logger.Info("Init Adx....")

	authMethod := keqing.ValueString("${sweet.adx.authMethod}")
	if keqing.IsEmpty(authMethod) {
		authMethod = "AAK"
	}
	appId := ""
	appKey := ""
	authorityID := ""
	if authMethod == "AAK" {
		appId = keqing.ValueString("${sweet.adx.appId}")
		if keqing.IsEmpty(appId) {
			panic("adx appId is empty")
		}
		appKey = keqing.ValueString("${sweet.adx.appKey}")
		if keqing.IsEmpty(appKey) {
			panic("adx appKey is empty")
		}
		authorityID = keqing.ValueString("${sweet.adx.authorityID}")
		if keqing.IsEmpty(authorityID) {
			panic("adx authorityID is empty")
		}
	}
	if authMethod == "SMI" {
		kcsb = kusto.NewConnectionStringBuilder(host).WithSystemManagedIdentity()
	} else {
		kcsb = kusto.NewConnectionStringBuilder(host).WithAadAppKey(appId, appKey, authorityID)
	}
	utils.LogFlag = keqing.ValueBool("${sweet.adx.logActive}")

	var err error
	utils.Client, err = kusto.New(kcsb)
	if err != nil {
		panic("add error handling")
	}
	defer utils.Client.Close()
}
