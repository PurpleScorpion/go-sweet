package utils

import (
	"context"
	"fmt"
	"github.com/Azure/azure-kusto-go/kusto"
	"github.com/Azure/azure-kusto-go/kusto/kql"
	"github.com/PurpleScorpion/go-sweet-json/jsonutil"
	"io"
	"log"
	"shared/logger"
	"strconv"
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

	adxDbName := getAdxDBName(baseFulleCode)

	ctx := context.Background()
	bl := &kql.Builder{}
	bl = bl.AddUnsafe(sql)
	// Query our database table "systemNodes" for the CollectionTimes and the NodeIds.
	iter, err := Client.Query(ctx, adxDbName, kql.FromBuilder(bl))
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

func KustoSelectRaw(kql string, fullBaseCode string) (map[string]interface{}, error) {
	result, err := sendKQLRaw(kql, fullBaseCode)
	return result, err
}

func sendKQLRaw(sql string, fullBaseCode string) (map[string]interface{}, error) {
	logger.Debug("sendKQLRaw - KQL: {}", sql)

	type NodeRec struct {
		// ID is the table's NodeId. We use the field tag here to instruct our client to convert NodeId to ID.
		ID int64 `kusto:"NodeId"`
		// CollectionTime is Go representation of the Kusto datetime type.
		CollectionTime time.Time
	}

	adxDbName := getAdxDBName(fullBaseCode)

	ctx := context.Background()
	bl := &kql.Builder{}
	bl = bl.AddUnsafe(sql)
	// Query our database table "systemNodes" for the CollectionTimes and the NodeIds.
	iter, err := Client.Query(ctx, adxDbName, kql.FromBuilder(bl))
	if err != nil {
		return nil, err
	}

	var columns []map[string]interface{}
	var rows [][]interface{}
	columnsInitialized := false

	for {
		row, Err := iter.Next()
		if Err != nil {
			if Err == io.EOF {
				break
			}
			log.Println("An exception occurred during the iteration process:", Err)
		}

		// 初始化列信息（只需要执行一次）
		if !columnsInitialized {
			columnNames := row.ColumnNames()
			columnTypes := row.ColumnTypes
			for i, columnName := range columnNames {
				column := make(map[string]interface{})
				column["ColumnName"] = columnName
				if i < len(columnTypes) {
					column["ColumnType"] = fmt.Sprintf("%v", columnTypes[i].Type)
				} else {
					column["ColumnType"] = "unknown"
				}
				columns = append(columns, column)
			}
			columnsInitialized = true
		}

		// 添加行数据
		var rowData []interface{}
		for _, value := range row.Values {
			// 先尝试将字符串转换为数值，如果失败则保持原字符串
			valueStr := value.String()

			// 尝试转换为 float64
			if floatVal, err := strconv.ParseFloat(valueStr, 64); err == nil {
				rowData = append(rowData, floatVal)
			} else {
				// 如果不是数值，则保持原字符串
				rowData = append(rowData, valueStr)
			}
		}
		rows = append(rows, rowData)
	}

	result := map[string]interface{}{
		"Columns": columns,
		"Rows":    rows,
	}

	defer iter.Stop()
	return result, nil
}

func getAdxDBName(fullBaseCode string) string {
	return "dbname"
}
