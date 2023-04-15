// Copyright 2023 Nautes Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	"context"
	"errors"
	"time"

	secretclient "github.com/nautes-labs/cluster-operator/pkg/secretclient/interface"
	clustercrd "github.com/nautes-labs/pkg/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Reconcile", func() {
	var cluster *clustercrd.Cluster

	BeforeEach(func() {
		cluster = &clustercrd.Cluster{
			ObjectMeta: v1.ObjectMeta{
				Name:      "test-cluster",
				Namespace: "nautes",
				Annotations: map[string]string{
					"ver": "1",
				},
			},
			Spec: clustercrd.ClusterSpec{
				ApiServer:   "https://127.0.0.1",
				ClusterType: clustercrd.CLUSTER_TYPE_PHYSICAL,
				ClusterKind: clustercrd.CLUSTER_KIND_KUBERNETES,
				Usage:       clustercrd.CLUSTER_USAGE_HOST,
				HostCluster: "",
			},
		}
		wantErr = nil
		wantResult = &secretclient.SyncResult{
			SecretID: "1",
		}
	})

	AfterEach(func() {
		err := k8sClient.Delete(context.Background(), cluster)
		Expect(client.IgnoreNotFound(err)).Should(BeNil())
		err = WaitForDelete(cluster)
		Expect(err).Should(BeNil())
	})

	It("sync secret", func() {
		startTime := time.Now()
		err := k8sClient.Create(context.Background(), cluster)
		Expect(err).Should(BeNil())
		cdt, err := WaitForCondition(cluster, startTime)
		Expect(err).Should(BeNil())
		Expect(string(cdt.Status)).Should(Equal("True"))

		nc := &clustercrd.Cluster{}
		err = k8sClient.Get(context.Background(), types.NamespacedName{
			Namespace: cluster.Namespace,
			Name:      cluster.Name,
		}, nc)
		Expect(err).Should(BeNil())
		Expect(nc.Status.MgtAuthStatus.SecretID).Should(Equal("1"))

	})

	It("sync failed", func() {
		startTime := time.Now()
		wantErr = errors.New("new error")
		err := k8sClient.Create(context.Background(), cluster)
		Expect(err).Should(BeNil())
		cdt, err := WaitForCondition(cluster, startTime)
		Expect(err).Should(BeNil())
		Expect(string(cdt.Status)).Should(Equal("False"))
		wantErr = nil
	})

	It("object is not legal", func() {
		startTime := time.Now()
		cluster.Spec.HostCluster = "newCluster"
		err := k8sClient.Create(context.Background(), cluster)
		Expect(err).Should(BeNil())
		cdt, err := WaitForCondition(cluster, startTime)
		Expect(err).Should(BeNil())
		Expect(string(cdt.Status)).Should(Equal("False"))

	})

	It("delete failed", func() {
		startTime := time.Now()
		time.Sleep(time.Second)
		err := k8sClient.Create(context.Background(), cluster)
		Expect(err).Should(BeNil())
		cdt, err := WaitForCondition(cluster, startTime)
		Expect(err).Should(BeNil())
		Expect(string(cdt.Status)).Should(Equal("True"))

		startTime = time.Now()
		time.Sleep(time.Second)
		wantErr = errors.New("new error")
		k8sClient.Delete(context.Background(), cluster)
		cdt, err = WaitForCondition(cluster, startTime)
		Expect(err).Should(BeNil())
		Expect(string(cdt.Status)).Should(Equal("False"))

		wantErr = nil
	})

})
