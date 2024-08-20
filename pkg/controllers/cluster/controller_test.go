// Copyright 2023 Authors of kubean-io
// SPDX-License-Identifier: Apache-2.0

package cluster

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	k8stesting "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/record"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/config/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/kubean-io/kubean-api/apis"
	clusterv1alpha1 "github.com/kubean-io/kubean-api/apis/cluster/v1alpha1"
	clusteroperationv1alpha1 "github.com/kubean-io/kubean-api/apis/clusteroperation/v1alpha1"
	localartifactsetv1alpha1 "github.com/kubean-io/kubean-api/apis/localartifactset/v1alpha1"
	manifestv1alpha1 "github.com/kubean-io/kubean-api/apis/manifest/v1alpha1"
	"github.com/kubean-io/kubean-api/constants"
	clusterv1alpha1fake "github.com/kubean-io/kubean-api/generated/cluster/clientset/versioned/fake"
	clusteroperationv1alpha1fake "github.com/kubean-io/kubean-api/generated/clusteroperation/clientset/versioned/fake"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientsetfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func fetchTestingFake(obj interface{ RESTClient() rest.Interface }) *k8stesting.Fake {
	// https://stackoverflow.com/questions/69740891/mocking-errors-with-client-go-fake-client
	return reflect.Indirect(reflect.ValueOf(obj)).FieldByName("Fake").Interface().(*k8stesting.Fake)
}

func removeReactorFromTestingTake(obj interface{ RESTClient() rest.Interface }, verb, resource string) {
	if fakeObj := fetchTestingFake(obj); fakeObj != nil {
		newReactionChain := make([]k8stesting.Reactor, 0)
		fakeObj.Lock()
		defer fakeObj.Unlock()
		for i := range fakeObj.ReactionChain {
			reaction := fakeObj.ReactionChain[i]
			if simpleReaction, ok := reaction.(*k8stesting.SimpleReactor); ok && simpleReaction.Verb == verb && simpleReaction.Resource == resource {
				continue // ignore
			}
			newReactionChain = append(newReactionChain, reaction)
		}
		fakeObj.ReactionChain = newReactionChain
	}
}

func TestCompareClusterCondition(t *testing.T) {
	tests := []struct {
		name string
		args func(condA, conB clusterv1alpha1.ClusterCondition) bool
		want bool
	}{
		{
			name: "same",
			args: func(condA, condB clusterv1alpha1.ClusterCondition) bool {
				return CompareClusterCondition(condA, condB)
			},
			want: true,
		},
		{
			name: "same again",
			args: func(condA, condB clusterv1alpha1.ClusterCondition) bool {
				condA.Status = "123"
				condB.Status = "123"
				return CompareClusterCondition(condA, condB)
			},
			want: true,
		},
		{
			name: "clusterOps",
			args: func(condA, condB clusterv1alpha1.ClusterCondition) bool {
				condA.ClusterOps = "12"
				return CompareClusterCondition(condA, condB)
			},
			want: false,
		},
		{
			name: "status",
			args: func(condA, condB clusterv1alpha1.ClusterCondition) bool {
				condA.Status = "121212"
				return CompareClusterCondition(condA, condB)
			},
			want: false,
		},
		{
			name: "startTime",
			args: func(condA, condB clusterv1alpha1.ClusterCondition) bool {
				condA.StartTime = &metav1.Time{Time: time.Now()}
				return CompareClusterCondition(condA, condB)
			},
			want: false,
		},
		{
			name: "endTime",
			args: func(condA, condB clusterv1alpha1.ClusterCondition) bool {
				condA.EndTime = &metav1.Time{Time: time.Now()}
				return CompareClusterCondition(condA, condB)
			},
			want: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.args(clusterv1alpha1.ClusterCondition{}, clusterv1alpha1.ClusterCondition{}) != test.want {
				t.Fatal()
			}
		})
	}
}

