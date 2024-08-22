package sweetyml

import (
	"fmt"
	"github.com/Azure/azure-kusto-go/kusto"
	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"github.com/PurpleScorpion/go-sweet-orm/mapper"
	"github.com/beego/beego/orm"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
	"sweet-common/utils"
)

var (
	kcsb *kusto.ConnectionStringBuilder
)

func initMySQL() {
	host := keqing.ValueString("${sweet.mysql.host}")
	if keqing.IsEmpty(host) {
		return
	}
	logs.Info("Init MySQL....")
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

	if maxIdleConns == 0 {
		maxIdleConns = 50
	}
	if maxOpenConns == 0 {
		maxOpenConns = 100
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", username, password, host, port, dbName)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", connStr)
	orm.SetMaxIdleConns("default", maxIdleConns)
	orm.SetMaxOpenConns("default", maxOpenConns)
	orm.Debug = false
	mapper.InitMapper(mapper.MySQL, true)
}

func initAdx() {
	host := keqing.ValueString("${sweet.adx.host}")
	if keqing.IsEmpty(host) {
		return
	}
	logs.Info("Init Adx....")

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
	fmt.Println(appId, appKey, authorityID)
	if authMethod == "SMI" {
		kcsb = kusto.NewConnectionStringBuilder(host).WithSystemManagedIdentity()
	} else {
		kcsb = kusto.NewConnectionStringBuilder(host).
			WithAadAppKey(appId, appKey, authorityID)
	}
	utils.LogFlag = keqing.ValueBool("${sweet.adx.logActive}")

	var err error
	utils.Client, err = kusto.New(kcsb)
	if err != nil {
		panic("add error handling")
	}
	defer utils.Client.Close()
}
