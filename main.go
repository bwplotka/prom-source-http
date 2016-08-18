package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "prom-source-http",
		Short: "Prometheus source HTTP Server",
		Long:  `[prom-source-http] is a HTTP Server which serves a content of the file specified in '-f' option.
It can be useful for debug and test purposes. It can be attached as a prometheus sorce by adding scrape_config job:

[prometheus.yml]:

scrape_configs:
  # The job name is added as a label 'job=<job_name>' to any timeseries scraped from this config.
  - job_name: 'fabric'
    metrics_path: /_metrics

	static_configs:
      - targets: ['localhost:9091']

It provides two endpoints:
- GET /metrics -> which serves an output of the file.
- GET /metrics.json?url=<some prometheus endpoint> -> which serves prometheus metrics fetched from url in raw fomat
and render them in readable JSON format.`,
		Run:   runHttp,
	}
	cfgPort string
	fileToServe string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgPort, "port", "p", "9091", "Port to serve a server on.")
	rootCmd.PersistentFlags().StringVarP(&fileToServe, "file", "f", "./_metrics", "File which content will be served")
}

func runHttp(cmd *cobra.Command, args []string) {
	ServeHTTP(cfgPort)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
