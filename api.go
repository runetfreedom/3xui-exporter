package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"msg"`
}

type ClientStat struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
	Up    uint64 `json:"up"`
	Down  uint64 `json:"down"`
	Total uint64 `json:"total"`
}

type Inbounds struct {
	Id    int64  `json:"id"`
	Up    uint64 `json:"up"`
	Down  uint64 `json:"down"`
	Total uint64 `json:"total"`

	ClientStats []ClientStat `json:"clientStats"`
}

type InboundsListResponse struct {
	Response
	Inbounds []Inbounds `json:"obj"`
}

func login() error {
	resp, err := httpClient.PostForm(*panelUrl+"/login", url.Values{
		"username": {*panelUserName},
		"password": {*panelPassWord},
	})

	if err != nil {
		return err
	}

	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	var loginResp Response
	if err = json.Unmarshal(data, &loginResp); err != nil {
		return err
	}

	if !loginResp.Success {
		return fmt.Errorf("login failed: %s", loginResp.Message)
	}

	return nil
}

func getInbounds() ([]Inbounds, error) {
	resp, err := httpClient.Get(*panelUrl + "/panel/api/inbounds/list")

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var obj InboundsListResponse
	if err = json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	if !obj.Success {
		return nil, fmt.Errorf("get inbounds list failed: %s", obj.Message)
	}

	return obj.Inbounds, nil
}
