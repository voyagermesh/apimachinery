---
title: Voyager Run
menu:
  product_voyager_5.0.0-rc.11:
    identifier: voyager-run
    name: Voyager Run
    parent: reference
product_name: voyager
menu_name: product_voyager_5.0.0-rc.11
section_menu_id: reference
---
## voyager run

Run operator

### Synopsis

Run operator

```
voyager run [flags]
```

### Options

```
      --address string                        Address to listen on for web interface and telemetry. (default ":56790")
      --burst int                             The maximum burst for throttle (default 1000000)
      --cloud-config string                   The path to the cloud provider configuration file.  Empty string for no configuration file.
  -c, --cloud-provider string                 Name of cloud provider
      --custom-templates string               Glob pattern of custom HAProxy template files used to override built-in templates
      --exporter-sidecar-image string         Docker image containing Prometheus exporter (default "appscode/voyager:5.0.0-rc.11")
      --haproxy-image string                  Docker image containing HAProxy binary (default "appscode/haproxy:1.7.9-5.0.0-rc.11")
      --haproxy.server-metric-fields string   Comma-separated list of exported server metrics. See http://cbonte.github.io/haproxy-dconv/configuration-1.5.html#9.1 (default "2,3,4,5,6,7,8,9,13,14,15,16,17,18,21,24,33,35,38,39,40,41,42,43,44")
      --haproxy.timeout duration              Timeout for trying to get stats from HAProxy. (default 5s)
  -h, --help                                  help for run
      --ingress-class string                  Ingress class handled by voyager. Unset by default. Set to voyager to only handle ingress with annotation kubernetes.io/ingress.class=voyager.
      --kubeconfig string                     Path to kubeconfig file with authorization information (the master location is set by the master flag).
      --master string                         The address of the Kubernetes API server (overrides any value in kubeconfig)
      --operator-service string               Name of service used to expose voyager operator (default "voyager-operator")
      --prometheus-crd-apigroup string        prometheus CRD  API group name (default "monitoring.coreos.com")
      --prometheus-crd-kinds CrdKinds          - EXPERIMENTAL (could be removed in future releases) - customize CRD kind names
      --qps float32                           The maximum QPS to the master from this client (default 1e+06)
      --rbac                                  Enable RBAC for operator & offshoot Kubernetes objects
      --restrict-to-operator-namespace        If true, voyager operator will only handle Kubernetes objects in its own namespace.
      --resync-period duration                If non-zero, will re-list this often. Otherwise, re-list will be delayed aslong as possible (until the upstream source closes the watch or times out. (default 5m0s)
```

### Options inherited from parent commands

```
      --alsologtostderr                  log to standard error as well as files
      --analytics                        Send analytical events to Google Analytics (default true)
      --log.format logFormatFlag         Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true" (default "logger:stderr")
      --log.level levelFlag              Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal] (default "info")
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --stderrthreshold severity         logs at or above this threshold go to stderr
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO

* [voyager](/docs/reference/voyager.md)	 - Voyager by Appscode - Secure Ingress Controller for Kubernetes

