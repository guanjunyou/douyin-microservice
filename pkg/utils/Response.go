package utils

type Response struct {
	CommonEntity
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type MsgResponse struct {
	CommonEntity
	StatusCode string `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
