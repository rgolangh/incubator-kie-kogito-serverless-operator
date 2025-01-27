// Copyright 2023 Red Hat, Inc. and/or its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"

	framework "github.com/apache/incubator-kie-kogito-serverless-operator/bddframework/pkg/framework"
	"github.com/apache/incubator-kie-kogito-serverless-operator/testbdd/executor"
	"github.com/apache/incubator-kie-kogito-serverless-operator/testbdd/meta"
)

func TestMain(m *testing.M) {
	// Create kube client
	if err := framework.InitKubeClient(meta.GetRegisteredSchema()); err != nil {
		panic(err)
	}

	executor.ExecuteBDDTests(nil)
}
