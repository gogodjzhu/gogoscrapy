package downloader

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gogoproxypool/src"
	"gogoscrapy/src/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MysqlProxyFactory struct {
	proxyMapper *src.ProxyMapper
}

func NewMysqlProxyFactory(config *src.MySqlConfig) (*MysqlProxyFactory, error) {
	persist, err := src.NewPersistence(config)
	if err != nil {
		return nil, err
	}
	return &MysqlProxyFactory{
		proxyMapper: src.NewProxyMapper(persist),
	}, nil
}

func (this *MysqlProxyFactory) GetProxy() (IProxy, error) {
	proxy, err := this.proxyMapper.Get()
	if err != nil {
		return nil, err
	}
	return NewProxy(proxy.Id, proxy.Host, proxy.Port, proxy.Username, proxy.Password), nil
}

func (this *MysqlProxyFactory) ReturnProxy(proxy IProxy) {
	if err := this.proxyMapper.ReturnCache(proxy.GetId()); err != nil {
		LOG.Errorf("failed to return proxy, err:%+v", err)
	}
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIN    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type ApiProxyFactory struct {
	oauth2Host      string
	proxyServerHost string

	oauth2ClientId     string
	oauth2ClientSecret string

	token             *AccessToken
	bucketServiceCode string
}

func NewApiProxyFactory(oauth2Client, oauth2ClientSecret, oauth2Host, proxyServerHost, bucketName string) (ApiProxyFactory, error) {
	token, err := getAccessToken(oauth2Host, proxyServerHost)
	if err != nil {
		return ApiProxyFactory{}, err
	}
	serviceCode, err := getServiceCode(proxyServerHost, token.AccessToken, bucketName)
	if err != nil {
		return ApiProxyFactory{}, err
	}
	fac := ApiProxyFactory{
		token:              token,
		bucketServiceCode:  serviceCode,
		oauth2Host:         oauth2Host,
		proxyServerHost:    proxyServerHost,
		oauth2ClientId:     oauth2Client,
		oauth2ClientSecret: oauth2ClientSecret,
	}
	go func() {
		for {
			time.Sleep(time.Duration(fac.token.ExpiresIN*9/10) * time.Second)
			err := fac.refreshToken()
			if err != nil {
				LOG.Errorf("Failed to refresh token err:%+v", err)
			}
			LOG.Infof("Update token to %+v", *fac.token)
		}
	}()
	return fac, nil
}

func (this *ApiProxyFactory) refreshToken() error {
	var token AccessToken
	err := utils.Post2Json(this.oauth2Host+"/v1/oauth2/refreshToken",
		url.Values{
			"grant_type":    {"refresh_token"},
			"client_id":     {this.oauth2ClientId},
			"client_secret": {this.oauth2ClientSecret},
			"domain":        {this.proxyServerHost},
			"refresh_token": {this.token.RefreshToken},
		}, &token)
	if err != nil {
		return err
	}
	this.token = &token
	return nil
}

func getOauth2LoginCookie(oauth2Host string) (*http.Cookie, error) {
	resp, err := http.Get(oauth2Host + "/v1/user/login?username=gogodjzhu&password=pass1123")
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if string(body) != "\"login success\"" {
		return nil, errors.New("Failed to login to oauth2, error:" + string(body))
	}
	for headKey, headValues := range resp.Header {
		if headKey == "Set-Cookie" {
			for _, headValue := range headValues {
				reg := regexp.MustCompile(`oauth2-session=.*?;`)
				sessionPair := reg.FindString(headValue)
				sessionPair = sessionPair[:len(sessionPair)-1] //去除';'
				if sessionPair != "" {
					cookie := http.Cookie{
						Name:  strings.Split(sessionPair, "=")[0],
						Value: strings.Split(sessionPair, "=")[1],
					}
					return &cookie, nil
				}
			}
		}
	}
	return nil, errors.New("Failed to get login cookie")
}

func getAuthorizeUrl(serverHost string) (string, error) {
	type grantResp struct {
		Msg  string
		Code int
		Data string
	}
	resp := grantResp{}
	if err := utils.Get2Json(serverHost+"/v1/auth/grant", &resp); err != nil {
		return "", err
	}
	if resp.Code != 0 {
		return "", errors.New(fmt.Sprintf("Failed to get authorizeUrl %+v", resp))
	}
	return resp.Data, nil
}

func getAccessToken(oauth2Host, serverHost string) (*AccessToken, error) {
	cookie, err := getOauth2LoginCookie(oauth2Host)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	defer client.CloseIdleConnections()
	authorizeUrl, err := getAuthorizeUrl(serverHost)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodGet, authorizeUrl, nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var accessTokenResp AccessToken
	if err = json.Unmarshal(body, &accessTokenResp); err != nil {
		return nil, err
	}
	return &accessTokenResp, nil
}

func getServiceCode(serverHost, accessToken, bucketName string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, serverHost+"/v1/proxy/bucket/get?bucketId="+bucketName, nil)
	if err != nil {
		return "", err
	}
	req.AddCookie(&http.Cookie{Name: "access_token", Value: accessToken})
	client := http.Client{}
	defer client.CloseIdleConnections()
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	type Bucket struct {
		ID                   int64     `gorm:"id"`
		Nick                 string    `gorm:"nick"`
		ValidateUrl          string    `gorm:"validate_url"`
		ValidateHtmlPattern  string    `gorm:"validate_html_pattern"`
		ValidateHtmlExpected string    `gorm:"validate_html_expected"`
		MaxExpiredTime       int64     `gorm:"max_expired_time"`
		CreateTime           time.Time `gorm:"create_time"`
		Locked               bool      `gorm:"locked"`
	}
	var retMap map[string]Bucket
	if err = json.Unmarshal(body, &retMap); err != nil {
		return "", err
	}
	for serviceCode := range retMap {
		return serviceCode, nil
	}
	return "", errors.New("Non bucket available, name:" + bucketName)
}

func (this *ApiProxyFactory) GetProxy() (IProxy, error) {
	req, err := http.NewRequest(http.MethodGet, this.proxyServerHost+"/v1/proxy/bucket/proxy?bucketServiceCode="+this.bucketServiceCode, nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "access_token", Value: this.token.AccessToken})
	client := http.Client{}
	defer client.CloseIdleConnections()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	type Resp struct {
		Msg  string
		Code int
		Data Proxy
	}
	var proxyResp Resp
	if err = json.Unmarshal(body, &proxyResp); err != nil {
		return nil, err
	}
	if proxyResp.Code != 0 {
		return nil, errors.New(fmt.Sprintf("Invalid proxy return: %+v", proxyResp))
	}
	return proxyResp.Data, nil
}

func (this *ApiProxyFactory) ReturnProxy(proxy IProxy) {
	req, err := http.NewRequest(http.MethodPost, this.proxyServerHost+"/v1/proxy/bucket/proxy", nil)
	req.PostForm = map[string][]string{
		"bucketServiceCode": {this.bucketServiceCode},
		"proxyId":           {strconv.Itoa(proxy.GetId())},
	}
	if err != nil {
		LOG.Error(err)
		return
	}
	req.AddCookie(&http.Cookie{Name: "access_token", Value: this.token.AccessToken})
	client := http.Client{}
	defer client.CloseIdleConnections()
	resp, err := client.Do(req)
	if err != nil {
		LOG.Error(err)
		return
	}
	defer resp.Body.Close()
}
