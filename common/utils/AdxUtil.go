package utils

import (
	"context"
	"fmt"
	"github.com/Azure/azure-kusto-go/kusto"
	"github.com/Azure/azure-kusto-go/kusto/kql"
	"github.com/PurpleScorpion/go-sweet-json/jsonutil"
	"io"
	"log"
	"sweet-common/logger"
	"time"
)

type AdxUtil struct {
}

var (
	Client  *kusto.Client
	LogFlag = false
)

func SelectKQL(kql string, baseFulleCode string) []jsonutil.JSONObject {
	result := sendKQL(kql, baseFulleCode)
	return result
}

func sendKQL(sql string, baseFulleCode string) []jsonutil.JSONObject {
	if LogFlag {
		logger.Info("kql - %s", sql)
	}

	type NodeRec struct {
		// ID is the table's NodeId. We use the field tag here to instruct our client to convert NodeId to ID.
		ID int64 `kusto:"NodeId"`
		// CollectionTime is Go representation of the Kusto datetime type.
		CollectionTime time.Time
	}

	dbName := fmt.Sprintf("ems-%s", baseFulleCode)

	ctx := context.Background()
	bl := &kql.Builder{}
	bl = bl.AddUnsafe(sql)
	// Query our database table "systemNodes" for the CollectionTimes and the NodeIds.
	iter, err := Client.Query(ctx, dbName, kql.FromBuilder(bl))
	if err != nil {
		panic(err.Error())
	}

	var list []jsonutil.JSONObject
	for {
		row, Err := iter.Next()
		if Err != nil {
			if Err == io.EOF {
				break
			}
			log.Println("An exception occurred during the iteration process:", Err)
		}
		// 获取最后一条数据的字段值
		columnNames := row.ColumnNames()
		// 遍历每个列名，获取对应的值并输出
		jSONObject := jsonutil.NewJSONObject()
		for i, columnName := range columnNames {
			value := row.Values[i]
			jSONObject.FluentPut(columnName, value.String())
		}
		list = append(list, jSONObject)
	}
	defer iter.Stop()
	return list
}
