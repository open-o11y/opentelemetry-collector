# OTLP Load Generator
`data.go`: writes metrics in the following format to a text file:
```  		 
    name, type, label1 labelvalue1 , value1 value2 value3 value4 value5
```
gauge and counter has only one value. Output file path, type, number of metrics, labels, and value bounds of generated metrics
are all defined in `otlp.go`.

`otlp.go`: build and send metrics to the collector. After the metric text is generated, it
waits for 10 seconds and start building OTLP metrics and sending.

`util.go`: utilities for building OTLP metrics.
