// Copyright 2020 The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheusremotewriteexporter

import (
	"context"
	"net/http"
	"plugin"

	"github.com/pkg/errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// The value of "type" key in configuration.
	typeStr = "prometheusremotewrite"
)

func NewFactory() component.ExporterFactory {
	return exporterhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		exporterhelper.WithMetrics(createMetricsExporter))
}

// Instantiates a pseudo-Cortex Exporter that adheres to the component MetricsExporter interface
func createMetricsExporter(_ context.Context, _ component.ExporterCreateParams,
	cfg configmodels.Exporter) (component.MetricsExporter, error) {

	cCfg, ok := cfg.(*Config)
	if !ok {
		return nil, errors.Errorf("invalid configuration")
	}
	client, _ := cCfg.HTTPClientSettings.ToClient()

	if len(cCfg.AuthPath) != 0 {
		auth, err := plugin.Open(cCfg.AuthPath)
		if err != nil {
			return nil, err
		}

		newAuth, err := auth.Lookup("NewAuth")
		if err != nil {
			return nil, err
		}

		client, err = newAuth.(func(*http.Client, string) (*http.Client, error))(client, cCfg.Region)
		if err != nil {
			return nil, err
		}
	}
	prwe, err := newPrwExporter(cCfg.Namespace, cCfg.HTTPClientSettings.Endpoint, client, cCfg.Headers)
	if err != nil {
		return nil, err
	}

	prwexp, err := exporterhelper.NewMetricsExporter(
		cfg,
		prwe.pushMetrics,
		exporterhelper.WithTimeout(cCfg.TimeoutSettings),
		exporterhelper.WithQueue(cCfg.QueueSettings),
		exporterhelper.WithRetry(cCfg.RetrySettings),
		exporterhelper.WithShutdown(prwe.shutdown),
	)

	if err != nil {
		return nil, err
	}

	return prwexp, nil
}

func createDefaultConfig() configmodels.Exporter {
	// TODO: Enable the queued settings.
	qs := exporterhelper.CreateDefaultQueueSettings()
	qs.Enabled = false

	return &Config{
		ExporterSettings: configmodels.ExporterSettings{
			TypeVal: typeStr,
			NameVal: typeStr,
		},
		Namespace:       "",
		Headers:         map[string]string{},
		AuthPath:        "",
		Region:          "",
		TimeoutSettings: exporterhelper.CreateDefaultTimeoutSettings(),
		RetrySettings:   exporterhelper.CreateDefaultRetrySettings(),
		QueueSettings:   qs,
		HTTPClientSettings: confighttp.HTTPClientSettings{
			Endpoint: "http://some.url:9411/api/prom/push",
			// We almost read 0 bytes, so no need to tune ReadBufferSize.
			ReadBufferSize:  0,
			WriteBufferSize: 512 * 1024,
		},
	}
}