func TestCompareClusterConditions(t *testing.T) {
	tests := []struct {
		name string
		args func() bool
		want bool
	}{
		{
			name: "zero length",
			args: func() bool {
				return CompareClusterConditions(nil, nil)
			},
			want: true,
		},
		{
			name: "different length",
			args: func() bool {
				return CompareClusterConditions(make([]clusterv1alpha1.ClusterCondition, 1), nil)
			},
			want: false,
		},
		{
			name: "one length",
			args: func() bool {
				return CompareClusterConditions(make([]clusterv1alpha1.ClusterCondition, 1), make([]clusterv1alpha1.ClusterCondition, 1))
			},
			want: true,
		},
		{
			name: "one length with different data",
			args: func() bool {
				condA := make([]clusterv1alpha1.ClusterCondition, 1)
				condB := make([]clusterv1alpha1.ClusterCondition, 1)
				condA[0].ClusterOps = "11"
				condB[0].ClusterOps = "22"
				return CompareClusterConditions(condA, condB)
			},
			want: false,
		},
	}

	for _, test := range tests {
		if test.args() != test.want {
			t.Fatal()
		}
	}
}

func TestSortClusterOperationsByCreation(t *testing.T) {
	controller := &Controller{
		Client:              newFakeClient(),
		ClientSet:           clientsetfake.NewSimpleClientset(),
		KubeanClusterSet:    clusterv1alpha1fake.NewSimpleClientset(),
		KubeanClusterOpsSet: clusteroperationv1alpha1fake.NewSimpleClientset(),
	}
	tests := []struct {
		name string
		args []clusteroperationv1alpha1.ClusterOperation
		want []clusteroperationv1alpha1.ClusterOperation
	}{
		{
			name: "empty slice",
			args: nil,
			want: nil,
		},
		{
			name: "unsorted slice",
			args: []clusteroperationv1alpha1.ClusterOperation{
				{ObjectMeta: metav1.ObjectMeta{CreationTimestamp: metav1.Unix(2, 0)}},
				{ObjectMeta: metav1.ObjectMeta{CreationTimestamp: metav1.Unix(1, 0)}},
				{ObjectMeta: metav1.ObjectMeta{CreationTimestamp: metav1.Unix(3, 0)}},
			},
			want: []clusteroperationv1alpha1.ClusterOperation{
				{ObjectMeta: metav1.ObjectMeta{CreationTimestamp: metav1.Unix(3, 0)}},
				{ObjectMeta: metav1.ObjectMeta{CreationTimestamp: metav1.Unix(2, 0)}},
				{ObjectMeta: metav1.ObjectMeta{CreationTimestamp: metav1.Unix(1, 0)}},
			},
		},
		{
			name: "eliminate-score",
			args: []clusteroperationv1alpha1.ClusterOperation{
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "100"}, CreationTimestamp: metav1.Unix(2, 0)}},
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "100"}, CreationTimestamp: metav1.Unix(1, 0)}},
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "100"}, CreationTimestamp: metav1.Unix(3, 0)}},
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "2"}, CreationTimestamp: metav1.Unix(300, 0)}},
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "2"}, CreationTimestamp: metav1.Unix(100, 0)}},
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "2"}, CreationTimestamp: metav1.Unix(200, 0)}},
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "101"}, CreationTimestamp: metav1.Unix(12, 0)}},
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "101"}, CreationTimestamp: metav1.Unix(11, 0)}},
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "101"}, CreationTimestamp: metav1.Unix(13, 0)}},
			},
			want: []clusteroperationv1alpha1.ClusterOperation{
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "2"}, CreationTimestamp: metav1.Unix(300, 0)}},
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "2"}, CreationTimestamp: metav1.Unix(200, 0)}},
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "2"}, CreationTimestamp: metav1.Unix(100, 0)}},
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "100"}, CreationTimestamp: metav1.Unix(3, 0)}},
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "100"}, CreationTimestamp: metav1.Unix(2, 0)}},
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "100"}, CreationTimestamp: metav1.Unix(1, 0)}},
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "101"}, CreationTimestamp: metav1.Unix(13, 0)}},
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "101"}, CreationTimestamp: metav1.Unix(12, 0)}},
				{ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{EliminateScoreAnno: "101"}, CreationTimestamp: metav1.Unix(11, 0)}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller.SortClusterOperationsByCreation(test.args)
			if !reflect.DeepEqual(test.args, test.want) {
				t.Fatal()
			}
		})
	}
}

