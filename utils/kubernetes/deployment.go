// Copyright 2023 Red Hat, Inc. and/or its affiliates
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

package kubernetes

import (
	"errors"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

const (
	// this const is available here https://github.com/kubernetes/kubernetes/blob/6e0cb243d57592c917fe449dde20b0e246bc66be/pkg/controller/deployment/util/deployment_util.go#L100
	// but it doesn't worth the dependency.
	deploymentMinimumReplicasUnavailable = "MinimumReplicasUnavailable"
)

// IsDeploymentAvailable verifies if the Deployment conditions match the Available status
func IsDeploymentAvailable(deployment *appsv1.Deployment) bool {
	return isDeploymentInCondition(deployment, appsv1.DeploymentAvailable, v1.ConditionTrue)
}

// IsDeploymentFailed returns true in case of Deployment not available (IsDeploymentAvailable returns false) or it has a condition of
// DeploymentReplicaFailure == true.
func IsDeploymentFailed(deployment *appsv1.Deployment) bool {
	if IsDeploymentAvailable(deployment) {
		return false
	}
	return isDeploymentInCondition(deployment, appsv1.DeploymentReplicaFailure, v1.ConditionTrue)
}

func isDeploymentInCondition(deployment *appsv1.Deployment, conditionType appsv1.DeploymentConditionType, status v1.ConditionStatus) bool {
	for _, condition := range deployment.Status.Conditions {
		if condition.Type == conditionType &&
			condition.Status == status {
			return true
		}
	}
	return false
}

// IsDeploymentMinimumReplicasUnavailable verifies if the deployment has the minimum replicas available
func IsDeploymentMinimumReplicasUnavailable(deployment *appsv1.Deployment) bool {
	for _, condition := range deployment.Status.Conditions {
		if condition.Type == appsv1.DeploymentAvailable &&
			condition.Status == v1.ConditionFalse &&
			condition.Reason == deploymentMinimumReplicasUnavailable {
			return true
		}
	}
	return false
}

// GetDeploymentUnavailabilityMessage returns a string explaining why the given deployment is unavailable. If empty, there's no replica failure.
// Note that the Deployment might be available, but a second replica failed to scale. Always check IsDeploymentAvailable.
//
// See: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#failed-deployment
func GetDeploymentUnavailabilityMessage(deployment *appsv1.Deployment) string {
	for _, condition := range deployment.Status.Conditions {
		if condition.Type == appsv1.DeploymentProgressing &&
			condition.Status == v1.ConditionFalse {
			return fmt.Sprintf("deployment %s unavailable: reason %s, message %s", deployment.Name, condition.Reason, condition.Message)
		}
		if condition.Type == appsv1.DeploymentReplicaFailure &&
			condition.Status == v1.ConditionTrue {
			return fmt.Sprintf("deployment %s unavailable: reason %s, message %s", deployment.Name, condition.Reason, condition.Message)
		}
	}
	return ""
}

// MarkDeploymentToRollout marks the given Deployment to restart now. The object must be updated.
// Code adapted from here: https://github.com/kubernetes/kubectl/blob/release-1.26/pkg/polymorphichelpers/objectrestarter.go#L44
func MarkDeploymentToRollout(deployment *appsv1.Deployment) error {
	if deployment.Spec.Paused {
		return errors.New("can't restart paused deployment (run rollout resume first)")
	}
	if deployment.Spec.Template.ObjectMeta.Annotations == nil {
		deployment.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
	return nil
}

// GetContainerByName returns a pointer to the Container within the given Deployment.
// If none found, returns nil.
// It also returns the position where the container was found, -1 if none
func GetContainerByName(name string, podSpec *v1.PodSpec) (*v1.Container, int) {
	if podSpec == nil {
		return nil, -1
	}
	for i, container := range podSpec.Containers {
		if container.Name == name {
			return &container, i
		}
	}
	return nil, -1
}

// AddOrReplaceContainer replace the existing container or add if it doesn't exist in the .spec.containers attribute
func AddOrReplaceContainer(containerName string, container v1.Container, podSpec *v1.PodSpec) {
	_, idx := GetContainerByName(containerName, podSpec)
	if idx < 0 {
		podSpec.Containers = append(podSpec.Containers, container)
	} else {
		podSpec.Containers[idx] = container
	}
}
