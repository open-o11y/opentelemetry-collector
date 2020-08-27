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

##`otlploadgenerator`
The load generator first creates a `data.txt` file. This file is need so that the querier knows what the input data is.
Each line in the file represents and OTLP metric. Then, it parse each line from the file and build OTLP metric. 
It then creates a gRPC connection to the Collector, and sends the metric it builds. 

