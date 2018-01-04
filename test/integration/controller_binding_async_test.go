/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package integration

import (
	"testing"

	// avoid error `servicecatalog/v1beta1 is not enabled`
	_ "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/install"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/test/util"
)

// TestAsyncCreateServiceBindingSuccess successful async binding
func TestAsyncCreateServiceBindingSuccess(t *testing.T) {
	cases := []struct {
		name string
	}{
		{
			name: "defaults",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ct := &controllerTest{
				t:                t,
				broker:           getTestBroker(),
				instance:         getTestInstance(),
				binding:          getTestBinding(),
				asyncForBindings: true,
			}
			ct.run(func(ct *controllerTest) {
				condition := v1beta1.ServiceBindingCondition{
					Type:   v1beta1.ServiceBindingConditionReady,
					Status: v1beta1.ConditionTrue,
				}
				if cond, err := util.WaitForBindingCondition(ct.client, testNamespace, testBindingName, condition); err != nil {
					t.Fatalf("error waiting for binding condition: %v\n"+"expecting: %+v\n"+"last seen: %+v", err, condition, cond)
				}
			})
		})
	}
}
