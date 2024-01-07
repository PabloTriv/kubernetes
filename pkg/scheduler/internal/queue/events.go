/*
Copyright 2019 The Kubernetes Authors.

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

package queue

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/kubernetes/pkg/api/v1/resource"
	"k8s.io/kubernetes/pkg/features"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

const (
	// PodAdd is the event when a new pod is added to API server.
	PodAdd = "PodAdd"
	// ScheduleAttemptFailure is the event when a schedule attempt fails.
	ScheduleAttemptFailure = "ScheduleAttemptFailure"
	// BackoffComplete is the event when a pod finishes backoff.
	BackoffComplete = "BackoffComplete"
	// ForceActivate is the event when a pod is moved from unschedulablePods/backoffQ
	// to activeQ. Usually it's triggered by plugin implementations.
	ForceActivate = "ForceActivate"
	// PodUpdate is the event when a pod is updated
	PodUpdate = "PodUpdate"
)

var (
	// AssignedPodAdd is the event when an assigned pod is added.
	AssignedPodAdd = framework.ClusterEvent{Resource: framework.Pod, ActionType: framework.Add, Label: "AssignedPodAdd"}
	// NodeAdd is the event when a new node is added to the cluster.
	NodeAdd = framework.ClusterEvent{Resource: framework.Node, ActionType: framework.Add, Label: "NodeAdd"}
	// NodeDelete is the event when a node is deleted from the cluster.
	NodeDelete = framework.ClusterEvent{Resource: framework.Node, ActionType: framework.Delete, Label: "NodeDelete"}
	// AssignedPodUpdate is the event when an assigned pod is updated.
	AssignedPodUpdate = framework.ClusterEvent{Resource: framework.Pod, ActionType: framework.Update, Label: "AssignedPodUpdate"}
	// UnscheduledPodAdd is the event when an unscheduled pod is added.
	UnscheduledPodAdd = framework.ClusterEvent{Resource: framework.Pod, ActionType: framework.Update, Label: "UnschedulablePodAdd"}
	// UnscheduledPodUpdate is the event when an unscheduled pod is updated.
	UnscheduledPodUpdate = framework.ClusterEvent{Resource: framework.Pod, ActionType: framework.Update, Label: "UnschedulablePodUpdate"}
	// UnscheduledPodDelete is the event when an unscheduled pod is deleted.
	UnscheduledPodDelete = framework.ClusterEvent{Resource: framework.Pod, ActionType: framework.Update, Label: "UnschedulablePodDelete"}
	// AssignedPodDelete is the event when an assigned pod is deleted.
	AssignedPodDelete = framework.ClusterEvent{Resource: framework.Pod, ActionType: framework.Delete, Label: "AssignedPodDelete"}
	// PodRequestChange is the event when a pod's resource request is changed.
	PodRequestChange = framework.ClusterEvent{Resource: framework.Pod, ActionType: framework.UpdatePodRequest, Label: "PodRequestChange"}
	// PodLabelChange is the event when a pod's label is changed.
	PodLabelChange = framework.ClusterEvent{Resource: framework.Pod, ActionType: framework.UpdatePodLabel, Label: "PodLabelChange"}
	// NodeSpecUnschedulableChange is the event when unschedulable node spec is changed.
	NodeSpecUnschedulableChange = framework.ClusterEvent{Resource: framework.Node, ActionType: framework.UpdateNodeTaint, Label: "NodeSpecUnschedulableChange"}
	// NodeAllocatableChange is the event when node allocatable is changed.
	NodeAllocatableChange = framework.ClusterEvent{Resource: framework.Node, ActionType: framework.UpdateNodeAllocatable, Label: "NodeAllocatableChange"}
	// NodeLabelChange is the event when node label is changed.
	NodeLabelChange = framework.ClusterEvent{Resource: framework.Node, ActionType: framework.UpdateNodeLabel, Label: "NodeLabelChange"}
	// NodeAnnotationChange is the event when node annotation is changed.
	NodeAnnotationChange = framework.ClusterEvent{Resource: framework.Node, ActionType: framework.UpdateNodeAnnotation, Label: "NodeAnnotationChange"}
	// NodeTaintChange is the event when node taint is changed.
	NodeTaintChange = framework.ClusterEvent{Resource: framework.Node, ActionType: framework.UpdateNodeTaint, Label: "NodeTaintChange"}
	// NodeConditionChange is the event when node condition is changed.
	NodeConditionChange = framework.ClusterEvent{Resource: framework.Node, ActionType: framework.UpdateNodeCondition, Label: "NodeConditionChange"}
	// PvAdd is the event when a persistent volume is added in the cluster.
	PvAdd = framework.ClusterEvent{Resource: framework.PersistentVolume, ActionType: framework.Add, Label: "PvAdd"}
	// PvUpdate is the event when a persistent volume is updated in the cluster.
	PvUpdate = framework.ClusterEvent{Resource: framework.PersistentVolume, ActionType: framework.Update, Label: "PvUpdate"}
	// PvcAdd is the event when a persistent volume claim is added in the cluster.
	PvcAdd = framework.ClusterEvent{Resource: framework.PersistentVolumeClaim, ActionType: framework.Add, Label: "PvcAdd"}
	// PvcUpdate is the event when a persistent volume claim is updated in the cluster.
	PvcUpdate = framework.ClusterEvent{Resource: framework.PersistentVolumeClaim, ActionType: framework.Update, Label: "PvcUpdate"}
	// StorageClassAdd is the event when a StorageClass is added in the cluster.
	StorageClassAdd = framework.ClusterEvent{Resource: framework.StorageClass, ActionType: framework.Add, Label: "StorageClassAdd"}
	// StorageClassUpdate is the event when a StorageClass is updated in the cluster.
	StorageClassUpdate = framework.ClusterEvent{Resource: framework.StorageClass, ActionType: framework.Update, Label: "StorageClassUpdate"}
	// CSINodeAdd is the event when a CSI node is added in the cluster.
	CSINodeAdd = framework.ClusterEvent{Resource: framework.CSINode, ActionType: framework.Add, Label: "CSINodeAdd"}
	// CSINodeUpdate is the event when a CSI node is updated in the cluster.
	CSINodeUpdate = framework.ClusterEvent{Resource: framework.CSINode, ActionType: framework.Update, Label: "CSINodeUpdate"}
	// CSIDriverAdd is the event when a CSI driver is added in the cluster.
	CSIDriverAdd = framework.ClusterEvent{Resource: framework.CSIDriver, ActionType: framework.Add, Label: "CSIDriverAdd"}
	// CSIDriverUpdate is the event when a CSI driver is updated in the cluster.
	CSIDriverUpdate = framework.ClusterEvent{Resource: framework.CSIDriver, ActionType: framework.Update, Label: "CSIDriverUpdate"}
	// CSIStorageCapacityAdd is the event when a CSI storage capacity is added in the cluster.
	CSIStorageCapacityAdd = framework.ClusterEvent{Resource: framework.CSIStorageCapacity, ActionType: framework.Add, Label: "CSIStorageCapacityAdd"}
	// CSIStorageCapacityUpdate is the event when a CSI storage capacity is updated in the cluster.
	CSIStorageCapacityUpdate = framework.ClusterEvent{Resource: framework.CSIStorageCapacity, ActionType: framework.Update, Label: "CSIStorageCapacityUpdate"}
	// WildCardEvent semantically matches all resources on all actions.
	WildCardEvent = framework.ClusterEvent{Resource: framework.WildCard, ActionType: framework.All, Label: "WildCardEvent"}
	// UnschedulableTimeout is the event when a pod stays in unschedulable for longer than timeout.
	UnschedulableTimeout = framework.ClusterEvent{Resource: framework.WildCard, ActionType: framework.All, Label: "UnschedulableTimeout"}
)

// PodSchedulingPropertiesChange interprets the update of a pod and returns corresponding UpdatePodXYZ event(s).
// Once we have other pod update events, we should update here as well.
func PodSchedulingPropertiesChange(newPod *v1.Pod, oldPod *v1.Pod) (events []framework.ClusterEvent) {
	podChangeExtracters := []podChangeExtractor{
		extractPodLabelsChange,
		extractPodResourceRequestChange,
	}

	for _, fn := range podChangeExtracters {
		if event := fn(newPod, oldPod); event != nil {
			events = append(events, *event)
		}
	}

	if len(events) == 0 {
		events = append(events, AssignedPodUpdate)
	}

	return
}

type podChangeExtractor func(newNode *v1.Pod, oldNode *v1.Pod) *framework.ClusterEvent

func extractPodResourceRequestChange(newPod, oldPod *v1.Pod) *framework.ClusterEvent {
	opt := resource.PodResourcesOptions{InPlacePodVerticalScalingEnabled: utilfeature.DefaultFeatureGate.Enabled(features.InPlacePodVerticalScaling)}
	if !equality.Semantic.DeepEqual(resource.PodRequests(newPod, opt), resource.PodRequests(oldPod, opt)) {
		return &PodRequestChange
	}
	return nil
}

func extractPodLabelsChange(newPod *v1.Pod, oldPod *v1.Pod) *framework.ClusterEvent {
	if isLabelChanged(newPod.GetLabels(), oldPod.GetLabels()) {
		return &PodLabelChange
	}
	return nil
}

// NodeSchedulingPropertiesChange interprets the update of a node and returns corresponding UpdateNodeXYZ event(s).
func NodeSchedulingPropertiesChange(newNode *v1.Node, oldNode *v1.Node) (events []framework.ClusterEvent) {
	nodeChangeExtracters := []nodeChangeExtractor{
		extractNodeSpecUnschedulableChange,
		extractNodeAllocatableChange,
		extractNodeLabelsChange,
		extractNodeTaintsChange,
		extractNodeConditionsChange,
		extractNodeAnnotationsChange,
	}

	for _, fn := range nodeChangeExtracters {
		if event := fn(newNode, oldNode); event != nil {
			events = append(events, *event)
		}
	}
	return
}

type nodeChangeExtractor func(newNode *v1.Node, oldNode *v1.Node) *framework.ClusterEvent

func extractNodeAllocatableChange(newNode *v1.Node, oldNode *v1.Node) *framework.ClusterEvent {
	if !equality.Semantic.DeepEqual(oldNode.Status.Allocatable, newNode.Status.Allocatable) {
		return &NodeAllocatableChange
	}
	return nil
}

func extractNodeLabelsChange(newNode *v1.Node, oldNode *v1.Node) *framework.ClusterEvent {
	if isLabelChanged(newNode.GetLabels(), oldNode.GetLabels()) {
		return &NodeLabelChange
	}
	return nil
}

func isLabelChanged(newLabels map[string]string, oldLabels map[string]string) bool {
	return !equality.Semantic.DeepEqual(newLabels, oldLabels)
}

func extractNodeTaintsChange(newNode *v1.Node, oldNode *v1.Node) *framework.ClusterEvent {
	if !equality.Semantic.DeepEqual(newNode.Spec.Taints, oldNode.Spec.Taints) {
		return &NodeTaintChange
	}
	return nil
}

func extractNodeConditionsChange(newNode *v1.Node, oldNode *v1.Node) *framework.ClusterEvent {
	strip := func(conditions []v1.NodeCondition) map[v1.NodeConditionType]v1.ConditionStatus {
		conditionStatuses := make(map[v1.NodeConditionType]v1.ConditionStatus, len(conditions))
		for i := range conditions {
			conditionStatuses[conditions[i].Type] = conditions[i].Status
		}
		return conditionStatuses
	}
	if !equality.Semantic.DeepEqual(strip(oldNode.Status.Conditions), strip(newNode.Status.Conditions)) {
		return &NodeConditionChange
	}
	return nil
}

func extractNodeSpecUnschedulableChange(newNode *v1.Node, oldNode *v1.Node) *framework.ClusterEvent {
	if newNode.Spec.Unschedulable != oldNode.Spec.Unschedulable && !newNode.Spec.Unschedulable {
		return &NodeSpecUnschedulableChange
	}
	return nil
}

func extractNodeAnnotationsChange(newNode *v1.Node, oldNode *v1.Node) *framework.ClusterEvent {
	if !equality.Semantic.DeepEqual(oldNode.GetAnnotations(), newNode.GetAnnotations()) {
		return &NodeAnnotationChange
	}
	return nil
}
