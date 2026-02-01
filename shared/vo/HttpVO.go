package vo

import "mime/multipart"

type HttpRespVO struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
}

type UploadForm struct {
	Name string                `form:"name"`
	File *multipart.FileHeader `form:"file" binding:"required"`
}
