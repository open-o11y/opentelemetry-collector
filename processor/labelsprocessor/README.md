# Labels Processor

Supported pipeline types: metrics

The labels processor can be used to add data point labels to all metrics that pass through it.
If any specified labels already exist in the metric, the value will be updated.

Please refer to [config.go](./config.go) for the config spec.

Example:

```yaml
processors:
  labels_processor:
    labels:
      - key: label1
        value: value1
      - key: label2
        value: value2
```

Refer to [config.yaml](./testdata/config.yaml) for detailed
examples on using the processor.
