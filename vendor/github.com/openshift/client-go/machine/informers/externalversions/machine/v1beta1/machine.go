// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	context "context"
	time "time"

	apimachinev1beta1 "github.com/openshift/api/machine/v1beta1"
	versioned "github.com/openshift/client-go/machine/clientset/versioned"
	internalinterfaces "github.com/openshift/client-go/machine/informers/externalversions/internalinterfaces"
	machinev1beta1 "github.com/openshift/client-go/machine/listers/machine/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// MachineInformer provides access to a shared informer and lister for
// Machines.
type MachineInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() machinev1beta1.MachineLister
}

type machineInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewMachineInformer constructs a new informer for Machine type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewMachineInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredMachineInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredMachineInformer constructs a new informer for Machine type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredMachineInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MachineV1beta1().Machines(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MachineV1beta1().Machines(namespace).Watch(context.TODO(), options)
			},
		},
		&apimachinev1beta1.Machine{},
		resyncPeriod,
		indexers,
	)
}

func (f *machineInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredMachineInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *machineInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apimachinev1beta1.Machine{}, f.defaultInformer)
}

func (f *machineInformer) Lister() machinev1beta1.MachineLister {
	return machinev1beta1.NewMachineLister(f.Informer().GetIndexer())
}
