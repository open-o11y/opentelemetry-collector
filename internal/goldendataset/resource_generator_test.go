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

package goldendataset

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/consumer/pdata"
)

func TestGenerateResource(t *testing.T) {
	resourceIds := []PICTInputResource{ResourceNil, ResourceEmpty, ResourceVMOnPrem, ResourceVMCloud, ResourceK8sOnPrem,
		ResourceK8sCloud, ResourceFaas, ResourceExec}
	for _, rscID := range resourceIds {
		rsc := GenerateResource(rscID)
		if rscID == ResourceNil || rscID == ResourceEmpty {
			assert.Equal(t, 0, rsc.Attributes().Len())
		} else {
			assert.True(t, rsc.Attributes().Len() > 0)
		}
		copy := pdata.NewResource()
		rsc.CopyTo(copy)
		assert.EqualValues(t, rsc.Attributes().Len(), copy.Attributes().Len())
	}
}
