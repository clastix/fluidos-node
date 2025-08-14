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

package localresourcemanager

import (
	"fmt"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog/v2"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"

	"github.com/fluidos-project/node/pkg/utils/models"
)

// GetNodeInfos returns the NodeInfo struct for a given node and its metrics.
func GetNodeInfos(node *corev1.Node, nodeMetrics *metricsv1beta1.NodeMetrics) (*models.NodeInfo, error) {
	// Check if the node and the node metrics match
	if node.Name != nodeMetrics.Name {
		klog.Info("Node and NodeMetrics do not match")
		return nil, fmt.Errorf("node and node metrics do not match")
	}

	metricsStruct := forgeResourceMetrics(nodeMetrics, node)
	nodeInfo := forgeNodeInfo(node, metricsStruct)

	return nodeInfo, nil
}

// forgeResourceMetrics creates from params a new ResourceMetrics Struct.
func forgeResourceMetrics(nodeMetrics *metricsv1beta1.NodeMetrics, node *corev1.Node) *models.ResourceMetrics {
	// Get the total and used resources
	cpuTotal := node.Status.Allocatable.Cpu().DeepCopy()
	cpuUsed := nodeMetrics.Usage.Cpu().DeepCopy()
	memoryTotal := node.Status.Allocatable.Memory().DeepCopy()
	memoryUsed := nodeMetrics.Usage.Memory().DeepCopy()
	podsTotal := node.Status.Allocatable.Pods().DeepCopy()
	podsUsed := nodeMetrics.Usage.Pods().DeepCopy()
	ephemeralStorage := nodeMetrics.Usage.StorageEphemeral().DeepCopy()

	// Compute the available resources
	cpuAvail := cpuTotal.DeepCopy()
	memAvail := memoryTotal.DeepCopy()
	podsAvail := podsTotal.DeepCopy()
	cpuAvail.Sub(cpuUsed)
	memAvail.Sub(memoryUsed)
	podsAvail.Sub(podsUsed)

	var gpuMetrics models.GPUMetrics

	annotations := node.ObjectMeta.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}

	if v, found := annotations["provider.fluidos.eu/name"]; found {
		gpuMetrics.Provider = v
	}
	if v, found := annotations["gpu.fluidos.eu/vendor"]; found {
		gpuMetrics.Vendor = v
	}
	if v, found := annotations["gpu.fluidos.eu/model"]; found {
		gpuMetrics.Model = v
	}
	if v, found := annotations["gpu.fluidos.eu/count"]; found {
		count, _ := strconv.Atoi(v)
		gpuMetrics.Count = int64(count)
	}
	if v, found := annotations["gpu.fluidos.eu/memory-per-gpu"]; found {
		qty := resource.MustParse(v)

		computed := resource.NewQuantity(0, resource.BinarySI)
		for range gpuMetrics.Count {
			computed.Add(qty)
		}

		gpuMetrics.MemoryTotal = *computed
	}
	if v, found := annotations["gpu.fluidos.eu/tier"]; found {
		gpuMetrics.Tier = v
	}
	if v, found := annotations["gpu.fluidos.eu/architecture"]; found {
		gpuMetrics.Architecture = v
	}
	if v, found := annotations["gpu.fluidos.eu/compute-capability"]; found {
		gpuMetrics.ComputeCapability = v
	}
	if v, found := annotations["nvidia.fluidos.eu/mig-capable"]; found {
		gpuMetrics.MultiInstance, _ = strconv.ParseBool(v)
	}
	if v, found := annotations["gpu.fluidos.eu/fp32-tflops"]; found {
		fp32tFlops, _ := strconv.ParseFloat(v, 64)
		gpuMetrics.FP32TFlops = fp32tFlops
	}
	if v, found := annotations["gpu.fluidos.eu/sharing-capable"]; found {
		gpuMetrics.Shared, _ = strconv.ParseBool(v)
	}
	if v, found := annotations["gpu.fluidos.eu/sharing-strategy"]; found {
		gpuMetrics.SharingStrategy = v
	}
	if v, found := annotations["gpu.fluidos.eu/interconnect"]; found {
		gpuMetrics.Interconnect = v
	}
	if v, found := annotations["gpu.fluidos.eu/interconnect-bandwidth-gbps"]; found {
		qty := resource.MustParse(v)
		gpuMetrics.InterconnectBandwidth = qty
	}
	if v, found := annotations["gpu.fluidos.eu/cores"]; found {
		qty := resource.MustParse(v)

		computed := resource.NewQuantity(0, resource.BinarySI)
		for range gpuMetrics.Count {
			computed.Add(qty)
		}

		gpuMetrics.CoresTotal = *computed
	}
	if v, found := annotations["gpu.fluidos.eu/clock-speed"]; found {
		gpuMetrics.ClockSpeed = resource.MustParse(v)
	}
	if v, found := annotations["gpu.fluidos.eu/interruptible"]; found {
		gpuMetrics.Interruptible, _ = strconv.ParseBool(v)
	}
	if v, found := annotations["gpu.fluidos.eu/dedicated"]; found {
		gpuMetrics.Dedicated, _ = strconv.ParseBool(v)
	}
	if v, found := annotations["gpu.fluidos.eu/topology"]; found {
		gpuMetrics.Topology = v
	}
	if v, found := annotations["gpu.fluidos.eu/multi-gpu-efficiency"]; found {
		gpuMetrics.MultiGPUEfficiency = v
	}
	if v, found := annotations["cost.fluidos.eu/hourly-rate"]; found {
		hourlyRate, _ := strconv.ParseFloat(v, 64)
		gpuMetrics.HourlyRate = hourlyRate
	}
	if v, found := annotations["provider.fluidos.eu/preemptible"]; found {
		gpuMetrics.PreEmptible, _ = strconv.ParseBool(v)
	}
	if v, found := annotations["workload.fluidos.eu/training-score"]; found {
		score, _ := strconv.ParseFloat(v, 64)
		gpuMetrics.TrainingScore = score
	}
	if v, found := annotations["workload.fluidos.eu/inference-score"]; found {
		score, _ := strconv.ParseFloat(v, 64)
		gpuMetrics.InferenceScore = score
	}
	if v, found := annotations["workload.fluidos.eu/hpc-score"]; found {
		score, _ := strconv.ParseFloat(v, 64)
		gpuMetrics.HPCScore = score
	}
	if v, found := annotations["workload.fluidos.eu/graphics-score"]; found {
		score, _ := strconv.ParseFloat(v, 64)
		gpuMetrics.GraphicsScore = score
	}
	if v, found := annotations["network.fluidos.eu/bandwidth-gbps"]; found {
		qty := resource.MustParse(v)
		gpuMetrics.NetworkBandwidth = qty
	}
	if v, found := annotations["network.fluidos.eu/latency-ms"]; found {
		ms, _ := strconv.ParseInt(v, 10, 64)
		gpuMetrics.NetworkLatencyMs = ms
	}
	if v, found := annotations["network.fluidos.eu/tier"]; found {
		gpuMetrics.NetworkTier = v
	}

	if v, found := annotations["location.fluidos.eu/zone"]; found {
		gpuMetrics.Region = v
	}

	if v, found := annotations["location.fluidos.eu/region"]; found {
		gpuMetrics.Zone = v
	}

	return &models.ResourceMetrics{
		CPUTotal:         cpuTotal,
		CPUAvailable:     cpuAvail,
		MemoryTotal:      memoryTotal,
		MemoryAvailable:  memAvail,
		PodsTotal:        podsTotal,
		PodsAvailable:    podsAvail,
		EphemeralStorage: ephemeralStorage,
		GPU:              gpuMetrics,
	}
}

// forgeNodeInfo creates from params a new NodeInfo struct.
func forgeNodeInfo(node *corev1.Node, metrics *models.ResourceMetrics) *models.NodeInfo {
	return &models.NodeInfo{
		UID:             string(node.UID),
		Name:            node.Name,
		Architecture:    node.Status.NodeInfo.Architecture,
		OperatingSystem: node.Status.NodeInfo.OperatingSystem,
		ResourceMetrics: *metrics,
	}
}
