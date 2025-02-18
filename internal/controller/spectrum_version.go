package controller

import "strings"

// Taken from https://www.ibm.com/docs/en/scalecontainernative/5.2.2?topic=planning-software-requirements
type StorageScaleData struct {
	CSIVersion                string   `json:"csi_version"`
	Architecture              []string `json:"architecture"`
	RemoteStorageClusterLevel string   `json:"remote_storage_cluster_level"`
	FileSystemVersion         string   `json:"file_system_version"`
	OpenShiftLevels           []string `json:"openshift_levels"`
}

// The dict key is the IBM Storage Scale Container Native version
var storageScaleTable = map[string]StorageScaleData{
	"5.1.5.0": {"2.7.0", []string{"x86_64", "ppc64le", "s390x"}, "5.1.3.0+", "29.00", []string{"4.9", "4.10", "4.11"}},
	"5.1.6.0": {"2.8.0", []string{"x86_64", "ppc64le", "s390x"}, "5.1.3.0+", "30.00", []string{"4.9", "4.10", "4.11"}},
	"5.1.7.0": {"2.9.0", []string{"x86_64", "ppc64le", "s390x"}, "5.1.3.0+", "31.00", []string{"4.10", "4.11", "4.12"}},
	"5.1.9.1": {"2.10.0", []string{"x86_64", "ppc64le", "s390x"}, "5.1.3.0+", "33.00", []string{"4.12", "4.13", "4.14"}},
	"5.1.9.3": {"2.10.1", []string{"x86_64", "ppc64le", "s390x"}, "5.1.3.0+", "33.00", []string{"4.12", "4.13", "4.14"}},
	"5.1.9.4": {"2.10.2", []string{"x86_64", "ppc64le", "s390x"}, "5.1.3.0+", "33.00", []string{"4.12", "4.13", "4.14"}},
	"5.1.9.5": {"2.10.3", []string{"x86_64", "ppc64le", "s390x"}, "5.1.3.0+", "33.00", []string{"4.12", "4.13", "4.14"}},
	"5.1.9.6": {"2.10.4", []string{"x86_64", "ppc64le", "s390x"}, "5.1.3.0+", "33.00", []string{"4.12", "4.13", "4.14"}},
	"5.1.9.7": {"2.10.5", []string{"x86_64", "ppc64le", "s390x"}, "5.1.3.0+", "33.00", []string{"4.12", "4.13", "4.14"}},
	"5.2.0.0": {"2.11.0", []string{"x86_64", "ppc64le", "s390x"}, "5.1.3.0+", "34.00", []string{"4.13", "4.14", "4.15"}},
	"5.2.0.1": {"2.11.1", []string{"x86_64", "ppc64le", "s390x"}, "5.1.3.0+", "34.00", []string{"4.13", "4.14", "4.15"}},
	"5.2.1.0": {"2.12.0", []string{"x86_64", "ppc64le", "s390x"}, "5.1.9.0+", "35.00", []string{"4.14", "4.15", "4.16"}},
	"5.2.1.1": {"2.12.1", []string{"x86_64", "ppc64le", "s390x"}, "5.1.9.0+", "35.00", []string{"4.14", "4.15", "4.16"}},
	"5.2.2.0": {"2.13.0", []string{"x86_64", "ppc64le", "s390x"}, "5.1.9.0+", "36.00", []string{"4.15", "4.16", "4.17"}},
}

func IsOpenShiftSupported(ibmStorageScaleVersion string, openShiftVersion string) bool {
	data, exists := storageScaleTable[ibmStorageScaleVersion]
	if !exists {
		return false
	}

	for _, version := range data.OpenShiftLevels {
		if strings.TrimSpace(version) == openShiftVersion {
			return true
		}
	}

	return false
}
