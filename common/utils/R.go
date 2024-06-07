package utils

type R struct {
	Code int32
	Msg  string
	Data interface{}
}

func Success(data interface{}) R {
	var r R
	r.Code = 200
	r.Msg = "success"
	r.Data = data
	return r
}

func Fail(code int32, msg string) R {
	var r R
	r.Code = code
	r.Msg = msg
	r.Data = ""
	return r
}
