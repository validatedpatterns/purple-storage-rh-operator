package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "diskmaker",
	Short: "Used to start the diskmaker daemon for the local-storage-operator",
}

var discoveryDaemonCmd = &cobra.Command{
	Use:   "discover",
	Short: "Used to start device discovery for the LocalVolumeDiscovery CR",
	RunE:  startDeviceDiscovery,
}

func main() {
	rootCmd.AddCommand(discoveryDaemonCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
