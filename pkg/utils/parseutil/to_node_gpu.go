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

package parseutil

import (
	nodecorev1alpha1 "github.com/fluidos-project/node/apis/nodecore/v1alpha1"
	"github.com/fluidos-project/node/pkg/utils/models"
)

func ToNodeCoreGPU(in models.GpuCharacteristics) *nodecorev1alpha1.GPU {
	return &nodecorev1alpha1.GPU{
		Model:                 in.Model,
		Cores:                 in.Cores,
		Memory:                in.Memory,
		Vendor:                in.Vendor,
		Tier:                  in.Tier,
		Count:                 in.Count,
		MultiInstance:         in.MultiInstance,
		Shared:                in.Shared,
		SharingStrategy:       in.SharingStrategy,
		Dedicated:             in.Dedicated,
		Interruptible:         in.Interruptible,
		NetworkBandwidth:      in.NetworkBandwidth,
		NetworkLatencyMs:      in.NetworkLatencyMs,
		NetworkTier:           in.NetworkTier,
		TrainingScore:         in.TrainingScore,
		InferenceScore:        in.InferenceScore,
		HPCScore:              in.HPCScore,
		GraphicsScore:         in.GraphicsScore,
		Architecture:          in.Architecture,
		Interconnect:          in.Interconnect,
		InterconnectBandwidth: in.InterconnectBandwidth,
		CoresTotal:            in.Cores,
		ComputeCapability:     in.ComputeCapability,
		ClockSpeed:            in.ClockSpeed,
		FP32TFlops:            in.FP32TFlops,
		Topology:              in.Topology,
		MultiGPUEfficiency:    in.MultiGPUEfficiency,
		Region:                in.Region,
		Zone:                  in.Zone,
		HourlyRate:            in.HourlyRate,
		Provider:              in.Provider,
		PreEmptible:           in.PreEmptible,
	}

}
