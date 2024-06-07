package controllers

import (
	"github.com/PurpleScorpion/go-sweet-json/jsonutil"
	"go-sweet/common/constants"
	"go-sweet/common/utils"
	sweetyml "go-sweet/common/yaml"
	"io"
	"os"
)

type FileController struct {
	BaseController
}

func (that *FileController) GetBaseImg() {
	that.Ok("Success", "")
}

func (that *FileController) UploadImg() {
	conf := sweetyml.GetYmlConf()

	files, _ := that.GetFiles("files")

	if len(files) == 0 {
		that.Error(constants.SYSTEM_ERROR, "Please select the file to upload", "")
	}

	fileHeader := files[0]

	fileName := ""

	if fileHeader != nil {

		fileName = utils.RandomImageName(fileHeader.Filename)

		savePath := getPath(conf) + fileName

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
	js.FluentPut("baseURL", conf.Sweet.Img.MappingUrl+"/"+fileName)
	js.FluentPut("url", conf.Sweet.Img.BaseUrl+conf.Sweet.Img.MappingUrl+"/"+fileName)

	that.Ok("Success", js.GetData())
}

func getPath(conf sweetyml.YmlConfig) string {
	// Create a folder if the folder does not exist
	if err := os.MkdirAll(conf.Sweet.Img.Path+"/", os.ModePerm); err != nil {
		panic(err)
	}
	return conf.Sweet.Img.Path + "/"
}
