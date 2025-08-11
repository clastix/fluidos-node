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
	"k8s.io/apimachinery/pkg/api/resource"
)

// NodeInfo represents a node and its resources.
type NodeInfo struct {
	UID             string          `json:"uid"`
	Name            string          `json:"name"`
	Architecture    string          `json:"architecture"`
	OperatingSystem string          `json:"os"`
	ResourceMetrics ResourceMetrics `json:"resources"`
}

// GPUMetrics represents GPU metrics.
type GPUMetrics struct {
	Vendor      string            `json:"vendor,omitempty"`
	Model       string            `json:"model,omitempty"`
	CountTotal  int64             `json:"total_count"`
	MemoryTotal resource.Quantity `json:"totalMemory"`
	Tier        string            `json:"tier,omitempty"`

	CoresAvailable  resource.Quantity `json:"availableCores"`
	MemoryAvailable resource.Quantity `json:"availableMemory"`
	CountAvailable  int64             `json:"available_count"`

	GPUSharingMetrics  `json:",inline"`
	GPUNetworkMetrics  `json:",inline"`
	GPUScoreMetrics    `json:",inline"`
	GPUSpecMetrics     `json:",inline"`
	GPURentingMetrics  `json:",inline"`
	GPUProviderMetrics `json:",inline"`
}

type GPUSharingMetrics struct {
	MultiInstance   bool   `json:"multi_instance,omitempty"`
	Shared          bool   `json:"shared,omitempty"`
	SharingStrategy string `json:"sharing_strategy,omitempty"`
	Dedicated       bool   `json:"dedicated,omitempty"`
	Interruptible   bool   `json:"interruptible,omitempty"`
}

type GPUNetworkMetrics struct {
	NetworkBandwidth resource.Quantity `json:"network_bandwidth,omitempty"`
	NetworkLatencyMs int64             `json:"network_latency_ms,omitempty"`
	NetworkTier      string            `json:"network_tier,omitempty"`
}

type GPUScoreMetrics struct {
	TrainingScore  float64 `json:"training_score,omitempty"`
	InferenceScore float64 `json:"inference_score,omitempty"`
	HPCScore       float64 `json:"hpc_score,omitempty"`
	GraphicsScore  float64 `json:"graphics_score,omitempty"`
}

type GPUSpecMetrics struct {
	Architecture          string            `json:"architecture,omitempty"`
	Interconnect          string            `json:"interconnect,omitempty"`
	InterconnectBandwidth resource.Quantity `json:"interconnect_bandwidth,omitempty"`
	CoresTotal            resource.Quantity `json:"totalCores"`
	ComputeCapability     string            `json:"compute_capability,omitempty"`
	ClockSpeed            resource.Quantity `json:"clock_speed,omitempty"`
	FP32TFlops            float64           `json:"fp32_tflops,omitempty"`
	Topology              string            `json:"topology,omitempty"`
	MultiGPUEfficiency    string            `json:"multi_gpu_efficiency,omitempty"`
}

type GPURentingMetrics struct {
	Region     string  `json:"region,omitempty"`
	Zone       string  `json:"zone,omitempty"`
	HourlyRate float64 `json:"hourly_rate,omitempty"`
}

type GPUProviderMetrics struct {
	Provider    string `json:"provider,omitempty"`
	PreEmptible bool   `json:"pre_emptible,omitempty"`
}

// ResourceMetrics represents resources of a certain node.
type ResourceMetrics struct {
	CPUTotal         resource.Quantity `json:"totalCPU"`
	CPUAvailable     resource.Quantity `json:"availableCPU"`
	MemoryTotal      resource.Quantity `json:"totalMemory"`
	MemoryAvailable  resource.Quantity `json:"availableMemory"`
	PodsTotal        resource.Quantity `json:"totalPods"`
	PodsAvailable    resource.Quantity `json:"availablePods"`
	EphemeralStorage resource.Quantity `json:"ephemeralStorage"`
	GPU              GPUMetrics        `json:"gpu"`
}
