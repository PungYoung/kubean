// Copyright 2023 Authors of kubean-io
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	clientsetfake "k8s.io/client-go/kubernetes/fake"

	"github.com/kubean-io/kubean-api/apis"
	clusterv1alpha1 "github.com/kubean-io/kubean-api/apis/cluster/v1alpha1"
	clusteroperationv1alpha1 "github.com/kubean-io/kubean-api/apis/clusteroperation/v1alpha1"
	localartifactsetv1alpha1 "github.com/kubean-io/kubean-api/apis/localartifactset/v1alpha1"
	manifestv1alpha1 "github.com/kubean-io/kubean-api/apis/manifest/v1alpha1"
	"github.com/kubean-io/kubean-api/cluster"
	"github.com/kubean-io/kubean-api/constants"
)

func TestNewSchema(t *testing.T) {
	aggregatedScheme := NewSchema()
	tests := []struct {
		name    string
		obj     runtime.Object
		wantGVK schema.GroupVersionKind
	}{
		{
			name:    "KuBeanInfoManifest gvk",
			obj:     &manifestv1alpha1.Manifest{},
			wantGVK: schema.GroupVersionKind{Group: "kubean.io", Version: "v1alpha1", Kind: "Manifest"},
		},
		{
			name:    "LocalArtifactSet gvk",
			obj:     &localartifactsetv1alpha1.LocalArtifactSet{},
			wantGVK: schema.GroupVersionKind{Group: "kubean.io", Version: "v1alpha1", Kind: "LocalArtifactSet"},
		},
		{
			name:    "Cluster gvk",
			obj:     &clusterv1alpha1.Cluster{},
			wantGVK: schema.GroupVersionKind{Group: "kubean.io", Version: "v1alpha1", Kind: "Cluster"},
		},
		{
			name:    "ClusterOperation gvk",
			obj:     &clusteroperationv1alpha1.ClusterOperation{},
			wantGVK: schema.GroupVersionKind{Group: "kubean.io", Version: "v1alpha1", Kind: "ClusterOperation"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if gvkArray, _, _ := aggregatedScheme.ObjectKinds(test.obj); gvkArray[0] != test.wantGVK {
				t.Fatal()
			}
		})
	}
}

func TestGetCurrentNSOrDefault(t *testing.T) {
	if namespace := GetCurrentNSOrDefault(); namespace == "" {
		t.Fatal()
	}
}

func TestGetCurrentNS(t *testing.T) {
	tests := []struct {
		name string
		args func() string
		want string
	}{
		{
			name: "get nothing",
			args: func() string {
				os.Setenv("POD_NAMESPACE", "")
				ns, _ := GetCurrentNS()
				return ns
			},
			want: "",
		},
		{
			name: "get from env",
			args: func() string {
				os.Setenv("POD_NAMESPACE", "pod-name-space-123")
				ns, _ := GetCurrentNS()
				return ns
			},
			want: "pod-name-space-123",
		},
		{
			name: "get from file",
			args: func() string {
				os.Setenv("POD_NAMESPACE", "")
				tempFile, err := os.CreateTemp(os.TempDir(), "kubean-test")
				if err != nil {
					return ""
				}
				tempFilePath, err := filepath.Abs(tempFile.Name())
				if err != nil {
					return ""
				}
				tempFile.WriteString("abc-namespace-123")
				tempFile.Sync()
				tempFile.Close()
				defer os.Remove(tempFilePath)
				ServiceAccountNamespaceFile = tempFilePath
				ns, _ := GetCurrentNS()
				return ns
			},
			want: "abc-namespace-123",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.args() != test.want {
				t.Fatal()
			}
		})
	}
}

func TestGetCurrentRunningPodName(t *testing.T) {
	os.Setenv("HOSTNAME", "abc")
	if GetCurrentRunningPodName() != "abc" {
		t.Fatal()
	}
}

