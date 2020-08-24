// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheusremotewriteexporter

import (
	"context"
	"errors"
	"net/http"
	"plugin"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// The value of "type" key in configuration.
	typeStr       = "prometheusremotewrite"
	pluginStr     = "plugin"
	newAuthStr    = "NewAuth"
	regionStr     = "region"
	origClientStr = "origClient"
)

func NewFactory() component.ExporterFactory {
	return exporterhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		exporterhelper.WithMetrics(createMetricsExporter))
}

func createMetricsExporter(_ context.Context, _ component.ExporterCreateParams,
	cfg configmodels.Exporter) (component.MetricsExporter, error) {

	prwCfg, ok := cfg.(*Config)
	if !ok {
		return nil, errors.New("invalid configuration")
	}

	client, cerr := prwCfg.HTTPClientSettings.ToClient()
	if cerr != nil {
		return nil, cerr
	}

	// 0. check if auth plugin is present
	auth := prwCfg.AuthCfg
	if auth != nil && auth[pluginStr] != "" {
		// 1. open the so file to load the symbols
		plug, err := plugin.Open(auth[pluginStr])
		if err != nil {
			return nil, err
		}
		// 2. look up NewAuth
		newAuth, err := plug.Lookup(newAuthStr)
		if err != nil {
			return nil, err
		}
		// 3. Assert that loaded symbol is of the type func (params map[string]interface{}) (http.RoundTripper, error)
		// cannot create an alias for func (params map[string]interface{}) (http.RoundTripper, error) because will
		// cause unexpected type error.
		newAuthFunc, ok := newAuth.(func(map[string]interface{}) (http.RoundTripper, error))
		if !ok {
			return nil, errors.New("unexpected type from plugin")
		}

		// 4. use the module
		params := map[string]interface{}{
			regionStr:     prwCfg.AuthCfg[regionStr],
			origClientStr: client,
		}
		roundTripper, err := newAuthFunc(params)
		if err != nil {
			return nil, err
		}
		client.Transport = roundTripper
	}

	prwe, err := newPrwExporter(prwCfg.Namespace, prwCfg.HTTPClientSettings.Endpoint, client)

	if err != nil {
		return nil, err
	}

	prwexp, err := exporterhelper.NewMetricsExporter(
		cfg,
		prwe.pushMetrics,
		exporterhelper.WithTimeout(prwCfg.TimeoutSettings),
		exporterhelper.WithQueue(prwCfg.QueueSettings),
		exporterhelper.WithRetry(prwCfg.RetrySettings),
		exporterhelper.WithShutdown(prwe.shutdown),
	)

	return prwexp, err
}

func createDefaultConfig() configmodels.Exporter {
	qs := exporterhelper.CreateDefaultQueueSettings()
	qs.Enabled = false

	return &Config{
		ExporterSettings: configmodels.ExporterSettings{
			TypeVal: typeStr,
			NameVal: typeStr,
		},
		Namespace: "",

		TimeoutSettings: exporterhelper.CreateDefaultTimeoutSettings(),
		RetrySettings:   exporterhelper.CreateDefaultRetrySettings(),
		QueueSettings:   qs,
		HTTPClientSettings: confighttp.HTTPClientSettings{
			Endpoint: "http://some.url:9411/api/prom/push",
			// We almost read 0 bytes, so no need to tune ReadBufferSize.
			ReadBufferSize:  0,
			WriteBufferSize: 512 * 1024,
			Timeout:         exporterhelper.CreateDefaultTimeoutSettings().Timeout,
			Headers:         map[string]string{},
		},
	}
}
