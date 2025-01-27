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

package common

import "time"

const (
	RequeueAfterFailure          = 3 * time.Minute
	RequeueAfterFollowDeployment = 5 * time.Second
	RequeueAfterIsRunning        = 1 * time.Minute
	// RecoverDeploymentErrorRetries how many times the operator should try to recover from a failure before giving up
	RecoverDeploymentErrorRetries = 3
	// RequeueRecoverDeploymentErrorInterval interval between recovering from failures
	RequeueRecoverDeploymentErrorInterval = RecoverDeploymentErrorInterval * time.Minute
	RecoverDeploymentErrorInterval        = 10
)
