package gapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
)

const snapshotAPI = "/api/snapshots"

type DashboardWithMeta struct {
	Meta      DashboardMeta `json:"meta"`
	Dashboard interface{}   `json:"dashboard"`
}

type DashboardSnapshot struct {
	Dashboard map[string]interface{} `json:"dashboard" binding:"Required"`
	Name      string                 `json:"name"`
	Expires   int64                  `json:"expires"`
	// these are passed when storing an external snapshot ref
	External  bool   `json:"external"`
	Key       string `json:"key"`
	DeleteKey string `json:"deleteKey"`
}

type CreateSnaphostResponse struct {
	Key       string `json:"key"`
	DeleteKey string `json:"deleteKey"`
	URL       string `json:"url"`
	DeleteURL string `json:"deleteUrl"`
}

type SharedOptionSnaphost struct {
	ExternalSnapshotURL  string `json:"externalSnapshotURL"`
	ExternalSnapshotName string `json:"externalSnapshotName"`
	ExternalEnabled      string `json:"externalEnabled"`
}

func (c *Client) GetByKey(key string) (*DashboardWithMeta, error) {
	req, err := c.newRequest("GET", "/api/snapshots/"+key, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &DashboardWithMeta{}
	err = json.Unmarshal(data, &result)
	return result, err
}

func (c *Client) Create(snapshot *DashboardSnapshot) (*CreateSnaphostResponse, error) {
	data, err := json.Marshal(snapshot)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest("POST", snapshotAPI, nil, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &CreateSnaphostResponse{}
	err = json.Unmarshal(data, &result)
	return result, err

}
