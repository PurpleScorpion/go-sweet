package controllers

import (
	"github.com/PurpleScorpion/go-sweet-json/jsonutil"
	"io"
	"os"
	"sweet-common/constants"
	"sweet-common/utils"
)

type FileController struct {
	BaseController
}

func (that *FileController) GetBaseImg() {
	that.Ok("Success", "")
}

func (that *FileController) UploadImg() {

	files, _ := that.GetFiles("files")

	if len(files) == 0 {
		that.Error(constants.SYSTEM_ERROR, "Please select the file to upload", "")
	}

	fileHeader := files[0]

	fileName := ""

	if fileHeader != nil {

		fileName = utils.RandomImageName(fileHeader.Filename)

		savePath := getPath() + fileName

		src, err := fileHeader.Open()
		if err != nil {
			that.Error(constants.SYSTEM_ERROR, "Failed to open file", "")
			return
		}
		defer src.Close()

		dst, err := os.Create(savePath)
		if err != nil {
			that.Error(constants.SYSTEM_ERROR, "Failed to create file", "")
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			that.Error(constants.SYSTEM_ERROR, "Failed to copy file", "")
			return
		}
	}

	js := jsonutil.NewJSONObject()
	js.FluentPut("baseURL", constants.MAPPING_URL+"/"+fileName)
	js.FluentPut("url", constants.IMG_BASE_URL+constants.MAPPING_URL+"/"+fileName)

	that.Ok("Success", js.GetData())
}

func getPath() string {
	// Create a folder if the folder does not exist
	if err := os.MkdirAll(constants.IMG_PATH+"/", os.ModePerm); err != nil {
		panic(err)
	}
	return constants.IMG_PATH + "/"
}
