// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build kubeapiserver && test

package kubeapiserver

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	fakediscovery "k8s.io/client-go/discovery/fake"
	fakeclientset "k8s.io/client-go/kubernetes/fake"

	"github.com/DataDog/datadog-agent/comp/core/config"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
)

func TestStoreGenerators(t *testing.T) {
	// Define tests
	tests := []struct {
		name                    string
		cfg                     map[string]interface{}
		expectedStoresGenerator []storeGenerator
	}{
		{
			name: "All configurations disabled",
			cfg: map[string]interface{}{
				"cluster_agent.collect_kubernetes_tags": false,
				"language_detection.reporting.enabled":  false,
				"language_detection.enabled":            false,
			},
			expectedStoresGenerator: []storeGenerator{},
		},
		{
			name: "All configurations disabled",
			cfg: map[string]interface{}{
				"cluster_agent.collect_kubernetes_tags": false,
				"language_detection.reporting.enabled":  false,
				"language_detection.enabled":            true,
			},
			expectedStoresGenerator: []storeGenerator{},
		},
		{
			name: "Kubernetes tags enabled",
			cfg: map[string]interface{}{
				"cluster_agent.collect_kubernetes_tags": true,
				"language_detection.reporting.enabled":  false,
				"language_detection.enabled":            true,
			},
			expectedStoresGenerator: []storeGenerator{newPodStore},
		},
		{
			name: "Language detection enabled",
			cfg: map[string]interface{}{
				"cluster_agent.collect_kubernetes_tags": false,
				"language_detection.reporting.enabled":  true,
				"language_detection.enabled":            true,
			},
			expectedStoresGenerator: []storeGenerator{newDeploymentStore},
		},
		{
			name: "Language detection enabled",
			cfg: map[string]interface{}{
				"cluster_agent.collect_kubernetes_tags": false,
				"language_detection.reporting.enabled":  true,
				"language_detection.enabled":            false,
			},
			expectedStoresGenerator: []storeGenerator{},
		},
		{
			name: "All configurations enabled",
			cfg: map[string]interface{}{
				"cluster_agent.collect_kubernetes_tags": true,
				"language_detection.reporting.enabled":  true,
				"language_detection.enabled":            true,
			},
			expectedStoresGenerator: []storeGenerator{newPodStore, newDeploymentStore},
		},
	}

	// Run test for each testcase
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := fxutil.Test[config.Component](t, fx.Options(
				config.MockModule(),
				fx.Replace(config.MockParams{Overrides: tt.cfg}),
			))
			expectedStores := collectResultStoreGenerator(tt.expectedStoresGenerator, cfg)
			stores := collectResultStoreGenerator(storeGenerators(cfg), cfg)

			assert.Equal(t, expectedStores, stores)
		})
	}
}

func collectResultStoreGenerator(funcs []storeGenerator, config config.Reader) []*reflectorStore {
	var stores []*reflectorStore
	for _, f := range funcs {
		_, s := f(nil, nil, config, nil)
		stores = append(stores, s)
	}
	return stores
}

