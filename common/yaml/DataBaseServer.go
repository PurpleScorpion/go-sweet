package sweetyml

import (
	"fmt"
	"github.com/Azure/azure-kusto-go/kusto"
	"github.com/PurpleScorpion/go-sweet-orm/mapper"
	"github.com/beego/beego/orm"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
	"go-sweet/common/utils"
)

var (
	kcsb *kusto.ConnectionStringBuilder
)

func initMySQL() {
	conf := GetYmlConf()
	if !conf.Sweet.MySqlConfig.Active {
		return
	}
	logs.Info("Init MySQL....")
	username := conf.Sweet.MySqlConfig.User
	password := conf.Sweet.MySqlConfig.Password
	host := conf.Sweet.MySqlConfig.Host
	port := conf.Sweet.MySqlConfig.Port
	dbName := conf.Sweet.MySqlConfig.DbName

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", username, password, host, port, dbName)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", connStr)
	orm.SetMaxIdleConns("default", conf.Sweet.MySqlConfig.MaxIdleConns)
	orm.SetMaxOpenConns("default", conf.Sweet.MySqlConfig.MaxOpenConns)
	orm.Debug = false
	mapper.InitMapper(mapper.MySQL, true)
}

func initAdx() {
	conf := GetYmlConf()
	if !conf.Sweet.Adx.Active {
		return
	}
	logs.Info("Init Adx....")
	if conf.Sweet.Adx.AuthMethod == "SMI" {
		kcsb = kusto.NewConnectionStringBuilder(conf.Sweet.Adx.Host).WithSystemManagedIdentity()
	} else {
		kcsb = kusto.NewConnectionStringBuilder(conf.Sweet.Adx.Host).
			WithAadAppKey(conf.Sweet.Adx.AppId, conf.Sweet.Adx.AppKey, conf.Sweet.Adx.AuthorityID)
	}
	utils.LogFlag = conf.Sweet.Adx.LogActive

	var err error
	utils.Client, err = kusto.New(kcsb)
	if err != nil {
		panic("add error handling")
	}
	defer utils.Client.Close()
}
