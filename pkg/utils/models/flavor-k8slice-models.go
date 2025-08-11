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

package models

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/api/resource"
)

// GpuCharacteristics represents the characteristics of a Gpu.
type GpuCharacteristics struct {
	Vendor                string            `json:"vendor"`
	Model                 string            `json:"model"`
	Count                 int64             `json:"total_count"`
	Tier                  string            `json:"tier"`
	MultiInstance         bool              `json:"multi_instance"`
	Shared                bool              `json:"shared"`
	SharingStrategy       string            `json:"sharing_strategy"`
	Dedicated             bool              `json:"dedicated"`
	Interruptible         bool              `json:"interruptible"`
	NetworkBandwidth      resource.Quantity `json:"network_bandwidth"`
	NetworkLatencyMs      int64             `json:"network_latency_ms"`
	NetworkTier           string            `json:"network_tier"`
	TrainingScore         float64           `json:"training_score"`
	InferenceScore        float64           `json:"inference_score"`
	HPCScore              float64           `json:"hpc_score"`
	GraphicsScore         float64           `json:"graphics_score"`
	Architecture          string            `json:"architecture"`
	Interconnect          string            `json:"interconnect"`
	InterconnectBandwidth resource.Quantity `json:"interconnect_bandwidth"`
	Cores                 resource.Quantity `json:"cores"`
	Memory                resource.Quantity `json:"memory"`
	ComputeCapability     string            `json:"compute_capability"`
	ClockSpeed            resource.Quantity `json:"clock_speed"`
	FP32TFlops            float64           `json:"fp32_tflops"`
	Topology              string            `json:"topology"`
	MultiGPUEfficiency    string            `json:"multi_gpu_efficiency"`
	Region                string            `json:"region"`
	Zone                  string            `json:"zone"`
	HourlyRate            float64           `json:"hourly_rate"`
	Provider              string            `json:"provider"`
	PreEmptible           bool              `json:"pre_emptible"`
}

// K8SliceCharacteristics represents the characteristics of a Kubernetes slice.
type K8SliceCharacteristics struct {
	Architecture string              `json:"architecture"`
	CPU          resource.Quantity   `json:"cpu"`
	Memory       resource.Quantity   `json:"memory"`
	Pods         resource.Quantity   `json:"pods"`
	Gpu          *GpuCharacteristics `json:"gpu,omitempty"`
	Storage      *resource.Quantity  `json:"storage,omitempty"`
}

// K8SliceProperties represents the properties of a Kubernetes slice.
type K8SliceProperties struct {
	Latency               int                        `json:"latency,omitempty"`
	SecurityStandards     []string                   `json:"securityStandards,omitempty"`
	CarbonFootprint       *CarbonFootprint           `json:"carbonFootprint,omitempty"`
	NetworkAuthorizations *NetworkAuthorizations     `json:"networkAuthorizations,omitempty"`
	AdditionalProperties  map[string]json.RawMessage `json:"additionalProperties,omitempty"`
}

// K8SlicePartitionability represents the partitionability of a Kubernetes slice.
type K8SlicePartitionability struct {
	CPUMin     resource.Quantity `json:"cpuMin"`
	MemoryMin  resource.Quantity `json:"memoryMin"`
	PodsMin    resource.Quantity `json:"podsMin"`
	CPUStep    resource.Quantity `json:"cpuStep"`
	MemoryStep resource.Quantity `json:"memoryStep"`
	PodsStep   resource.Quantity `json:"podsStep"`
}

// K8SlicePolicies represents the policies of a Kubernetes slice.
type K8SlicePolicies struct {
	Partitionability K8SlicePartitionability `json:"partitionability"`
}

// K8Slice represents a Kubernetes slice.
type K8Slice struct {
	Characteristics K8SliceCharacteristics `json:"charateristics"`
	Properties      K8SliceProperties      `json:"properties"`
	Policies        K8SlicePolicies        `json:"policies"`
}

// GetFlavorTypeName returns the type of the Flavor.
func (K8Slice) GetFlavorTypeName() FlavorTypeName {
	return K8SliceNameDefault
}

// K8SliceSelector represents the criteria for selecting a K8Slice Flavor.
type K8SliceSelector struct {
	Architecture *StringFilter           `json:"architecture,omitempty"`
	CPU          *ResourceQuantityFilter `scheme:"cpu,omitempty"`
	Memory       *ResourceQuantityFilter `scheme:"memory,omitempty"`
	Pods         *ResourceQuantityFilter `scheme:"pods,omitempty"`
	Storage      *ResourceQuantityFilter `scheme:"storage,omitempty"`
}

// GetSelectorType returns the type of the Selector.
func (ks K8SliceSelector) GetSelectorType() FlavorTypeName {
	return K8SliceNameDefault
}
