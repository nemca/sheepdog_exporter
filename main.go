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
	md        = "md_info"
	ns        = "node_stat"
)

var (
	mdInfoSize = prometheus.NewDesc(
		prometheus.BuildFQName(namespase, md, "size"),
		"Multi-disk total size in bytes by path.",
		[]string{"path"}, nil,
	)
	mdInfoUse = prometheus.NewDesc(
		prometheus.BuildFQName(namespase, md, "use"),
		"Multi-disk usage in percentage by path.",
		[]string{"path"}, nil,
	)
	mdInfoAvail = prometheus.NewDesc(
		prometheus.BuildFQName(namespase, md, "avail"),
		"Multi-disk available size in bytes by path.",
		[]string{"path"}, nil,
	)
	mdInfoUsed = prometheus.NewDesc(
		prometheus.BuildFQName(namespase, md, "used"),
		"Multi-disk used size in bytes by path.",
		[]string{"path"}, nil,
	)
	nodeStatActive = prometheus.NewDesc(
		prometheus.BuildFQName(namespase, ns, "active"),
		"Number of running requests by type.",
		[]string{"type"}, nil,
	)
	nodeStatTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespase, ns, "total"),
		"Total numbers of requests received by type.",
		[]string{"type"}, nil,
	)
	nodeStatWrite = prometheus.NewDesc(
		prometheus.BuildFQName(namespase, ns, "write"),
		"Number of write requests by type.",
		[]string{"type"}, nil,
	)
	nodeStatRead = prometheus.NewDesc(
		prometheus.BuildFQName(namespase, ns, "read"),
		"Number of read requests by type.",
		[]string{"type"}, nil,
	)
	nodeStatRemove = prometheus.NewDesc(
		prometheus.BuildFQName(namespase, ns, "remove"),
		"Number of remove requests by type.",
		[]string{"type"}, nil,
	)
	nodeStatFlush = prometheus.NewDesc(
		prometheus.BuildFQName(namespase, ns, "flush"),
		"Number of flush requests by type.",
		[]string{"type"}, nil,
	)
	nodeStatAllWrite = prometheus.NewDesc(
		prometheus.BuildFQName(namespase, ns, "write_all"),
		"Number of all write requests by type.",
		[]string{"type"}, nil,
	)
	nodeStatAllRead = prometheus.NewDesc(
		prometheus.BuildFQName(namespase, ns, "read_all"),
		"Number of all read requests by type.",
		[]string{"type"}, nil,
	)
)

type dogCollector struct {
	mdInfoUse        *prometheus.Desc
	mdInfoSize       *prometheus.Desc
	mdInfoAvail      *prometheus.Desc
	mdInfoUsed       *prometheus.Desc
	nodeStatActive   *prometheus.Desc
	nodeStatTotal    *prometheus.Desc
	nodeStatWrite    *prometheus.Desc
	nodeStatRead     *prometheus.Desc
	nodeStatRemove   *prometheus.Desc
	nodeStatFlush    *prometheus.Desc
	nodeStatAllWrite *prometheus.Desc
	nodeStatAllRead  *prometheus.Desc
}

func newDogCollector() *dogCollector {
	return &dogCollector{
		mdInfoSize:       mdInfoSize,
		mdInfoUse:        mdInfoUse,
		mdInfoAvail:      mdInfoAvail,
		mdInfoUsed:       mdInfoUsed,
		nodeStatActive:   nodeStatActive,
		nodeStatTotal:    nodeStatTotal,
		nodeStatWrite:    nodeStatWrite,
		nodeStatRead:     nodeStatRead,
		nodeStatRemove:   nodeStatRemove,
		nodeStatFlush:    nodeStatFlush,
		nodeStatAllWrite: nodeStatAllWrite,
		nodeStatAllRead:  nodeStatAllRead,
	}
}

func (c *dogCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.mdInfoUse
	ch <- c.mdInfoSize
	ch <- c.mdInfoAvail
	ch <- c.mdInfoUsed
	ch <- c.nodeStatActive
	ch <- c.nodeStatTotal
	ch <- c.nodeStatWrite
	ch <- c.nodeStatRead
	ch <- c.nodeStatRemove
	ch <- c.nodeStatFlush
	ch <- c.nodeStatAllWrite
	ch <- c.nodeStatAllRead
}

func (c *dogCollector) Collect(ch chan<- prometheus.Metric) {
	mds, err := getMdInfo()
	if err != nil {
		log.Error(err)
	}
	for _, md := range mds {
		ch <- prometheus.MustNewConstMetric(c.mdInfoSize, prometheus.GaugeValue, float64(md.Size), md.Path)
		ch <- prometheus.MustNewConstMetric(c.mdInfoUse, prometheus.GaugeValue, float64(md.Use), md.Path)
		ch <- prometheus.MustNewConstMetric(c.mdInfoAvail, prometheus.GaugeValue, float64(md.Avail), md.Path)
		ch <- prometheus.MustNewConstMetric(c.mdInfoUsed, prometheus.GaugeValue, float64(md.Used), md.Path)
	}

	ns, err := getNodeStat()
	if err != nil {
		log.Error(err)
	}
	for _, s := range ns {
		ch <- prometheus.MustNewConstMetric(c.nodeStatActive, prometheus.GaugeValue, float64(s.Active), s.Type)
		ch <- prometheus.MustNewConstMetric(c.nodeStatTotal, prometheus.GaugeValue, float64(s.Total), s.Type)
		ch <- prometheus.MustNewConstMetric(c.nodeStatWrite, prometheus.GaugeValue, float64(s.Write), s.Type)
		ch <- prometheus.MustNewConstMetric(c.nodeStatRead, prometheus.GaugeValue, float64(s.Read), s.Type)
		ch <- prometheus.MustNewConstMetric(c.nodeStatRemove, prometheus.GaugeValue, float64(s.Remove), s.Type)
		ch <- prometheus.MustNewConstMetric(c.nodeStatFlush, prometheus.GaugeValue, float64(s.Flush), s.Type)
		ch <- prometheus.MustNewConstMetric(c.nodeStatAllWrite, prometheus.GaugeValue, float64(s.AllWrite), s.Type)
		ch <- prometheus.MustNewConstMetric(c.nodeStatAllRead, prometheus.GaugeValue, float64(s.AllRead), s.Type)
	}
}

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

	dc := newDogCollector()
	prometheus.MustRegister(dc)

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
