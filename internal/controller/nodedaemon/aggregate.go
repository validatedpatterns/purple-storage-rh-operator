package nodedaemon

import (
	"context"
	"fmt"
	"sort"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	localv1alpha1 "github.com/darkdoc/purple-storage-rh-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *DaemonReconciler) aggregateDeamonInfo(ctx context.Context, request reconcile.Request) ([]corev1.Toleration, []metav1.OwnerReference, *corev1.NodeSelector, error) {
	// list
	lvSetList := localv1alpha1.PurpleStorageList{}
	err := r.Client.List(ctx, &lvSetList, client.InNamespace(request.Namespace))
	if err != nil {
		return []corev1.Toleration{}, []metav1.OwnerReference{}, nil, fmt.Errorf("could not fetch localvolumeset link: %w", err)
	}

	lvSets := lvSetList.Items
	tolerations, ownerRefs, terms := extractLVSetInfo(lvSets)

	var nodeSelector *corev1.NodeSelector = nil
	if len(terms) > 0 {
		nodeSelector = &corev1.NodeSelector{NodeSelectorTerms: terms}
	}

	return tolerations, ownerRefs, nodeSelector, err
}

func extractLVSetInfo(lvsets []localv1alpha1.PurpleStorage) ([]corev1.Toleration, []metav1.OwnerReference, []corev1.NodeSelectorTerm) {
	tolerations := make([]corev1.Toleration, 0)
	ownerRefs := make([]metav1.OwnerReference, 0)
	terms := make([]corev1.NodeSelectorTerm, 0)
	// if any one of the lvset nodeSelectors are nil, the terms should be empty to indicate matchAllNodes
	matchAllNodes := false

	// sort so that changing order doesn't cause unnecessary updates
	sort.SliceStable(lvsets, func(i, j int) bool {
		return lvsets[i].GetName() < lvsets[j].GetName()
	})
	for _, lvset := range lvsets {
		tolerations = append(tolerations, lvset.Spec.NodeSpec.Tolerations...)

		falseVar := false
		ownerRefs = append(ownerRefs, metav1.OwnerReference{
			UID:                lvset.GetUID(),
			Name:               lvset.GetName(),
			APIVersion:         lvset.APIVersion,
			Kind:               lvset.Kind,
			Controller:         &falseVar,
			BlockOwnerDeletion: &falseVar,
		})
		if lvset.Spec.NodeSpec.Selector != nil {
			terms = append(terms, lvset.Spec.NodeSpec.Selector.NodeSelectorTerms...)
		} else {
			matchAllNodes = true
		}
	}
	if matchAllNodes {
		terms = make([]corev1.NodeSelectorTerm, 0)
	}

	return tolerations, ownerRefs, terms
}
