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
	"context"
	"encoding/json"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	nodecorev1alpha1 "github.com/fluidos-project/node/apis/nodecore/v1alpha1"
	"github.com/fluidos-project/node/pkg/indexer"
	"github.com/fluidos-project/node/pkg/utils/flags"
	"github.com/fluidos-project/node/pkg/utils/getters"
	"github.com/fluidos-project/node/pkg/utils/models"
	"github.com/fluidos-project/node/pkg/utils/namings"
	"github.com/fluidos-project/node/pkg/utils/parseutil"
)

// ClusterRole
// +kubebuilder:rbac:groups=nodecore.fluidos.eu,resources=flavors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=endpoints,verbs=get;list;watch
// +kubebuilder:rbac:groups=metrics.k8s.io,resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups=metrics.k8s.io,resources=nodes,verbs=get;list;watch

// NodeReconciler reconciles a Node object and creates Flavor objects.
type NodeReconciler struct {
	client.Client
	Scheme              *runtime.Scheme
	EnableAutoDiscovery bool
	WebhookServer       webhook.Server
	FlavorIndexer       indexer.FlavorByNodeName
}

func (r *NodeReconciler) LabelSelector() labels.Selector {
	return labels.Set{flags.ResourceNodeLabel: "true"}.AsSelector()
}

// Reconcile reconciles a Node object to create Flavor objects.
func (r *NodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx, "node", req.NamespacedName)
	ctx = ctrl.LoggerInto(ctx, log)
	// Check if AutoDiscovery is enabled
	if !r.EnableAutoDiscovery {
		log.Info("AutoDiscovery is disabled")
		return ctrl.Result{}, nil
	}

	if r.WebhookServer != nil {
		// Check if the webhook server is running
		if err := r.WebhookServer.StartedChecker()(nil); err != nil {
			log.Info("Webhook server not started yet, requeuing the request")
			return ctrl.Result{Requeue: true}, nil
		}
	}
	// Fetch the Node instance
	var node corev1.Node
	if err := r.Get(ctx, req.NamespacedName, &node); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Info("Node not found")
			return ctrl.Result{}, nil
		}
	}
	// Check if the node has the label
	if !r.LabelSelector().Matches(labels.Set(node.GetLabels())) {
		log.Info("Node %s does not have the label %s", node.Name, flags.ResourceNodeLabel)

		return ctrl.Result{}, nil
	}

	var nodeMetrics metricsv1beta1.NodeMetrics
	// Get the node metrics referred to the node
	if err := r.Client.Get(ctx, client.ObjectKey{Name: node.Name}, &nodeMetrics); err != nil {
		log.Error(err, "error getting NodeMetrics", err)
		return ctrl.Result{}, err
	}
	// Get the NodeInfo struct for the node and its metrics
	nodeInfo, err := GetNodeInfos(&node, &nodeMetrics)
	if err != nil {
		log.Error(err, "error getting NodeInfo", err)
		return ctrl.Result{}, err
	}
	log.Info("NodeInfo created", "value", nodeInfo.Name)
	// Get NodeIdentity
	nodeIdentity := getters.GetNodeIdentity(ctx, r.Client)
	if nodeIdentity == nil {
		log.Info("error getting FLUIDOS Node identity")
		return ctrl.Result{}, nil
	}
	// Get all the Flavors owned by this node as kubernetes ownership:
	// iterating over a list is required since a Flavor can be owned by multiple resources.
	var matchFlavors nodecorev1alpha1.FlavorList
	if err = r.Client.List(ctx, &matchFlavors, client.MatchingFieldsSelector{Selector: fields.OneTermEqualSelector(r.FlavorIndexer.Field(), node.Name)}); err != nil {
		log.Error(err, "error listing Flavors")
		return ctrl.Result{}, nil
	}
	// Check if you have found any Flavor
	var flavor *nodecorev1alpha1.Flavor

	for _, i := range matchFlavors.Items {
		for _, or := range i.OwnerReferences {
			if or.Kind == "Node" {
				flavor = &i

				log.Info("Flavor found", "namespacedName", client.ObjectKeyFromObject(flavor), "node", node.Name)

				break
			}
		}
	}

	if err = r.createOrUpdateFlavor(ctx, flavor, nodeInfo, *nodeIdentity, &node); err != nil {
		log.Error(err, "error creating or updating Flavor", err)
		return ctrl.Result{Requeue: true}, nil
	}

	log.Info("Flavor reconciliation completed")

	return ctrl.Result{}, nil
}