func Test_metadataCollectionGVRs_WithFunctionalDiscovery(t *testing.T) {
	tests := []struct {
		name                  string
		apiServerResourceList []*metav1.APIResourceList
		expectedGVRs          []schema.GroupVersionResource
		cfg                   map[string]interface{}
	}{
		{
			name:                  "no requested resources, no resources at all!",
			apiServerResourceList: []*metav1.APIResourceList{},
			expectedGVRs:          []schema.GroupVersionResource{},
			cfg: map[string]interface{}{
				"cluster_agent.kube_metadata_collection.enabled":   true,
				"cluster_agent.kube_metadata_collection.resources": "",
			},
		},
		{
			name:                  "requested resources, but no resources at all!",
			apiServerResourceList: []*metav1.APIResourceList{},
			expectedGVRs:          []schema.GroupVersionResource{},
			cfg: map[string]interface{}{
				"cluster_agent.kube_metadata_collection.enabled":   true,
				"cluster_agent.kube_metadata_collection.resources": "apps/deployments",
			},
		},
		{
			name: "only one resource (deployments), only one version, correct resource requested",
			apiServerResourceList: []*metav1.APIResourceList{
				{
					GroupVersion: "apps/v1",
					APIResources: []metav1.APIResource{
						{
							Name:       "deployments",
							Kind:       "Deployment",
							Namespaced: true,
						},
					},
				},
			},
			expectedGVRs: []schema.GroupVersionResource{{Resource: "deployments", Group: "apps", Version: "v1"}},
			cfg: map[string]interface{}{
				"cluster_agent.kube_metadata_collection.enabled":   true,
				"cluster_agent.kube_metadata_collection.resources": "apps/deployments",
			},
		},
		{
			name: "only one resource (deployments), only one version, correct resource requested, but version is empty (with double slash)",
			apiServerResourceList: []*metav1.APIResourceList{
				{
					GroupVersion: "apps/v1",
					APIResources: []metav1.APIResource{
						{
							Name:       "deployments",
							Kind:       "Deployment",
							Namespaced: true,
						},
					},
				},
			},
			expectedGVRs: []schema.GroupVersionResource{{Resource: "deployments", Group: "apps", Version: "v1"}},
			cfg: map[string]interface{}{
				"cluster_agent.kube_metadata_collection.enabled":   true,
				"cluster_agent.kube_metadata_collection.resources": "apps//deployments",
			},
		},
		{
			name:                  "only one resource with specific version",
			apiServerResourceList: []*metav1.APIResourceList{},
			expectedGVRs:          []schema.GroupVersionResource{{Resource: "foo", Group: "g", Version: "v"}},
			cfg: map[string]interface{}{
				"cluster_agent.kube_metadata_collection.enabled":   true,
				"cluster_agent.kube_metadata_collection.resources": "g/v/foo",
			},
		},
		{
			name: "only one resource (deployments), only one version, wrong resource requested",
			apiServerResourceList: []*metav1.APIResourceList{
				{
					GroupVersion: "apps/v1",
					APIResources: []metav1.APIResource{
						{
							Name:       "deployments",
							Kind:       "Deployment",
							Namespaced: true,
						},
					},
				},
			},
			expectedGVRs: []schema.GroupVersionResource{},
			cfg: map[string]interface{}{
				"cluster_agent.kube_metadata_collection.enabled":   true,
				"cluster_agent.kube_metadata_collection.resources": "apps/daemonsets",
			},
		},
		{
			name: "multiple resources (deployments, statefulsets), multiple versions, all resources requested",
			apiServerResourceList: []*metav1.APIResourceList{
				{
					GroupVersion: "apps/v1",
					APIResources: []metav1.APIResource{
						{
							Name:       "deployments",
							Kind:       "Deployment",
							Namespaced: true,
						},
					},
				},
				{
					GroupVersion: "apps/v1beta1",
					APIResources: []metav1.APIResource{
						{
							Name:       "deployments",
							Kind:       "Deployment",
							Namespaced: true,
						},
					},
				},
				{
					GroupVersion: "apps/v1",
					APIResources: []metav1.APIResource{
						{
							Name:       "statefulsets",
							Kind:       "StatefulSet",
							Namespaced: true,
						},
					},
				},
				{
					GroupVersion: "apps/v1beta1",
					APIResources: []metav1.APIResource{
						{
							Name:       "statefulsets",
							Kind:       "StatefulSet",
							Namespaced: true,
						},
					},
				},
			},
			expectedGVRs: []schema.GroupVersionResource{
				{Resource: "deployments", Group: "apps", Version: "v1"},
				{Resource: "statefulsets", Group: "apps", Version: "v1"},
			},
			cfg: map[string]interface{}{
				"cluster_agent.kube_metadata_collection.enabled":   true,
				"cluster_agent.kube_metadata_collection.resources": "apps/deployments apps/statefulsets",
			},
		},
		{
			name: "multiple resources (deployments, statefulsets), multiple versions, only one resource requested",
			apiServerResourceList: []*metav1.APIResourceList{
				{
					GroupVersion: "apps/v1",
					APIResources: []metav1.APIResource{
						{
							Name:       "deployments",
							Kind:       "Deployment",
							Namespaced: true,
						},
					},
				},
				{
					GroupVersion: "apps/v1beta1",
					APIResources: []metav1.APIResource{
						{
							Name:       "deployments",
							Kind:       "Deployment",
							Namespaced: true,
						},
					},
				},
				{
					GroupVersion: "apps/v1",
					APIResources: []metav1.APIResource{
						{
							Name:       "statefulsets",
							Kind:       "StatefulSet",
							Namespaced: true,
						},
					},
				},
				{
					GroupVersion: "apps/v1beta1",
					APIResources: []metav1.APIResource{
						{
							Name:       "statefulsets",
							Kind:       "StatefulSet",
							Namespaced: true,
						},
					},
				},
			},
			expectedGVRs: []schema.GroupVersionResource{{Resource: "deployments", Group: "apps", Version: "v1"}},
			cfg: map[string]interface{}{
				"cluster_agent.kube_metadata_collection.enabled":   true,
				"cluster_agent.kube_metadata_collection.resources": "apps/deployments",
			},
		},
		{
			name: "multiple resources (deployments, statefulsets), multiple versions, two resources requested (one with a typo)",
			apiServerResourceList: []*metav1.APIResourceList{
				{
					GroupVersion: "apps/v1",
					APIResources: []metav1.APIResource{
						{
							Name:       "deployments",
							Kind:       "Deployment",
							Namespaced: true,
						},
					},
				},
				{
					GroupVersion: "apps/v1beta1",
					APIResources: []metav1.APIResource{
						{
							Name:       "deployments",
							Kind:       "Deployment",
							Namespaced: true,
						},
					},
				},
				{
					GroupVersion: "apps/v1",
					APIResources: []metav1.APIResource{
						{
							Name:       "statefulsets",
							Kind:       "StatefulSet",
							Namespaced: true,
						},
					},
				},
				{
					GroupVersion: "apps/v1beta1",
					APIResources: []metav1.APIResource{
						{
							Name:       "statefulsets",
							Kind:       "StatefulSet",
							Namespaced: true,
						},
					},
				},
			},
			expectedGVRs: []schema.GroupVersionResource{
				{Resource: "deployments", Group: "apps", Version: "v1"},
			},
			cfg: map[string]interface{}{
				"cluster_agent.kube_metadata_collection.enabled":   true,
				"cluster_agent.kube_metadata_collection.resources": "apps/deployments apps/statefulsetsy",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			cfg := fxutil.Test[config.Component](t, fx.Options(
				config.MockModule(),
				fx.Replace(config.MockParams{Overrides: test.cfg}),
			))

			client := fakeclientset.NewSimpleClientset()
			fakeDiscoveryClient, ok := client.Discovery().(*fakediscovery.FakeDiscovery)
			assert.Truef(t, ok, "Failed to initialise fake discovery client")

			fakeDiscoveryClient.Resources = test.apiServerResourceList

			discoveredGVRs, err := metadataCollectionGVRs(cfg, fakeDiscoveryClient)
			require.NoErrorf(t, err, "Function should not have returned an error")

			assert.Truef(t, reflect.DeepEqual(discoveredGVRs, test.expectedGVRs), "Expected %v but got %v.", test.expectedGVRs, discoveredGVRs)
		})
	}
}

