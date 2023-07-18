package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const CloudBase = "https://cloud.fermyon.com"

type LogsResponse struct {
	Logs []string `json:"logs"`
}

type Channel struct {
	Id string `json:"id"`
}

type Item struct {
	Channels []Channel `json:"channels"`
	Name     string    `json:"name"`
}

type AppsResponse struct {
	Items []Item `json:"items"`
}

type client struct {
	httpclient *http.Client
	token      string
}

func (c *client) getChannelId(appName string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/apps", CloudBase), nil)
	if err != nil {
		return "", err
	}

	raw, err := c.do(req)
	if err != nil {
		return "", err
	}

	appsResp := &AppsResponse{}
	err = json.Unmarshal(raw, appsResp)
	if err != nil {
		return "", err
	}

	for _, item := range appsResp.Items {
		if item.Name != appName {
			continue
		}

		if len(item.Channels) == 0 {
			return "", fmt.Errorf("no channel found for app %s", item.Name)
		}

		return item.Channels[0].Id, nil
	}

	return "", fmt.Errorf("failed to find channel id for app %s", appName)
}

func (c *client) getLogs(channelId string) ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/channels/%s/logs?max=0", CloudBase, channelId), nil)
	if err != nil {
		return nil, err
	}

	raw, err := c.do(req)
	if err != nil {
		return nil, err
	}

	logsresp := &LogsResponse{}
	err = json.Unmarshal(raw, logsresp)
	if err != nil {
		return nil, err
	}

	return logsresp.Logs, nil
}

func (c *client) do(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.httpclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected: %d, got: %d", http.StatusOK, resp.StatusCode)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return raw, nil
}
