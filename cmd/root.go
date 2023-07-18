package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type CLI struct {
	appName string
	envName string
}

var cli = &CLI{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "spin tail",
	Short: "tail the Fermyon Cloud logs",
	Run: func(cmd *cobra.Command, args []string) {
		err := cli.run()
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.Flags().StringVar(&cli.appName, "app-name", "", "app name for which you want to run tail")
	rootCmd.MarkFlagRequired("app-name")

	rootCmd.Flags().StringVar(&cli.envName, "env-name", "", "env name to use for reading token")
}

func (cli *CLI) run() error {
	envName := "config"
	if cli.envName != "" {
		envName = cli.envName
	}

	token, err := getToken(envName)
	if err != nil {
		return fmt.Errorf("failed to find token %v", err)
	}

	fc := &client{
		token:      token,
		httpclient: &http.Client{Timeout: 15 * time.Second},
	}

	channelId, err := fc.getChannelId(cli.appName)
	if err != nil {
		return err
	}

	logsSoFar, err := fc.getLogs(channelId)
	if err != nil {
		return err
	}
	printLogs(logsSoFar)

	ticker := time.NewTicker(5 * time.Second)
	for {
		<-ticker.C

		newLogs, err := fc.returnNewLogs(channelId, len(logsSoFar))
		if err != nil {
			return err
		}

		printLogs(newLogs)
		logsSoFar = append(logsSoFar, newLogs...)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
