package service

import (
	"bytes"
	"encoding/json"
	"go-fcm-example/admin/src/constant"
	"go-fcm-example/admin/src/define"
	"net/http"
)

/**
* @Author: caishi13202
* @Date: 2021/9/27 3:52 下午
 */
type Notification struct {
	loginUser  map[string]string
	httpClient *http.Client
}

func (n Notification) ListAccount() map[string]string {
	return n.loginUser
}

func (n Notification) Login(req *define.AccessTokenReq) bool {
	n.loginUser[req.AccountId] = req.Token
	return true
}

func (n Notification) Send(req *define.SendReq) bool {
	to := n.loginUser[req.AccountId]
	if to == "" {
		return false
	}
	msg := &define.Msg{
		Data: req,
		To:   to,
	}
	json, _ := json.Marshal(msg)
	request, _ := http.NewRequest(http.MethodPost, constant.FCM_SEND_URL, bytes.NewReader(json))
	request.Header.Add("Authorization", constant.FCM_AUTHORIZATION)
	request.Header.Add("Content-Type", "application/json")
	resp, _ := n.httpClient.Do(request)
	if resp.StatusCode == http.StatusOK {
		return true
	} else {
		return false
	}
}

func (n *Notification) SetLoginUser(loginUser map[string]string) {
	n.loginUser = loginUser
}

func (n *Notification) SetHttpClient(httpClient *http.Client) {
	n.httpClient = httpClient
}
