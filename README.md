# Sheepdog Exporter

Export Sheepdog service health to Prometheus.

#### Build

```bash
make
./sheepdog_exporter [flags]
```

#### Usage

```bash
./sheepdog_exporter -h
usage: sheepdog_exporter [<flags>]

Flags:
  -h, --help                  Show context-sensitive help (also try --help-long and --help-man).
      --web.listen-address=":9525"
                              Address to listen on for web interface and telemetry.
      --web.telemetry-path="/metrics"
                              Path under which to expose metrics.
      --sheepdog.pid-file=""  Path to Sheepdog's pid file to export process information.
      --log.level="info"      Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]
      --log.format="logger:stderr"
                              Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true"
      --version               Show application version.
```

## Exported Metrics
| Metric | Description | Type | Labels |
| ------ | ----------- | ---- | ------ |
| sheepdog_md_info_avail | Multi-disk available size in bytes | gauge | path |
| sheepdog_md_info_size | Multi-disk total size in bytes | gauge | path |
| sheepdog_md_info_use | Multi-disk usage in percentage | gauge | path |
| sheepdog_md_info_used | Multi-disk used size in bytes | gauge | path |
| sheepdog_node_stat_active | Number of running requests | gauge | type |
| sheepdog_node_stat_flush | Number of flush requests | gauge | type |
| sheepdog_node_stat_read | Number of read requests | gauge | type |
| sheepdog_node_stat_read_all | Number of all read requests | gauge | type |
| sheepdog_node_stat_remove | Number of remove requests | gauge | type |
| sheepdog_node_stat_total | Total numbers of requests received | gauge | type |
| sheepdog_node_stat_write | Number of write requests | gauge | type |
| sheepdog_node_stat_write_all | Number of all write requests | gauge | type |
| sheepdog_process_cpu_seconds_total | Total user and system CPU time spent in seconds | counter | |
| sheepdog_process_max_fds | Maximum number of open file descriptors | gauge | |
| sheepdog_process_open_fds | Number of open file descriptors | gauge | |
| sheepdog_process_resident_memory_bytes | Resident memory size in bytes | gauge | |
| sheepdog_process_virtual_memory_bytes | Virtual memory size in bytes | gauge | |
| sheepdog_process_virtual_memory_max_bytes | Maximum amount of virtual memory available in bytes | gauge | |
| sheepdog_process_start_time_seconds | Start time of the process since unix epoch in seconds | gauge | |
