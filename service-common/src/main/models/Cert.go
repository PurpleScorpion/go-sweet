package models

import "time"

type Cert struct {
	Id           int32     `json:"id" tableId:"id"`
	OriginalName string    `json:"originalName"`
	FileName     string    `json:"fileName"`
	FileMd5      string    `json:"fileMd5"`
	FileSize     int32     `json:"fileSize"`
	DeleteFlag   int32     `json:"deleteFlag"`
	CreatedBy    string    `json:"createdBy"`
	CreatedTime  time.Time `json:"createdTime"`
}

func (Cert) TableName() string {
	return "tbl_cert"
}