func (r *NodeReconciler) createOrUpdateFlavor(ctx context.Context, flavor *nodecorev1alpha1.Flavor, nodeInfo *models.NodeInfo, nodeIdentity nodecorev1alpha1.NodeIdentity, owner client.Object) (err error) {
	log := ctrl.LoggerFrom(ctx)
	// Forge the Flavor from the NodeInfo and NodeIdentity
	shouldCreate := flavor == nil

	if shouldCreate {
		flavor = &nodecorev1alpha1.Flavor{}
		flavor.Name = namings.ForgeFlavorName(string(nodecorev1alpha1.TypeK8Slice), nodeIdentity.Domain)
	}
	flavor.Namespace = flags.FluidosNamespace

	log.Info("ready to handle Flavor", "namespacedName", client.ObjectKeyFromObject(flavor), "type", nodecorev1alpha1.TypeK8Slice)
	// Creating a new flavor custom resource from the metrics of the node.
	res, err := controllerutil.CreateOrUpdate(ctx, r.Client, flavor, func() error {
		var k8sSliceType nodecorev1alpha1.K8Slice
		if len(flavor.Spec.FlavorType.TypeData.Raw) > 0 {
			if unmarshalErr := json.Unmarshal(flavor.Spec.FlavorType.TypeData.Raw, &k8sSliceType); unmarshalErr != nil {
				return unmarshalErr
			}
		}

		if shouldCreate {
			k8sSliceType.Characteristics.CPU = nodeInfo.ResourceMetrics.CPUAvailable
			k8sSliceType.Characteristics.Memory = nodeInfo.ResourceMetrics.MemoryAvailable
			k8sSliceType.Characteristics.Pods = nodeInfo.ResourceMetrics.PodsAvailable
			k8sSliceType.Characteristics.Storage = &nodeInfo.ResourceMetrics.EphemeralStorage
		}

		k8sSliceType.Characteristics.Architecture = nodeInfo.Architecture
		k8sSliceType.Characteristics.Gpu = &nodecorev1alpha1.GPU{
			Model:                 nodeInfo.ResourceMetrics.GPU.Model,
			Memory:                nodeInfo.ResourceMetrics.GPU.MemoryTotal,
			Vendor:                nodeInfo.ResourceMetrics.GPU.Vendor,
			Tier:                  nodeInfo.ResourceMetrics.GPU.Tier,
			MultiInstance:         nodeInfo.ResourceMetrics.GPU.MultiInstance,
			Shared:                nodeInfo.ResourceMetrics.GPU.Shared,
			SharingStrategy:       nodeInfo.ResourceMetrics.GPU.SharingStrategy,
			Dedicated:             nodeInfo.ResourceMetrics.GPU.Dedicated,
			Interruptible:         nodeInfo.ResourceMetrics.GPU.Interruptible,
			NetworkBandwidth:      nodeInfo.ResourceMetrics.GPU.NetworkBandwidth,
			NetworkLatencyMs:      nodeInfo.ResourceMetrics.GPU.NetworkLatencyMs,
			NetworkTier:           nodeInfo.ResourceMetrics.GPU.NetworkTier,
			TrainingScore:         nodeInfo.ResourceMetrics.GPU.TrainingScore,
			InferenceScore:        nodeInfo.ResourceMetrics.GPU.InferenceScore,
			HPCScore:              nodeInfo.ResourceMetrics.GPU.HPCScore,
			GraphicsScore:         nodeInfo.ResourceMetrics.GPU.GraphicsScore,
			Architecture:          nodeInfo.ResourceMetrics.GPU.Architecture,
			Interconnect:          nodeInfo.ResourceMetrics.GPU.Interconnect,
			InterconnectBandwidth: nodeInfo.ResourceMetrics.GPU.InterconnectBandwidth,
			CoresTotal:            nodeInfo.ResourceMetrics.GPU.CoresTotal,
			ComputeCapability:     nodeInfo.ResourceMetrics.GPU.ComputeCapability,
			ClockSpeed:            nodeInfo.ResourceMetrics.GPU.ClockSpeed,
			FP32TFlops:            nodeInfo.ResourceMetrics.GPU.FP32TFlops,
			Topology:              nodeInfo.ResourceMetrics.GPU.Topology,
			MultiGPUEfficiency:    nodeInfo.ResourceMetrics.GPU.MultiGPUEfficiency,
			Region:                nodeInfo.ResourceMetrics.GPU.Region,
			Zone:                  nodeInfo.ResourceMetrics.GPU.Zone,
			HourlyRate:            nodeInfo.ResourceMetrics.GPU.HourlyRate,
			Provider:              nodeInfo.ResourceMetrics.GPU.Provider,
			PreEmptible:           nodeInfo.ResourceMetrics.GPU.PreEmptible,
		}
		k8sSliceType.Policies = nodecorev1alpha1.Policies{
			Partitionability: nodecorev1alpha1.Partitionability{
				CPUMin:     parseutil.ParseQuantityFromString(flags.CPUMin),
				MemoryMin:  parseutil.ParseQuantityFromString(flags.MemoryMin),
				PodsMin:    parseutil.ParseQuantityFromString(flags.PodsMin),
				CPUStep:    parseutil.ParseQuantityFromString(flags.CPUStep),
				MemoryStep: parseutil.ParseQuantityFromString(flags.MemoryStep),
				PodsStep:   parseutil.ParseQuantityFromString(flags.PodsStep),
			},
		}
		// Serialize K8SliceType to JSON
		k8SliceTypeJSON, marshalErr := json.Marshal(k8sSliceType)
		if marshalErr != nil {
			return marshalErr
		}

		flavor.Spec.ProviderID = nodeIdentity.NodeID
		flavor.Spec.FlavorType = nodecorev1alpha1.FlavorType{
			TypeIdentifier: nodecorev1alpha1.TypeK8Slice,
			TypeData:       runtime.RawExtension{Raw: k8SliceTypeJSON},
		}
		flavor.Spec.Owner = nodeIdentity
		// The following options can be changed at runtime by the provider,
		// avoiding to overwriting them at every reconciliation.
		if shouldCreate {
			flavor.Spec.Price.Amount = flags.AMOUNT
			flavor.Spec.Price.Currency = flags.CURRENCY
			flavor.Spec.Price.Period = flags.PERIOD
			flavor.Spec.Availability = true
			// FIXME: NetworkPropertyType should be taken in a smarter way
			flavor.Spec.NetworkPropertyType = "networkProperty"
			// FIXME: Location should be taken in a smarter way
			flavor.Spec.Location = &nodecorev1alpha1.Location{
				Latitude:        "10",
				Longitude:       "58",
				Country:         "Italy",
				City:            "Turin",
				AdditionalNotes: "None",
			}
		}

		return controllerutil.SetOwnerReference(owner, flavor, r.Client.Scheme())
	})
	if err != nil {
		return err
	}

	log.Info("Flavor handling completed", "namespacedName", client.ObjectKeyFromObject(flavor), "res", res)

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Node{}, builder.WithPredicates(predicate.NewPredicateFuncs(func(object client.Object) bool {
			return r.LabelSelector().Matches(labels.Set(object.GetLabels()))
		}))).
		Owns(&nodecorev1alpha1.Flavor{}, builder.MatchEveryOwner).
		Complete(r)
}
