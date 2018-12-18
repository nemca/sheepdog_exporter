package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

const (
	namespase = "sheepdog"
	exporter  = "sheepdog_exporter"
)

func main() {
	var (
		showVersion     = flag.Bool("version", false, "Print version information.")
		listenAddress   = flag.String("web.listen-address", ":9525", "Address to listen on for web interface and telemetry.")
		metricsPath     = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		sheepdogPidFile = flag.String("sheepdog.pid-file", "", "Path to Sheepdog's pid file to export process information.")
	)
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print(exporter))
		os.Exit(0)
	}

	log.Infoln("Starting", exporter, version.Info())
	log.Infoln("Build context", version.BuildContext())

	if *sheepdogPidFile != "" {
		procExporter := prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{
			func() (int, error) {
				content, err := ioutil.ReadFile(*sheepdogPidFile)
				if err != nil {
					return 0, fmt.Errorf("Can't read pid file: %s", err)
				}
				value, err := strconv.Atoi(strings.TrimSpace(string(content)))
				if err != nil {
					return 0, fmt.Errorf("Can't parse pid file: %s", err)
				}
				return value, nil
			}, namespase, false})
		prometheus.MustRegister(procExporter)
	}

	log.Info("Starting Server: ", *listenAddress)
	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Sheepdog Exporter</title></head>
			<body>
			<h1>Sheepdog Exporter</h1>
			<p><a href='` + *metricsPath + `'>Metrics</a></p>
			</body>
			</html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