func TestResourcesWithMetadataCollectionEnabled(t *testing.T) {
	tests := []struct {
		name              string
		cfg               map[string]interface{}
		expectedResources []string
	}{
		{
			name: "no resources requested",
			cfg: map[string]interface{}{
				"cluster_agent.kube_metadata_collection.enabled":   true,
				"cluster_agent.kube_metadata_collection.resources": "",
			},
			expectedResources: []string{"nodes"},
		},
		{
			name: "deployments needed for language detection should be excluded from metadata collection",
			cfg: map[string]interface{}{
				"language_detection.enabled":                       true,
				"language_detection.reporting.enabled":             true,
				"cluster_agent.kube_metadata_collection.enabled":   true,
				"cluster_agent.kube_metadata_collection.resources": "apps/daemonsets apps/deployments",
			},
			expectedResources: []string{"apps/daemonsets", "nodes"},
		},
		{
			name: "pods needed for autoscaling should be excluded from metadata collection",
			cfg: map[string]interface{}{
				"autoscaling.workload.enabled":                     true,
				"cluster_agent.kube_metadata_collection.enabled":   true,
				"cluster_agent.kube_metadata_collection.resources": "apps/daemonsets pods",
			},
			expectedResources: []string{"apps/daemonsets", "nodes"},
		},
		{
			name: "resources explicitly requested",
			cfg: map[string]interface{}{
				"cluster_agent.kube_metadata_collection.enabled":   true,
				"cluster_agent.kube_metadata_collection.resources": "apps/deployments apps/statefulsets",
			},
			expectedResources: []string{"nodes", "apps/deployments", "apps/statefulsets"},
		},
		{
			name: "namespaces needed for namespace labels as tags",
			cfg: map[string]interface{}{
				"kubernetes_namespace_labels_as_tags": map[string]string{
					"label1": "tag1",
				},
			},
			expectedResources: []string{"nodes", "namespaces"},
		},
		{
			name: "namespaces needed for namespace annotations as tags",
			cfg: map[string]interface{}{
				"kubernetes_namespace_annotations_as_tags": map[string]string{
					"annotation1": "tag1",
				},
			},
			expectedResources: []string{"nodes", "namespaces"},
		},
		{
			name: "namespaces needed for namespace labels and annotations as tags",
			cfg: map[string]interface{}{
				"kubernetes_namespace_labels_as_tags": map[string]string{
					"label1": "tag1",
				},
				"kubernetes_namespace_annotations_as_tags": map[string]string{
					"annotation1": "tag2",
				},
			},
			expectedResources: []string{"nodes", "namespaces"},
		},
		{
			name: "resources explicitly requested and also needed for namespace labels as tags",
			cfg: map[string]interface{}{
				"cluster_agent.kube_metadata_collection.enabled":   true,
				"cluster_agent.kube_metadata_collection.resources": "namespaces apps/deployments",
				"kubernetes_namespace_labels_as_tags": map[string]string{
					"label1": "tag1",
				},
			},
			expectedResources: []string{"nodes", "namespaces", "apps/deployments"}, // namespaces are not duplicated
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			cfg := fxutil.Test[config.Component](t, fx.Options(
				config.MockModule(),
				fx.Replace(config.MockParams{Overrides: test.cfg}),
			))

			assert.ElementsMatch(t, test.expectedResources, resourcesWithMetadataCollectionEnabled(cfg))
		})
	}
}
