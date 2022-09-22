package wechat

import "github.com/medivhzhan/weapp/v2"

type Service struct {
	AppID     string
	AppSecret string
}

func (s *Service) Resolve(code string) (string, error) {
	resp, err := weapp.Login(s.AppID, s.AppSecret, code)
	if err != nil {
		return "", err
	}
	if err := resp.GetResponseError(); err != nil {
		return "", err
	}
	return resp.OpenID, nil
}
