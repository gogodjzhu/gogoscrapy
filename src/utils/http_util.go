package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func Get2String(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func Get2Json(url string, res interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("status code is " + strconv.Itoa(resp.StatusCode))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, res); err != nil {
		return err
	}
	return nil
}

func Post2Json(url string, data url.Values, res interface{}) error {
	resp, err := http.PostForm(url, data)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("status code is " + strconv.Itoa(resp.StatusCode))
	}
	defer resp.Body.Close() //TODO 检查错误不关闭
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, res); err != nil {
		return err
	}
	return nil
}
