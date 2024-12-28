package httpUtils

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"errors"
	_http "github.com/deeptest-com/deeptest-next/pkg/libs/http"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	_str "github.com/deeptest-com/deeptest-next/pkg/libs/string"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	Verbose = true
)

func Get(url string, headers map[string]string) (ret []byte, err error) {
	return gets(url, "GET", headers)
}

func Post(url string, data interface{}, headers map[string]string) (ret []byte, err error) {
	return posts(url, "POST", data, headers)
}
func Put(url string, data interface{}, headers map[string]string) (ret []byte, err error) {
	return posts(url, "PUT", data, headers)
}
func Delete(url string, headers map[string]string) (ret []byte, err error) {
	return gets(url, "DELETE", headers)
}
func PostFile(url string, data interface{}, filePath string, headers map[string]string) (ret []byte, err error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 8 * time.Second,
	}

	var dataBytes []byte

	dataBytes, err = json.Marshal(data)
	formData := []BodyFormDataItem{}
	formItemText := BodyFormDataItem{
		Name:  "data",
		Value: string(dataBytes),
		Type:  FormDataTypeText,
	}
	formData = append(formData, formItemText)

	formItemFile := BodyFormDataItem{
		Name:  "file",
		Value: filePath,
		Type:  FormDataTypeFile,
	}
	formData = append(formData, formItemFile)

	bodyFormData := GenBodyFormDataFromItems(formData)

	formDataWriter, _ := MultipartEncoder(bodyFormData)
	formDataContentType := MultipartContentType(formDataWriter)

	dataBytes = formDataWriter.Payload.Bytes()

	if err != nil {
		_logUtils.Infof(color.RedString("marshal httpReq failed, error: %s.", err.Error()))
		return
	}

	httpReq, reqErr := http.NewRequest("POST", url, bytes.NewReader(dataBytes))
	if reqErr != nil {
		_logUtils.Error(reqErr.Error())
		return
	}

	httpReq.Header.Set("Content-Type", formDataContentType)
	for key, value := range headers {
		httpReq.Header.Set(key, value)
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		_logUtils.Infof(color.RedString("post request failed, error: %s.", err.Error()))
		return
	}

	defer resp.Body.Close()

	reader := resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, _ = gzip.NewReader(resp.Body)
	}

	unicodeContent, _ := io.ReadAll(reader)
	ret, _ = _str.UnescapeUnicode(unicodeContent)

	// check response status
	if !_http.IsSuccessCode(resp.StatusCode) {
		_logUtils.Infof(color.RedString("post request return '%s'.", resp.Status))
		err = errors.New(resp.Status)
	}

	return
}

func gets(url, method string, headers map[string]string) (ret []byte, err error) {
	if Verbose {
		_logUtils.Infof("===DEBUG===  request: %s", url)
	}

	client := &http.Client{
		Timeout: 8 * time.Second,
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		_logUtils.Infof(color.RedString("get request failed, error: %s.", err.Error()))
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		_logUtils.Infof(color.RedString("get request failed, error: %s.", err.Error()))
		return
	}
	defer resp.Body.Close()

	if !_http.IsSuccessCode(resp.StatusCode) {
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
func posts(url string, method string, data interface{}, headers map[string]string) (ret []byte, err error) {
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
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		_logUtils.Infof(color.RedString("post request failed, error: %s.", err.Error()))
		return
	}

	if !_http.IsSuccessCode(resp.StatusCode) {
		_logUtils.Infof(color.RedString("post request return '%s'.", resp.Status))
		err = errors.New(resp.Status)
		return
	}

	reader := resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, _ = gzip.NewReader(resp.Body)
	}

	unicodeContent, _ := io.ReadAll(reader)
	ret, _ = _str.UnescapeUnicode(unicodeContent)

	return
}

func GenBodyFormDataFromItems(items []BodyFormDataItem) (formData []BodyFormDataItem) {
	mp := map[string]bool{}

	if items != nil {
		for _, item := range items {
			key := item.Name
			if _, ok := mp[key]; ok { // skip duplicate one
				continue
			}

			formData = append(formData, item)
			mp[key] = true
		}
	}

	return
}
