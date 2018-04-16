/*
Copyright 2018 The Kubernetes Authors.

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

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilfeature "k8s.io/apiserver/pkg/util/feature"

	// avoid error `servicecatalog/v1beta1 is not enabled`
	_ "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/install"
	scfeatures "github.com/kubernetes-incubator/service-catalog/pkg/features"
	fakeosb "github.com/pmorie/go-open-service-broker-client/v2/fake"

	"time"

	"fmt"

	"github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	"github.com/kubernetes-incubator/service-catalog/test/util"
	"github.com/pmorie/go-open-service-broker-client/v2/generator"
)

//
// TODO(nicholss): add more complicated integration tests.
//func TestClusterServicePlanRemovedFromCatalogAfterFiltered(t *testing.T) {
//	ct := &controllerTest{
//		t:      t,
//		broker: getTestBroker(),
//	}
//
//	ct.run(func(ct *controllerTest) {
//		removedPlan := getTestClusterServicePlanRemoved()
//		removedPlan, err := ct.client.ClusterServicePlans().Create(removedPlan)
//		if err != nil {
//			t.Fatalf("error creating ClusterServicePlan: %v", err)
//		}
//
//		err = util.WaitForClusterServicePlanToExist(ct.client, testRemovedClusterServicePlanGUID)
//		if err != nil {
//			t.Fatalf("error waiting for ClusterServicePlan to exist: %v", err)
//		}
//
//		t.Log("updating ClusterServiceClass status")
//		removedPlan.Status.RemovedFromBrokerCatalog = true
//		_, err = ct.client.ClusterServicePlans().UpdateStatus(removedPlan)
//		if err != nil {
//			t.Fatalf("error marking ClusterServicePlan as removed from catalog: %v", err)
//		}
//
//		err = util.WaitForClusterServicePlanToNotExist(ct.client, testRemovedClusterServicePlanGUID)
//		if err != nil {
//			t.Fatalf("error waiting for remove ClusterServicePlan to not exist: %v", err)
//		}
//	})
//}

func TestClusterServiceClassRemovedFromCatalogAfterFiltering(t *testing.T) {

	name := "Archonei"
	uuid := generator.IDFrom(name)

	broker := getTestBroker()
	broker.Spec.RelistDuration = &metav1.Duration{Duration: time.Millisecond * 100}
	broker.Spec.CatalogRestrictions = &v1beta1.CatalogRestrictions{
		ServiceClass: []string{"name!=" + uuid},
	}

	ct := &controllerTest{
		t:      t,
		broker: broker,
	}

	ct.run(func(ct *controllerTest) {

		ct.osbClient.CatalogReaction = &fakeosb.CatalogReaction{
			Response: getTestLargeCatalogResponse(),
		}

		err := util.WaitForClusterServiceClassToExist(ct.client, uuid)
		if err != nil {
			t.Fatalf("error waiting for ClusterServiceClass to exist: %v", err)
		}
		t.Log("class added")

		// turn on the catalog restrictions feature.
		//
		utilfeature.DefaultFeatureGate.Set(fmt.Sprintf("%v=true", scfeatures.CatalogRestrictions))
		defer utilfeature.DefaultFeatureGate.Set(fmt.Sprintf("%v=false", scfeatures.CatalogRestrictions))

		err = util.WaitForClusterServiceClassToNotExist(ct.client, uuid)
		if err != nil {
			t.Fatalf("error waiting for remove ClusterServiceClass to not exist: %v", err)
		}
		t.Log("class removed")
	})
}