package cmd

import "fmt"

func (c *client) returnNewLogs(channelId string, existingEntriesCount int) ([]string, error) {
	allLogs, err := c.getLogs(channelId)
	if err != nil {
		return nil, err
	}

	return allLogs[existingEntriesCount:], nil
}

func printLogs(logs []string) {
	for _, l := range logs {
		fmt.Println(l)
	}
}
