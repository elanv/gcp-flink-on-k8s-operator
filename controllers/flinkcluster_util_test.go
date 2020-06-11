/*
Copyright 2019 Google LLC.

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

package controllers

import (
	"github.com/googlecloudplatform/flink-operator/controllers/flinkclient"
	"github.com/googlecloudplatform/flink-operator/controllers/history"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
	"time"

	v1beta1 "github.com/googlecloudplatform/flink-operator/api/v1beta1"
	"gotest.tools/assert"
)

func TestTimeConverter(t *testing.T) {
	var tc = &TimeConverter{}

	var str1 = "2019-10-23T05:10:36Z"
	var tm1 = tc.FromString(str1)
	var str2 = tc.ToString(tm1)
	assert.Assert(t, str1 == str2)

	var str3 = "2019-10-24T09:57:18+09:00"
	var tm2 = tc.FromString(str3)
	var str4 = tc.ToString(tm2)
	assert.Assert(t, str3 == str4)
}

func TestShouldRestartJob(t *testing.T) {
	var restartOnFailure = v1beta1.JobRestartPolicyFromSavepointOnFailure
	var jobStatus1 = v1beta1.JobStatus{
		State:             v1beta1.JobStateFailed,
		SavepointLocation: "gs://my-bucket/savepoint-123",
	}
	var restart1 = shouldRestartJob(&restartOnFailure, &jobStatus1)
	assert.Equal(t, restart1, true)

	var jobStatus2 = v1beta1.JobStatus{
		State: v1beta1.JobStateFailed,
	}
	var restart2 = shouldRestartJob(&restartOnFailure, &jobStatus2)
	assert.Equal(t, restart2, false)

	var neverRestart = v1beta1.JobRestartPolicyNever
	var jobStatus3 = v1beta1.JobStatus{
		State:             v1beta1.JobStateFailed,
		SavepointLocation: "gs://my-bucket/savepoint-123",
	}
	var restart3 = shouldRestartJob(&neverRestart, &jobStatus3)
	assert.Equal(t, restart3, false)
}

func TestGetRetryCount(t *testing.T) {
	var data1 = map[string]string{}
	var result1, _ = getRetryCount(data1)
	assert.Equal(t, result1, "1")

	var data2 = map[string]string{"retries": "1"}
	var result2, _ = getRetryCount(data2)
	assert.Equal(t, result2, "2")
}

func TestNewRevision(t *testing.T) {
	var jmReplicas int32 = 1
	var rpcPort int32 = 8001
	var blobPort int32 = 8002
	var queryPort int32 = 8003
	var uiPort int32 = 8004
	var dataPort int32 = 8005
	var memoryOffHeapRatio int32 = 25
	var memoryOffHeapMin = resource.MustParse("600M")
	var parallelism int32 = 2
	var savepointDir = "/savepoint_dir"
	var flinkCluster = v1beta1.FlinkCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mycluster",
			Namespace: "default",
		},
		Spec: v1beta1.FlinkClusterSpec{
			Image: v1beta1.ImageSpec{
				Name:       "flink:1.8.1",
				PullPolicy: corev1.PullPolicy("Always"),
			},
			JobManager: v1beta1.JobManagerSpec{
				Replicas:    &jmReplicas,
				AccessScope: v1beta1.AccessScopeVPC,
				Ports: v1beta1.JobManagerPorts{
					RPC:   &rpcPort,
					Blob:  &blobPort,
					Query: &queryPort,
					UI:    &uiPort,
				},
				MemoryOffHeapRatio: &memoryOffHeapRatio,
				MemoryOffHeapMin:   memoryOffHeapMin,
			},
			TaskManager: v1beta1.TaskManagerSpec{
				Replicas: 3,
				Ports: v1beta1.TaskManagerPorts{
					RPC:   &rpcPort,
					Data:  &dataPort,
					Query: &queryPort,
				},
				MemoryOffHeapRatio: &memoryOffHeapRatio,
				MemoryOffHeapMin:   memoryOffHeapMin,
			},
			Job: &v1beta1.JobSpec{
				JarFile:       "gs://my-bucket/myjob.jar",
				Parallelism:   &parallelism,
				SavepointsDir: &savepointDir,
			},
		},
	}
	var collisionCount int32 = 0
	var controller = true
	var blockOwnerDeletion = true
	var raw, _ = getPatch(&flinkCluster)
	var revision, _ = newRevision(&flinkCluster, 1, &collisionCount)
	var expectedRevision = appsv1.ControllerRevision{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mycluster-fd56cb4fb",
			Namespace: "default",
			Labels: map[string]string{
				"flinkoperator.k8s.io/hash":       "fd56cb4fb",
				"flinkoperator.k8s.io/managed-by": "mycluster",
			},
			Annotations: map[string]string{},
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion:         "flinkoperator.k8s.io/v1beta1",
				Kind:               "FlinkCluster",
				Name:               "mycluster",
				Controller:         &controller,
				BlockOwnerDeletion: &blockOwnerDeletion,
			}},
		},
		Revision: 1,
		Data:     runtime.RawExtension{Raw: raw},
	}
	assert.Assert(t, revision != nil)
	assert.DeepEqual(
		t,
		*revision,
		expectedRevision)
}

func TestCanTakeSavepoint(t *testing.T) {
	// session cluster
	var cluster = v1beta1.FlinkCluster{
		Spec: v1beta1.FlinkClusterSpec{},
	}
	var take = canTakeSavepoint(cluster)
	assert.Equal(t, take, false)

	// no savepointDir and job status
	cluster = v1beta1.FlinkCluster{
		Spec: v1beta1.FlinkClusterSpec{
			Job: &v1beta1.JobSpec{},
		},
	}
	take = canTakeSavepoint(cluster)
	assert.Equal(t, take, false)

	// no job status, job is to be started
	savepointDir := "/savepoints"
	cluster = v1beta1.FlinkCluster{
		Spec: v1beta1.FlinkClusterSpec{
			Job: &v1beta1.JobSpec{SavepointsDir: &savepointDir},
		},
	}
	take = canTakeSavepoint(cluster)
	assert.Equal(t, take, true)

	// running job and no progressing savepoint
	savepointDir = "/savepoints"
	cluster = v1beta1.FlinkCluster{
		Spec: v1beta1.FlinkClusterSpec{
			Job: &v1beta1.JobSpec{SavepointsDir: &savepointDir},
		},
		Status: v1beta1.FlinkClusterStatus{Components: v1beta1.FlinkClusterComponentsStatus{
			Job: &v1beta1.JobStatus{State: "Running"},
		}},
	}
	take = canTakeSavepoint(cluster)
	assert.Equal(t, take, true)

	// progressing savepoint
	savepointDir = "/savepoints"
	cluster = v1beta1.FlinkCluster{
		Spec: v1beta1.FlinkClusterSpec{
			Job: &v1beta1.JobSpec{SavepointsDir: &savepointDir},
		},
		Status: v1beta1.FlinkClusterStatus{
			Components: v1beta1.FlinkClusterComponentsStatus{
				Job: &v1beta1.JobStatus{State: "Running"},
			},
			Savepoint: &v1beta1.SavepointStatus{State: v1beta1.SavepointStateInProgress},
		},
	}
	take = canTakeSavepoint(cluster)
	assert.Equal(t, take, false)
}

func TestShouldUpdateJob(t *testing.T) {
	// should update
	var tc = &TimeConverter{}
	var savepointTime = time.Now()
	var observeTime = savepointTime.Add(time.Second * 100)
	var observed = ObservedClusterState{
		observeTime: observeTime,
		cluster: &v1beta1.FlinkCluster{
			Status: v1beta1.FlinkClusterStatus{
				Components: v1beta1.FlinkClusterComponentsStatus{Job: &v1beta1.JobStatus{
					State:             v1beta1.JobStateRunning,
					LastSavepointTime: tc.ToString(savepointTime),
					SavepointLocation: "gs://my-bucket/savepoint-123",
				}},
				CurrentRevision: "1", NextRevision: "2",
			},
		},
	}
	var update = shouldUpdateJob(observed)
	assert.Equal(t, update, true)

	// should update when update triggered and job failed.
	observed = ObservedClusterState{
		cluster: &v1beta1.FlinkCluster{
			Status: v1beta1.FlinkClusterStatus{
				Components: v1beta1.FlinkClusterComponentsStatus{Job: &v1beta1.JobStatus{
					State: v1beta1.JobStateFailed,
				}},
				CurrentRevision: "1", NextRevision: "2",
			},
		},
	}
	update = shouldUpdateJob(observed)
	assert.Equal(t, update, true)

	// cannot update with old savepoint
	tc = &TimeConverter{}
	savepointTime = time.Now()
	observeTime = savepointTime.Add(time.Second * 500)
	observed = ObservedClusterState{
		observeTime: observeTime,
		cluster: &v1beta1.FlinkCluster{
			Status: v1beta1.FlinkClusterStatus{
				Components: v1beta1.FlinkClusterComponentsStatus{Job: &v1beta1.JobStatus{
					State:             v1beta1.JobStateRunning,
					LastSavepointTime: tc.ToString(savepointTime),
					SavepointLocation: "gs://my-bucket/savepoint-123",
				}},
				CurrentRevision: "1", NextRevision: "2",
			},
		},
	}
	update = shouldUpdateJob(observed)
	assert.Equal(t, update, false)

	// cannot update without savepointLocation
	tc = &TimeConverter{}
	savepointTime = time.Now()
	observeTime = savepointTime.Add(time.Second * 500)
	observed = ObservedClusterState{
		observeTime: observeTime,
		cluster: &v1beta1.FlinkCluster{
			Status: v1beta1.FlinkClusterStatus{
				Components: v1beta1.FlinkClusterComponentsStatus{Job: &v1beta1.JobStatus{
					State: v1beta1.JobStateUpdating,
				}},
				CurrentRevision: "1", NextRevision: "2",
			},
		},
	}
	update = shouldUpdateJob(observed)
	assert.Equal(t, update, false)
}

func TestGetNextRevisionNumber(t *testing.T) {
	var revisions []*appsv1.ControllerRevision
	var nextRevision = getNextRevisionNumber(revisions)
	assert.Equal(t, nextRevision, int64(1))

	revisions = []*appsv1.ControllerRevision{{Revision: 1}, {Revision: 2}}
	nextRevision = getNextRevisionNumber(revisions)
	assert.Equal(t, nextRevision, int64(3))
}

func TestIsJobTerminated(t *testing.T) {
	var jobStatus = v1beta1.JobStatus{
		State: v1beta1.JobStateSucceeded,
	}
	var terminated = isJobTerminated(nil, &jobStatus)
	assert.Equal(t, terminated, true)

	var restartOnFailure = v1beta1.JobRestartPolicyFromSavepointOnFailure
	jobStatus = v1beta1.JobStatus{
		State:             v1beta1.JobStateFailed,
		SavepointLocation: "gs://my-bucket/savepoint-123",
	}
	terminated = isJobTerminated(&restartOnFailure, &jobStatus)
	assert.Equal(t, terminated, false)
}

func TestIsSavepointUpToDate(t *testing.T) {
	var tc = &TimeConverter{}
	var savepointTime = time.Now()
	var observeTime = savepointTime.Add(time.Second * 100)
	var jobStatus = v1beta1.JobStatus{
		State:             v1beta1.JobStateFailed,
		LastSavepointTime: tc.ToString(savepointTime),
		SavepointLocation: "gs://my-bucket/savepoint-123",
	}
	var update = isSavepointUpToDate(observeTime, jobStatus)
	assert.Equal(t, update, true)

	// old
	savepointTime = time.Now()
	observeTime = savepointTime.Add(time.Second * 500)
	jobStatus = v1beta1.JobStatus{
		State:             v1beta1.JobStateFailed,
		LastSavepointTime: tc.ToString(savepointTime),
		SavepointLocation: "gs://my-bucket/savepoint-123",
	}
	update = isSavepointUpToDate(observeTime, jobStatus)
	assert.Equal(t, update, false)

	// Fails without savepointLocation
	savepointTime = time.Now()
	observeTime = savepointTime.Add(time.Second * 500)
	jobStatus = v1beta1.JobStatus{
		State:             v1beta1.JobStateFailed,
		LastSavepointTime: tc.ToString(savepointTime),
	}
	update = isSavepointUpToDate(observeTime, jobStatus)
	assert.Equal(t, update, false)
}

func TestIsComponentUpdated(t *testing.T) {
	var cluster = v1beta1.FlinkCluster{
		Status: v1beta1.FlinkClusterStatus{NextRevision: "2"},
	}
	var cluster2 = v1beta1.FlinkCluster{
		Spec: v1beta1.FlinkClusterSpec{
			JobManager: v1beta1.JobManagerSpec{Ingress: &v1beta1.JobManagerIngressSpec{}},
			Job:        &v1beta1.JobSpec{},
		},
		Status: v1beta1.FlinkClusterStatus{NextRevision: "2"},
	}
	var deploy = &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{
		history.ControllerRevisionHashLabel: "2",
	}}}
	var update = isComponentUpdated(deploy, cluster)
	assert.Equal(t, update, true)

	deploy = &appsv1.Deployment{}
	update = isComponentUpdated(deploy, cluster)
	assert.Equal(t, update, false)

	deploy = nil
	update = isComponentUpdated(deploy, cluster)
	assert.Equal(t, update, false)

	var job = &batchv1.Job{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{
		history.ControllerRevisionHashLabel: "2",
	}}}
	update = isComponentUpdated(job, cluster2)
	assert.Equal(t, update, true)

	job = &batchv1.Job{}
	update = isComponentUpdated(job, cluster2)
	assert.Equal(t, update, false)

	job = nil
	update = isComponentUpdated(job, cluster2)
	assert.Equal(t, update, false)

	job = nil
	update = isComponentUpdated(job, cluster)
	assert.Equal(t, update, true)
}

func TestIsFlinkAPIReady(t *testing.T) {
	var observed = ObservedClusterState{
		cluster: &v1beta1.FlinkCluster{
			Spec: v1beta1.FlinkClusterSpec{
				JobManager: v1beta1.JobManagerSpec{Ingress: &v1beta1.JobManagerIngressSpec{}},
				Job:        &v1beta1.JobSpec{},
			},
			Status: v1beta1.FlinkClusterStatus{NextRevision: "2"},
		},
		configMap:    &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
		jmDeployment: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
		tmDeployment: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
		jmService:    &corev1.Service{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
		flinkJobList: &flinkclient.JobStatusList{},
	}
	var ready = isFlinkAPIReady(observed)
	assert.Equal(t, ready, true)

	// flinkJobList is nil
	observed = ObservedClusterState{
		cluster: &v1beta1.FlinkCluster{
			Spec: v1beta1.FlinkClusterSpec{
				JobManager: v1beta1.JobManagerSpec{Ingress: &v1beta1.JobManagerIngressSpec{}},
				Job:        &v1beta1.JobSpec{},
			},
			Status: v1beta1.FlinkClusterStatus{NextRevision: "2"},
		},
		configMap:    &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
		jmDeployment: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
		tmDeployment: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
		jmService:    &corev1.Service{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
	}
	ready = isFlinkAPIReady(observed)
	assert.Equal(t, ready, false)

	// jmDeployment is not observed
	observed = ObservedClusterState{
		cluster: &v1beta1.FlinkCluster{
			Spec: v1beta1.FlinkClusterSpec{
				JobManager: v1beta1.JobManagerSpec{Ingress: &v1beta1.JobManagerIngressSpec{}},
				Job:        &v1beta1.JobSpec{},
			},
			Status: v1beta1.FlinkClusterStatus{NextRevision: "2"},
		},
		configMap:    &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
		tmDeployment: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
		jmService:    &corev1.Service{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
	}
	ready = isFlinkAPIReady(observed)
	assert.Equal(t, ready, false)

	// jmDeployment is not updated
	observed = ObservedClusterState{
		cluster: &v1beta1.FlinkCluster{
			Spec: v1beta1.FlinkClusterSpec{
				JobManager: v1beta1.JobManagerSpec{Ingress: &v1beta1.JobManagerIngressSpec{}},
				Job:        &v1beta1.JobSpec{},
			},
			Status: v1beta1.FlinkClusterStatus{NextRevision: "2"},
		},
		configMap:    &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
		jmDeployment: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "1"}}},
		tmDeployment: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
		jmService:    &corev1.Service{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
	}
	ready = isFlinkAPIReady(observed)
	assert.Equal(t, ready, false)
}

func TestGetUpdateState(t *testing.T) {
	var observed = ObservedClusterState{
		cluster: &v1beta1.FlinkCluster{
			Spec: v1beta1.FlinkClusterSpec{
				JobManager: v1beta1.JobManagerSpec{Ingress: &v1beta1.JobManagerIngressSpec{}},
				Job:        &v1beta1.JobSpec{},
			},
			Status: v1beta1.FlinkClusterStatus{CurrentRevision: "2", NextRevision: "3"},
		},
		job:          &batchv1.Job{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
		configMap:    &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
		jmDeployment: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
		tmDeployment: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
		jmService:    &corev1.Service{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
	}
	var state = getUpdateState(observed)
	assert.Equal(t, state, UpdateStateStoppingJob)

	observed = ObservedClusterState{
		cluster: &v1beta1.FlinkCluster{
			Spec: v1beta1.FlinkClusterSpec{
				JobManager: v1beta1.JobManagerSpec{Ingress: &v1beta1.JobManagerIngressSpec{}},
				Job:        &v1beta1.JobSpec{},
			},
			Status: v1beta1.FlinkClusterStatus{CurrentRevision: "2", NextRevision: "3"},
		},
		jmDeployment: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "3"}}},
		tmDeployment: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
		jmService:    &corev1.Service{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "2"}}},
	}
	state = getUpdateState(observed)
	assert.Equal(t, state, UpdateStateUpdating)

	observed = ObservedClusterState{
		cluster: &v1beta1.FlinkCluster{
			Spec: v1beta1.FlinkClusterSpec{
				JobManager: v1beta1.JobManagerSpec{Ingress: &v1beta1.JobManagerIngressSpec{}},
				Job:        &v1beta1.JobSpec{},
			},
			Status: v1beta1.FlinkClusterStatus{CurrentRevision: "2", NextRevision: "3"},
		},
		job:          &batchv1.Job{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "3"}}},
		configMap:    &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "3"}}},
		jmDeployment: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "3"}}},
		tmDeployment: &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "3"}}},
		jmService:    &corev1.Service{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "3"}}},
		jmIngress:    &extensionsv1beta1.Ingress{ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{history.ControllerRevisionHashLabel: "3"}}},
	}
	state = getUpdateState(observed)
	assert.Equal(t, state, UpdateStateFinished)
}
