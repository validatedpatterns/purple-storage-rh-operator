package main

import (
	"runtime"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/validatedpatterns/purple-storage-rh-operator/internal/diskmaker/discovery"
	"github.com/validatedpatterns/purple-storage-rh-operator/internal/localmetrics"
	"k8s.io/klog/v2"
)

func startDeviceDiscovery(cmd *cobra.Command, args []string) error {
	printVersion()
	// start local server to emit custom metrics
	err := localmetrics.NewConfigBuilder().
		WithCollectors(localmetrics.LVDMetricsList).
		Build()
	if err != nil {
		return errors.Wrap(err, "failed to discover devices")
	}

	discoveryObj, err := discovery.NewDeviceDiscovery()
	if err != nil {
		return errors.Wrap(err, "failed to discover devices")
	}

	err = discoveryObj.Start()
	if err != nil {
		return errors.Wrap(err, "failed to discover devices")
	}

	return nil
}

func printVersion() {
	klog.Infof("Go Version: %s", runtime.Version())
	klog.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
}
