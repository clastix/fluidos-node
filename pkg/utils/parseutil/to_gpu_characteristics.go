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
	"github.com/fluidos-project/node/apis/nodecore/v1alpha1"
	"github.com/fluidos-project/node/pkg/utils/models"
)

func ToGpuCharacteristics(input v1alpha1.GPU) *models.GpuCharacteristics {
	return &models.GpuCharacteristics{
		Vendor:                input.Vendor,
		Model:                 input.Model,
		Count:                 input.Count,
		Tier:                  input.Tier,
		MultiInstance:         input.MultiInstance,
		Shared:                input.Shared,
		SharingStrategy:       input.SharingStrategy,
		Dedicated:             input.Dedicated,
		Interruptible:         input.Interruptible,
		NetworkBandwidth:      input.NetworkBandwidth,
		NetworkLatencyMs:      input.NetworkLatencyMs,
		NetworkTier:           input.NetworkTier,
		TrainingScore:         input.TrainingScore,
		InferenceScore:        input.InferenceScore,
		HPCScore:              input.HPCScore,
		GraphicsScore:         input.GraphicsScore,
		Architecture:          input.Architecture,
		Interconnect:          input.Interconnect,
		InterconnectBandwidth: input.InterconnectBandwidth,
		Cores:                 input.CoresTotal,
		Memory:                input.Memory,
		ComputeCapability:     input.ComputeCapability,
		ClockSpeed:            input.ClockSpeed,
		FP32TFlops:            input.FP32TFlops,
		Topology:              input.Topology,
		MultiGPUEfficiency:    input.MultiGPUEfficiency,
		Region:                input.Region,
		Zone:                  input.Zone,
		HourlyRate:            input.HourlyRate,
		Provider:              input.Provider,
		PreEmptible:           input.PreEmptible,
	}
}