func Test_CleanExcessClusterOps(t *testing.T) {
	OpsBackupNum := 5
	controller := &Controller{
		Client:              newFakeClient(),
		ClientSet:           clientsetfake.NewSimpleClientset(),
		KubeanClusterSet:    clusterv1alpha1fake.NewSimpleClientset(),
		KubeanClusterOpsSet: clusteroperationv1alpha1fake.NewSimpleClientset(),
	}
	exampleCluster := &clusterv1alpha1.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Cluster",
			APIVersion: "kubean.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster1",
		},
		Spec: clusterv1alpha1.Spec{
			HostsConfRef: &apis.ConfigMapRef{NameSpace: "kubean-system", Name: "hosts-a"},
			VarsConfRef:  &apis.ConfigMapRef{NameSpace: "kubean-system", Name: "vars-a"},
		},
	}
	tests := []struct {
		name string
		args func() bool
		want bool
	}{
		{
			name: "get nothing",
			args: func() bool {
				result, _ := controller.CleanExcessClusterOps(exampleCluster, OpsBackupNum)
				return result
			},
			want: false,
		},
		{
			name: "OpsBackupNum clusterOperations",
			args: func() bool {
				for i := 0; i < OpsBackupNum; i++ {
					clusterOperation := &clusteroperationv1alpha1.ClusterOperation{
						TypeMeta: metav1.TypeMeta{
							Kind:       "ClusterOperation",
							APIVersion: "kubean.io/v1alpha1",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:              "my_kubean_ops_cluster_1_" + fmt.Sprint(i),
							Labels:            map[string]string{constants.KubeanClusterLabelKey: "cluster1"},
							CreationTimestamp: metav1.Unix(int64(i), 0),
						},
					}
					controller.KubeanClusterOpsSet.KubeanV1alpha1().ClusterOperations().Create(context.Background(), clusterOperation, metav1.CreateOptions{})
				}
				result, _ := controller.CleanExcessClusterOps(exampleCluster, OpsBackupNum)
				return result
			},
			want: false,
		},
		{
			name: "clean clusterOperations",
			args: func() bool {
				for i := 0; i < OpsBackupNum*2; i++ {
					clusterOperation := &clusteroperationv1alpha1.ClusterOperation{
						TypeMeta: metav1.TypeMeta{
							Kind:       "ClusterOperation",
							APIVersion: "kubean.io/v1alpha1",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:              "my_kubean_ops_cluster_2_" + fmt.Sprint(i),
							Labels:            map[string]string{constants.KubeanClusterLabelKey: "cluster1"},
							CreationTimestamp: metav1.Unix(int64(i), 0),
						},
					}
					controller.KubeanClusterOpsSet.KubeanV1alpha1().ClusterOperations().Create(context.Background(), clusterOperation, metav1.CreateOptions{})
				}
				clusterOperationRunning := &clusteroperationv1alpha1.ClusterOperation{
					TypeMeta: metav1.TypeMeta{
						Kind:       "ClusterOperation",
						APIVersion: "kubean.io/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:              "my_kubean_ops_cluster_2_" + "running",
						Labels:            map[string]string{constants.KubeanClusterLabelKey: "cluster1"},
						CreationTimestamp: metav1.Unix(int64(10), 0),
					},
				}
				clusterOperationRunning, _ = controller.KubeanClusterOpsSet.KubeanV1alpha1().ClusterOperations().Create(context.Background(), clusterOperationRunning, metav1.CreateOptions{})
				clusterOperationRunning.Status = clusteroperationv1alpha1.Status{
					Status: clusteroperationv1alpha1.RunningStatus,
				}
				controller.KubeanClusterOpsSet.KubeanV1alpha1().ClusterOperations().UpdateStatus(context.Background(), clusterOperationRunning, metav1.UpdateOptions{})
				result, _ := controller.CleanExcessClusterOps(exampleCluster, OpsBackupNum)
				return result
			},
			want: true,
		},
		{
			name: "get error",
			args: func() bool {
				fetchTestingFake(controller.KubeanClusterOpsSet.KubeanV1alpha1()).PrependReactor("list", "clusteroperations", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, fmt.Errorf("this is error")
				})
				_, err := controller.CleanExcessClusterOps(exampleCluster, 5)
				removeReactorFromTestingTake(controller.KubeanClusterOpsSet.KubeanV1alpha1(), "list", "clusteroperations")
				return err != nil && err.Error() == "this is error"
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

func Test_UpdateStatus(t *testing.T) {
	controller := &Controller{
		Client:              newFakeClient(),
		ClientSet:           clientsetfake.NewSimpleClientset(),
		KubeanClusterSet:    clusterv1alpha1fake.NewSimpleClientset(),
		KubeanClusterOpsSet: clusteroperationv1alpha1fake.NewSimpleClientset(),
	}
	tests := []struct {
		name string
		args func() bool
		want bool
	}{
		{
			name: "get nothing",
			args: func() bool {
				exampleCluster := &clusterv1alpha1.Cluster{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Cluster",
						APIVersion: "kubean.io/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster1",
					},
					Spec: clusterv1alpha1.Spec{
						HostsConfRef: &apis.ConfigMapRef{NameSpace: "kubean-system", Name: "hosts-a"},
						VarsConfRef:  &apis.ConfigMapRef{NameSpace: "kubean-system", Name: "vars-a"},
					},
				}
				return controller.UpdateStatus(exampleCluster) == nil
			},
			want: true,
		},
		{
			name: "get some ops",
			args: func() bool {
				exampleCluster := &clusterv1alpha1.Cluster{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Cluster",
						APIVersion: "kubean.io/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster1",
					},
					Spec: clusterv1alpha1.Spec{
						HostsConfRef: &apis.ConfigMapRef{NameSpace: "kubean-system", Name: "hosts-a"},
						VarsConfRef:  &apis.ConfigMapRef{NameSpace: "kubean-system", Name: "vars-a"},
					},
				}
				clusterOps := &clusteroperationv1alpha1.ClusterOperation{
					TypeMeta: metav1.TypeMeta{
						Kind:       "ClusterOperation",
						APIVersion: "kubean.io/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:   "cluster1-ops",
						Labels: map[string]string{constants.KubeanClusterLabelKey: exampleCluster.Name},
					},
				}
				controller.Client.Create(context.Background(), clusterOps)
				controller.Client.Create(context.Background(), exampleCluster)
				controller.KubeanClusterSet.KubeanV1alpha1().Clusters().Create(context.Background(), exampleCluster, metav1.CreateOptions{})
				controller.KubeanClusterOpsSet.KubeanV1alpha1().ClusterOperations().Create(context.Background(), clusterOps, metav1.CreateOptions{})
				return controller.UpdateStatus(exampleCluster) == nil
			},
			want: true,
		},
		{
			name: "get one error",
			args: func() bool {
				exampleCluster := &clusterv1alpha1.Cluster{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Cluster",
						APIVersion: "kubean.io/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster1",
					},
					Spec: clusterv1alpha1.Spec{
						HostsConfRef: &apis.ConfigMapRef{NameSpace: "kubean-system", Name: "hosts-a"},
						VarsConfRef:  &apis.ConfigMapRef{NameSpace: "kubean-system", Name: "vars-a"},
					},
				}
				fetchTestingFake(controller.KubeanClusterOpsSet.KubeanV1alpha1()).PrependReactor("list", "clusteroperations", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, fmt.Errorf("this is error")
				})
				err := controller.UpdateStatus(exampleCluster)
				removeReactorFromTestingTake(controller.KubeanClusterOpsSet.KubeanV1alpha1(), "list", "clusteroperations")
				return err != nil && err.Error() == "this is error"
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

func newFakeClient() client.Client {
	sch := scheme.Scheme
	if err := clusteroperationv1alpha1.AddToScheme(sch); err != nil {
		panic(err)
	}
	if err := clusterv1alpha1.AddToScheme(sch); err != nil {
		panic(err)
	}
	client := fake.NewClientBuilder().WithScheme(sch).WithRuntimeObjects(&clusteroperationv1alpha1.ClusterOperation{}).WithRuntimeObjects(&clusterv1alpha1.Cluster{}).Build()
	return client
}

func TestReconcile(t *testing.T) {
	genController := func() *Controller {
		return &Controller{
			Client:              newFakeClient(),
			ClientSet:           clientsetfake.NewSimpleClientset(),
			KubeanClusterSet:    clusterv1alpha1fake.NewSimpleClientset(),
			KubeanClusterOpsSet: clusteroperationv1alpha1fake.NewSimpleClientset(),
		}
	}
	tests := []struct {
		name        string
		args        func() bool
		needRequeue bool
	}{
		{
			name: "cluster not found",
			args: func() bool {
				controller := genController()
				result, _ := controller.Reconcile(context.Background(), controllerruntime.Request{NamespacedName: types.NamespacedName{Name: "cluster1"}})
				return result.Requeue
			},
			needRequeue: false,
		},
		{
			name: "cluster found successfully",
			args: func() bool {
				controller := genController()
				exampleCluster := &clusterv1alpha1.Cluster{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Cluster",
						APIVersion: "kubean.io/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster1",
					},
					Spec: clusterv1alpha1.Spec{
						HostsConfRef: &apis.ConfigMapRef{NameSpace: "kubean-system", Name: "hosts-a"},
						VarsConfRef:  &apis.ConfigMapRef{NameSpace: "kubean-system", Name: "vars-a"},
					},
				}
				clusterOps := &clusteroperationv1alpha1.ClusterOperation{
					TypeMeta: metav1.TypeMeta{
						Kind:       "ClusterOperation",
						APIVersion: "kubean.io/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:   "cluster1-ops",
						Labels: map[string]string{constants.KubeanClusterLabelKey: exampleCluster.Name},
					},
				}
				controller.Client.Create(context.Background(), clusterOps)
				controller.Client.Create(context.Background(), exampleCluster)
				controller.KubeanClusterSet.KubeanV1alpha1().Clusters().Create(context.Background(), exampleCluster, metav1.CreateOptions{})
				controller.KubeanClusterOpsSet.KubeanV1alpha1().ClusterOperations().Create(context.Background(), clusterOps, metav1.CreateOptions{})
				result, _ := controller.Reconcile(context.Background(), controllerruntime.Request{NamespacedName: types.NamespacedName{Name: "cluster1"}})
				return result.RequeueAfter > 0
			},
			needRequeue: true,
		},
		{
			name: "CleanExcessClusterOps with error",
			args: func() bool {
				controller := genController()
				exampleCluster := &clusterv1alpha1.Cluster{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Cluster",
						APIVersion: "kubean.io/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster1",
					},
					Spec: clusterv1alpha1.Spec{
						HostsConfRef: &apis.ConfigMapRef{NameSpace: "kubean-system", Name: "hosts-a"},
						VarsConfRef:  &apis.ConfigMapRef{NameSpace: "kubean-system", Name: "vars-a"},
					},
				}
				clusterOps := &clusteroperationv1alpha1.ClusterOperation{
					TypeMeta: metav1.TypeMeta{
						Kind:       "ClusterOperation",
						APIVersion: "kubean.io/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:   "cluster1-ops",
						Labels: map[string]string{constants.KubeanClusterLabelKey: exampleCluster.Name},
					},
				}
				controller.Client.Create(context.Background(), clusterOps)
				controller.Client.Create(context.Background(), exampleCluster)
				controller.KubeanClusterSet.KubeanV1alpha1().Clusters().Create(context.Background(), exampleCluster, metav1.CreateOptions{})
				controller.KubeanClusterOpsSet.KubeanV1alpha1().ClusterOperations().Create(context.Background(), clusterOps, metav1.CreateOptions{})
				fetchTestingFake(controller.KubeanClusterOpsSet.KubeanV1alpha1()).PrependReactor("list", "clusteroperations", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, fmt.Errorf("this is error")
				})
				result, _ := controller.Reconcile(context.Background(), controllerruntime.Request{NamespacedName: types.NamespacedName{Name: "cluster1"}})
				removeReactorFromTestingTake(controller.KubeanClusterOpsSet.KubeanV1alpha1(), "list", "clusteroperations")
				return result.RequeueAfter > 0
			},
			needRequeue: true,
		},
		{
			name: "CleanExcessClusterOps with true",
			args: func() bool {
				controller := genController()
				exampleCluster := &clusterv1alpha1.Cluster{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Cluster",
						APIVersion: "kubean.io/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "cluster1",
					},
					Spec: clusterv1alpha1.Spec{
						HostsConfRef: &apis.ConfigMapRef{NameSpace: "kubean-system", Name: "hosts-a"},
						VarsConfRef:  &apis.ConfigMapRef{NameSpace: "kubean-system", Name: "vars-a"},
					},
				}
				clusterOps := &clusteroperationv1alpha1.ClusterOperation{
					TypeMeta: metav1.TypeMeta{
						Kind:       "ClusterOperation",
						APIVersion: "kubean.io/v1alpha1",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:   "cluster1-ops",
						Labels: map[string]string{constants.KubeanClusterLabelKey: exampleCluster.Name},
					},
				}
				controller.Client.Create(context.Background(), clusterOps)
				controller.Client.Create(context.Background(), exampleCluster)
				controller.KubeanClusterSet.KubeanV1alpha1().Clusters().Create(context.Background(), exampleCluster, metav1.CreateOptions{})
				controller.KubeanClusterOpsSet.KubeanV1alpha1().ClusterOperations().Create(context.Background(), clusterOps, metav1.CreateOptions{})
				for i := 0; i < 100; i++ {
					clusterOperation := &clusteroperationv1alpha1.ClusterOperation{
						TypeMeta: metav1.TypeMeta{
							Kind:       "ClusterOperation",
							APIVersion: "kubean.io/v1alpha1",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:              "my_kubean_ops_cluster_2_" + fmt.Sprint(i),
							Labels:            map[string]string{constants.KubeanClusterLabelKey: "cluster1"},
							CreationTimestamp: metav1.Unix(int64(i), 0),
						},
					}
					controller.KubeanClusterOpsSet.KubeanV1alpha1().ClusterOperations().Create(context.Background(), clusterOperation, metav1.CreateOptions{})
				}
				result, _ := controller.Reconcile(context.Background(), controllerruntime.Request{NamespacedName: types.NamespacedName{Name: "cluster1"}})
				return result.RequeueAfter > 0
			},
			needRequeue: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.args() != test.needRequeue {
				t.Fatal()
			}
		})
	}
}

