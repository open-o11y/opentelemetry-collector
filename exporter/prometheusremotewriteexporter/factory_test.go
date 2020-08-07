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
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configcheck"
	"go.opentelemetry.io/collector/config/configmodels"
)

//Tests whether or not the default Exporter factory can instantiate a properly interfaced Exporter with default conditions
func TestCreateDefaultConfig(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	assert.NotNil(t, cfg, "failed to create default config")
	assert.NoError(t, configcheck.ValidateConfig(cfg))
}

//Tests whether or not a correct Metrics Exporter from the default Config parameters
func TestCreateMetricsExporter(t *testing.T) {
	factory := NewFactory()
	tests := []struct {
		name        string
		cfg         configmodels.Exporter
		params      component.ExporterCreateParams
		returnError bool
	}{
		{"success_case",
			factory.CreateDefaultConfig(),
			component.ExporterCreateParams{},
			false,
		},
	}
	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := factory.CreateMetricsExporter(context.Background(), tt.params, tt.cfg)
			if tt.returnError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}