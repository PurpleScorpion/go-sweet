package service

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"service-common/src/main/models"
	"shared/constants"
	"shared/utils"
	"shared/vo"
	"strings"

	"github.com/PurpleScorpion/go-sweet-keqing/keqing"
	"github.com/PurpleScorpion/go-sweet-orm/v3/mapper"
)

type DemoService struct {
}

func (s DemoService) Demo1(info vo.UserVO) utils.R {
	return utils.Success("http://localhost:28666/static/1769825432_krGjH.webp")
}

func (s DemoService) UploadFile(fileHeader *multipart.FileHeader) utils.R {

	if fileHeader == nil {
		return utils.Fail(constants.FILE_OPEN_ERROR, "无法打开文件。请确认文件完整性后，再试一次。")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return utils.Fail(constants.FILE_OPEN_ERROR, "无法打开文件。请确认文件完整性后，再试一次。")
	}
	defer src.Close()

	originalName := fileHeader.Filename
	msg := checkFileType(originalName)
	if msg != "ok" {
		return utils.Fail(constants.FILE_OPEN_ERROR, msg)
	}
	fileName := utils.RandomFileName(originalName)
	fileSize := fileHeader.Size

	// 计算文件 MD5 值
	fileMD5 := getFileMD5(src)
	if keqing.IsEmpty(fileMD5) {
		return utils.Fail(constants.FILE_OPEN_ERROR, "无法打开文件。请确认文件完整性后，再试一次。")
	}
	certPath := utils.GetCertPath()

	objectName := utils.SaveFile(certPath, fileName, src)
	if keqing.IsEmpty(objectName) {
		return utils.Fail(constants.FILE_WRITE_ERROR, "文件保存失败。请稍后再试。")
	}

	// 确立自身
	var obj models.Cert
	obj.OriginalName = originalName
	obj.FileName = fileName
	obj.FileSize = int32(fileSize)
	obj.FileMd5 = fileMD5
	obj.CreatedTime = keqing.NowDate()
	mapper.Insert[models.Cert](&obj, nil)

	return utils.Success("")
}

func (s DemoService) DownloadFile() (utils.R, string) {

	list := mapper.SelectList[models.Cert](mapper.BuilderQueryWrapper().
		Eq(true, "delete_flag", constants.NO_DELETE_CODE),
	)

	if keqing.IsEmpty(list) {
		return utils.Fail(constants.FILE_NOT_EXIST, "未找到可下载的文件。"), ""
	}

	fileName := list[0].FileName

	certPath := utils.GetCertPath() + "/" + fileName

	return utils.Success(certPath), list[0].OriginalName
}

func (s DemoService) Demo2(code string) utils.R {
	return utils.Success(code)
}

func (s DemoService) Demo3(pageVO vo.UserVO) utils.R {
	return utils.Success(pageVO)
}

func checkFileType(fileName string) string {
	// 获取文件扩展名并转换为小写
	ext := strings.ToLower(filepath.Ext(fileName))

	// 检查文件扩展名是否为 zip/rar/7z 格式
	validExtensions := []string{".crt", ".pem"}
	isValidExt := false
	for _, validExt := range validExtensions {
		if ext == validExt {
			isValidExt = true
			break
		}
	}

	if !isValidExt {
		return "文件格式不支持。请上传cert或pem格式的文件。"
	}
	return "ok"
}

func getFileMD5(src multipart.File) string {
	hash := md5.New()
	_, err := io.Copy(hash, src)
	if err != nil {
		return ""
	}
	// 重要：重置文件指针到开头
	_, err = src.Seek(0, io.SeekStart)
	if err != nil {
		return ""
	}
	md5Sum := fmt.Sprintf("%x", hash.Sum(nil))
	return md5Sum
}
