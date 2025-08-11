// Copyright 2022-2025 FLUIDOS Project
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

package indexer

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	nodecorev1alpha1 "github.com/fluidos-project/node/apis/nodecore/v1alpha1"
)

// FlavorByNodeName is a controller-runtime Manager indexer returning
// a list of nodecorev1alpha1.Flavor by Owner Reference name.
type FlavorByNodeName struct{}

func (FlavorByNodeName) Object() client.Object {
	return &nodecorev1alpha1.Flavor{}
}

func (FlavorByNodeName) Field() string {
	return "metadata.ownerReferences.name"
}

func (FlavorByNodeName) IndexerFunc() client.IndexerFunc {
	return func(obj client.Object) []string {
		ownerReferences := obj.GetOwnerReferences()
		keys := make([]string, 0, len(ownerReferences))

		for _, or := range ownerReferences {
			keys = append(keys, or.Name)
		}

		return keys
	}
}
