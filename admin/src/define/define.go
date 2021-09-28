package define

/**
* @Author: caishi13202
* @Date: 2021/9/27 3:53 下午
 */
type SendReq struct {
	AccountId string `json:"accountId,omitempty"`
	Title     string `json:"title,omitempty"`
	Body      string `json:"body,omitempty"`
	Url       string `json:"url,omitempty"`
}
type AccessTokenReq struct {
	AccountId string `json:"accountId,omitempty"`
	Token     string `json:"token,omitempty"`
}
type Msg struct {
	Data *SendReq `json:"data,omitempty"`
	To   string   `json:"to,omitempty"`
}
