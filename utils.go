package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type RequestBuilder = func(category, interval, metrics string) string

func NewRequestBuilder(resourceUrl string) RequestBuilder {
	return func(category, interval, metrics string) string {
		return fmt.Sprintf("%s/stats?category=%s&interval=%s&metrics=%s", resourceUrl, category, interval, metrics)
	}
}

func NewRequestBuilderWithoutInterval(resourceUrl string) RequestBuilder {
	return func(category, _, metrics string) string {
		return fmt.Sprintf("%s/stats?category=%s&metrics=%s", resourceUrl, category, metrics)
	}
}

type respTab struct {
	Headers []string          `json:"headers"`
	Value   []json.RawMessage `json:"value"`
}

func GetRespVal(rawResp []byte) ([]json.RawMessage, error) {
	var tabObj respTab
	err := json.Unmarshal(rawResp, &tabObj)
	retObj := tabObj.Value
	return retObj, err
}

func GetRespRaw(url string) ([]byte, error) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status + string(body))
	}
	return body, nil
}

func GetVal(url string) ([]map[string]string, error) {
	resp, err := GetRespRaw(url)
	if err != nil {
		return nil, err
	}
	xs, err := GetRespVal(resp)
	if err != nil {
		return nil, err
	}
	return ValToStrMap(xs)
}

func ValToStrMap(xs []json.RawMessage) ([]map[string]string, error) {
	var ret []map[string]string
	for _, x := range xs {
		var xMap map[string]string
		if err := json.Unmarshal(x, &xMap); err != nil {
			return nil, err
		} else {
			ret = append(ret, xMap)
		}
	}
	return ret, nil
}

func NewTickerSec(n int) *time.Ticker {
	return time.NewTicker(time.Duration(n) * time.Second)
}
