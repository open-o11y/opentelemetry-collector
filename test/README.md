#Testing Pipeline for Prometheus Remote Write Exporter
The this package contains utilities for testing the Prometheus remote write exporter. 

- `otlploadgenerator` generates and
sends metric to OTLP receiver of the Collector. 
- `querier` validates the correctness of the metric by querying a backend.
- `otel-collector-config.yaml` specifies the configuration of the OpenTelemetry Collector.

To start a Collector instance and send to it using `otlploadgenerator`, run the following command:

```
make testaps
```