func TestStart(t *testing.T) {
	controller := &Controller{
		Client:              newFakeClient(),
		ClientSet:           clientsetfake.NewSimpleClientset(),
		KubeanClusterSet:    clusterv1alpha1fake.NewSimpleClientset(),
		KubeanClusterOpsSet: clusteroperationv1alpha1fake.NewSimpleClientset(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	controller.Start(ctx)
}

func TestSetupWithManager(t *testing.T) {
	controller := &Controller{
		Client:              newFakeClient(),
		ClientSet:           clientsetfake.NewSimpleClientset(),
		KubeanClusterSet:    clusterv1alpha1fake.NewSimpleClientset(),
		KubeanClusterOpsSet: clusteroperationv1alpha1fake.NewSimpleClientset(),
	}
	if controller.SetupWithManager(MockManager{}) != nil {
		t.Fatal()
	}
}

type MockClusterForManager struct {
	_ string
}

func (MockClusterForManager) SetFields(interface{}) error { return nil }

func (MockClusterForManager) GetConfig() *rest.Config { return &rest.Config{} }

func (MockClusterForManager) GetScheme() *runtime.Scheme {
	sch := scheme.Scheme
	if err := manifestv1alpha1.AddToScheme(sch); err != nil {
		panic(err)
	}
	if err := localartifactsetv1alpha1.AddToScheme(sch); err != nil {
		panic(err)
	}
	return sch
}

func (MockClusterForManager) GetClient() client.Client { return nil }

func (MockClusterForManager) GetFieldIndexer() client.FieldIndexer { return nil }

func (MockClusterForManager) GetCache() cache.Cache { return nil }

func (MockClusterForManager) GetEventRecorderFor(name string) record.EventRecorder { return nil }

func (MockClusterForManager) GetRESTMapper() meta.RESTMapper { return nil }

func (MockClusterForManager) GetAPIReader() client.Reader { return nil }

func (MockClusterForManager) Start(ctx context.Context) error { return nil }

type MockManager struct {
	MockClusterForManager
}

func (MockManager) Add(manager.Runnable) error { return nil }

func (MockManager) Elected() <-chan struct{} { return nil }

func (MockManager) AddMetricsExtraHandler(path string, handler http.Handler) error { return nil }

func (MockManager) AddHealthzCheck(name string, check healthz.Checker) error { return nil }

func (MockManager) AddReadyzCheck(name string, check healthz.Checker) error { return nil }

func (MockManager) Start(ctx context.Context) error { return nil }

func (MockManager) GetWebhookServer() *webhook.Server { return nil }

func (MockManager) GetLogger() logr.Logger { return logr.Logger{} }

func (MockManager) GetControllerOptions() v1alpha1.ControllerConfigurationSpec {
	return v1alpha1.ControllerConfigurationSpec{}
}
