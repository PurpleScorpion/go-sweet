package vo

type DefaultPageVO struct {
	Current  int32 `json:"current"`
	PageSize int32 `json:"pageSize"`
	UserId   int32 `json:"userId"`
}