func TestUpdateOwnReference(t *testing.T) {
	tests := []struct {
		name string
		args func() bool
		want bool
	}{
		{
			name: "already ownerReferences exist",
			args: func() bool {
				fakeClient := clientsetfake.NewSimpleClientset()
				configMapList := []*apis.ConfigMapRef{{Name: "abc", NameSpace: "abc"}}
				fakeClient.CoreV1().ConfigMaps(configMapList[0].NameSpace).Create(context.Background(), &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:            configMapList[0].Name,
						OwnerReferences: []metav1.OwnerReference{{}},
					},
				}, metav1.CreateOptions{})
				secretList := []*apis.SecretRef{{Name: "cba", NameSpace: "cba"}}
				fakeClient.CoreV1().Secrets(secretList[0].NameSpace).Create(context.Background(), &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:            secretList[0].Name,
						OwnerReferences: []metav1.OwnerReference{{}},
					},
				}, metav1.CreateOptions{})
				return UpdateOwnReference(fakeClient, configMapList, secretList, metav1.OwnerReference{}) == nil
			},
			want: true,
		},
		{
			name: "empty RefData",
			args: func() bool {
				fakeClient := clientsetfake.NewSimpleClientset()
				configMapList := []*apis.ConfigMapRef{{Name: "", NameSpace: ""}}
				secretList := []*apis.SecretRef{{Name: "", NameSpace: ""}}
				return UpdateOwnReference(fakeClient, configMapList, secretList, metav1.OwnerReference{}) == nil
			},
			want: true,
		},
		{
			name: "not found",
			args: func() bool {
				fakeClient := clientsetfake.NewSimpleClientset()
				configMapList := []*apis.ConfigMapRef{{Name: "cm1", NameSpace: "abc"}}
				secretList := []*apis.SecretRef{{Name: "secret1", NameSpace: "abc"}}
				return UpdateOwnReference(fakeClient, configMapList, secretList, metav1.OwnerReference{}) == nil
			},
			want: true,
		},
		{
			name: "record exists",
			args: func() bool {
				fakeClient := clientsetfake.NewSimpleClientset()
				fakeClient.CoreV1().ConfigMaps("abc").Create(context.Background(), &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
					Name:      "cm1",
					Namespace: "abc",
				}}, metav1.CreateOptions{})
				fakeClient.CoreV1().Secrets("abc").Create(context.Background(), &corev1.Secret{ObjectMeta: metav1.ObjectMeta{
					Name:      "secret1",
					Namespace: "abc",
				}}, metav1.CreateOptions{})
				configMapList := []*apis.ConfigMapRef{{Name: "cm1", NameSpace: "abc"}}
				secretList := []*apis.SecretRef{{Name: "secret1", NameSpace: "abc"}}
				return UpdateOwnReference(fakeClient, configMapList, secretList, metav1.OwnerReference{}) == nil
			},
			want: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.args() != test.want {
				t.Fatal()
			}
		})
	}
}

func TestFetchKubeanConfigProperty(t *testing.T) {
	genClient := func() kubernetes.Interface {
		return clientsetfake.NewSimpleClientset()
	}
	tests := []struct {
		name   string
		client kubernetes.Interface
		args   func(client kubernetes.Interface)
		want   cluster.ConfigProperty
	}{
		{
			name:   "no config and default value",
			client: genClient(),
			args: func(client kubernetes.Interface) {
				// do nothing
			},
			want: cluster.ConfigProperty{},
		},
		{
			name:   "the config has been set to 10000",
			client: genClient(),
			args: func(client kubernetes.Interface) {
				os.Setenv("POD_NAMESPACE", "")
				configMap := &corev1.ConfigMap{}
				configMap.Name = constants.KubeanConfigMapName
				configMap.Namespace = GetCurrentNSOrDefault()
				configMap.Data = map[string]string{"CLUSTER_OPERATIONS_BACKEND_LIMIT": "10000"}
				client.CoreV1().ConfigMaps(GetCurrentNSOrDefault()).Create(context.Background(), configMap, metav1.CreateOptions{})
			},
			want: cluster.ConfigProperty{
				ClusterOperationsBackEndLimit: "10000",
				SprayJobImageRegistry:         "",
			},
		},
		{
			name:   "the config has been set to 10000",
			client: genClient(),
			args: func(client kubernetes.Interface) {
				os.Setenv("POD_NAMESPACE", "")
				configMap := &corev1.ConfigMap{}
				configMap.Name = constants.KubeanConfigMapName
				configMap.Namespace = GetCurrentNSOrDefault()
				configMap.Data = map[string]string{"CLUSTER_OPERATIONS_BACKEND_LIMIT": "10000", "SPRAY_JOB_IMAGE_REGISTRY": "ghcr.io"}
				client.CoreV1().ConfigMaps(GetCurrentNSOrDefault()).Create(context.Background(), configMap, metav1.CreateOptions{})
			},
			want: cluster.ConfigProperty{
				ClusterOperationsBackEndLimit: "10000",
				SprayJobImageRegistry:         "ghcr.io",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.args(test.client)
			result := FetchKubeanConfigProperty(test.client)
			if result.ClusterOperationsBackEndLimit != test.want.ClusterOperationsBackEndLimit || result.SprayJobImageRegistry != test.want.SprayJobImageRegistry {
				t.Fatal()
			}
		})
	}
}

func TestGetClusterOperationsBackEndLimit(t *testing.T) {
	tests := []struct {
		name     string
		limit    string
		expected int
	}{
		{
			name:     "valid limit",
			limit:    "5",
			expected: 5,
		},
		{
			name:     "invalid limit, use default",
			limit:    "0",
			expected: constants.DefaultClusterOperationsBackEndLimit,
		},
		{
			name:     "invalid limit, use max",
			limit:    "300",
			expected: constants.MaxClusterOperationsBackEndLimit,
		},
		{
			name:     "invalid format, use default",
			limit:    "abc",
			expected: constants.DefaultClusterOperationsBackEndLimit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &cluster.ConfigProperty{
				ClusterOperationsBackEndLimit: tt.limit,
			}
			if got := config.GetClusterOperationsBackEndLimit(); got != tt.expected {
				t.Errorf("GetClusterOperationsBackEndLimit() = %v, want %v", got, tt.expected)
			}
		})
	}
}
