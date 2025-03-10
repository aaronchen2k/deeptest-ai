package _http

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	_str "github.com/deeptest-com/deeptest-next/pkg/libs/string"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	Verbose = true
)

func Get(url string) (ret []byte, err error) {
	if Verbose {
		_logUtils.Infof("===DEBUG===  request: %s", url)
	}

	client := &http.Client{
		Timeout: 8 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		_logUtils.Infof(color.RedString("get request failed, error: %s.", err.Error()))
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		_logUtils.Infof(color.RedString("get request failed, error: %s.", err.Error()))
		return
	}
	defer resp.Body.Close()

	if !IsSuccessCode(resp.StatusCode) {
		_logUtils.Infof(color.RedString("read response failed, StatusCode: %d.", resp.StatusCode))
		err = errors.New(resp.Status)
		return
	}

	reader := resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, _ = gzip.NewReader(resp.Body)
	}

	unicodeContent, _ := ioutil.ReadAll(reader)
	ret, _ = _str.UnescapeUnicode(unicodeContent)

	return
}

func Post(url string, data interface{}) (ret []byte, err error) {
	return PostOrPut(url, "POST", data)
}
func Put(url string, data interface{}) (ret []byte, err error) {
	return PostOrPut(url, "PUT", data)
}
func Delete(url string) (ret []byte, err error) {
	if Verbose {
		_logUtils.Infof("===DEBUG===  request: %s", url)
	}

	client := &http.Client{
		Timeout: 8 * time.Second,
	}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		_logUtils.Infof(color.RedString("get request failed, error: %s.", err.Error()))
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		_logUtils.Infof(color.RedString("get request failed, error: %s.", err.Error()))
		return
	}
	defer resp.Body.Close()

	if !IsSuccessCode(resp.StatusCode) {
		_logUtils.Infof(color.RedString("read response failed, StatusCode: %d.", resp.StatusCode))
		err = errors.New(resp.Status)
		return
	}

	reader := resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, _ = gzip.NewReader(resp.Body)
	}

	unicodeContent, _ := ioutil.ReadAll(reader)
	ret, _ = _str.UnescapeUnicode(unicodeContent)

	return
}

func PostOrPut(url string, method string, data interface{}) (ret []byte, err error) {
	if Verbose {
		_logUtils.Infof("===DEBUG===  request: %s", url)
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	dataBytes, err := json.Marshal(data)
	if Verbose {
		_logUtils.Infof("===DEBUG===     data: %s", string(dataBytes))
	}

	if err != nil {
		_logUtils.Infof(color.RedString("marshal request failed, error: %s.", err.Error()))
		return
	}

	dataStr := string(dataBytes)

	req, err := http.NewRequest(method, url, strings.NewReader(dataStr))
	if err != nil {
		_logUtils.Infof(color.RedString("post request failed, error: %s.", err.Error()))
		return
	}

	//req.Header.SetVariable("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		_logUtils.Infof(color.RedString("post request failed, error: %s.", err.Error()))
		return
	}

	if !IsSuccessCode(resp.StatusCode) {
		_logUtils.Infof(color.RedString("post request return '%s'.", resp.Status))
		err = errors.New(resp.Status)
		return
	}

	reader := resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, _ = gzip.NewReader(resp.Body)
	}

	unicodeContent, _ := ioutil.ReadAll(reader)
	ret, _ = _str.UnescapeUnicode(unicodeContent)

	return
}

func IsSuccessCode(code int) (success bool) {
	return code >= 200 && code <= 299
}

func GenUrlWithParams(pth string, params map[string]interface{}, baseUrl string) (url string) {
	uri := pth

	index := 0
	for key, val := range params {
		if index == 0 {
			uri += "?"
		} else {
			uri += "&"
		}

		uri += fmt.Sprintf("%v=%v", key, val)
		index++
	}

	url = baseUrl + uri

	return
}
