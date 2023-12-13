module go.opentelemetry.io/collector/pdata

go 1.19

require (
	github.com/gogo/protobuf v1.3.2
	github.com/json-iterator/go v1.1.12
	github.com/stretchr/testify v1.8.2
	go.uber.org/multierr v1.11.0
	google.golang.org/grpc v1.60.0
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.16.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231002182017-d307bd883b97 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v1.0.0-rc10 // RC version scheme discovered to be alphabetical, use v1.0.0-rcv0011 instead
	v0.57.1 // Release failed, use v0.57.2
	v0.57.0 // Release failed, use v0.57.2
)
