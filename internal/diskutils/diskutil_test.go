//nolint:lll
package diskutils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var lsblkOut string
var blkidOut string

const (
	lsblkOutput1 = `{"blockdevices": [
				{"name": "sda", "rota": true, "type": "disk", "size": 62914560000, "model": "VBOX HARDDISK", "vendor": "ATA", "ro": false, "rm": false, "state": "running", "kname": "sda", "serial": "", "partlabel": "", "wwn": "0x500a07512b9f5254"},
				{"name": "sda1", "rota": true, "type": "part", "size": 62913494528, "model": "", "vendor": "", "ro": false, "rm": false, "state": "running", "kname": "sda1", "serial": "", "partlabel": "BIOS-BOOT", "wwn": "0x55cd2e41563851e9"}]}`
	lsblkOutput2 = `{"blockdevices": [
				{"name": "sdc","rota": true, "type": "disk", "size": 62914560000, "model": "VBOX HARDDISK", "vendor": "ATA", "ro": false, "rm": true, "state": "", "kname": "sdc", "serial": "", "partlabel": null, "wwn": "0x500a07512b9f5254"},
				{"name": "sdc3", "rota": true, "type": "part", "size": 62913494528, "model": "", "vendor": "", "ro": false, "rm": true, "state": "", "kname": "sdc3", "serial": "", "partlabel": null, "wwn": "0x55cd2e41563851e9"}]}`

	blkIDOutput1 = `/dev/sdc: TYPE="ext4"
/dev/sdc3: TYPE="ext2"
`
)

type mockCmdExec struct {
	stdout []string
	count  int
}

func (m *mockCmdExec) Execute(name string, args ...string) Command {
	return m
}

func (m *mockCmdExec) CombinedOutput() ([]byte, error) {
	o := m.stdout[m.count]
	m.count++
	return []byte(o), nil
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	defer os.Exit(0)
	switch os.Getenv("COMMAND") {
	case "lsblk":
		fmt.Fprint(os.Stdout, os.Getenv("LSBLKOUT"))
	case "blkid":
		fmt.Fprint(os.Stdout, os.Getenv("BLKIDOUT"))
	}
}

func TestListBlockDevices(t *testing.T) {
	testcases := []struct {
		label             string
		lsblkOutput       string
		blkIDOutput       string
		totalBlockDevices int
		totalBadRows      int
		expected          []BlockDevice
	}{
		{
			label:             "Case 1: block devices with no filesystems",
			lsblkOutput:       lsblkOutput1,
			blkIDOutput:       "",
			totalBlockDevices: 2,
			totalBadRows:      0,
			expected: []BlockDevice{
				{
					Name:       "sda",
					FSType:     "",
					Type:       "disk",
					Size:       62914560000,
					Model:      "VBOX HARDDISK",
					Vendor:     "ATA",
					Serial:     "",
					Rotational: true,
					ReadOnly:   false,
					Removable:  false,
					State:      "running",
					PartLabel:  "",
				},
				{

					Name:       "sda1",
					FSType:     "",
					Type:       "part",
					Size:       62913494528,
					Model:      "",
					Vendor:     "",
					Serial:     "",
					Rotational: true,
					ReadOnly:   false,
					Removable:  false,
					State:      "running",
					PartLabel:  "BIOS-BOOT",
				},
			},
		},
		{
			label:             "Case 2: block devices with filesystems",
			lsblkOutput:       lsblkOutput2,
			blkIDOutput:       blkIDOutput1,
			totalBlockDevices: 2,
			totalBadRows:      0,
			expected: []BlockDevice{
				{
					Name:       "sdc",
					FSType:     "ext4",
					Type:       "disk",
					Size:       62914560000,
					Model:      "VBOX HARDDISK",
					Vendor:     "ATA",
					Serial:     "",
					Rotational: true,
					ReadOnly:   false,
					Removable:  true,
					State:      "running",
					PartLabel:  "",
				},
				{

					Name:       "sdc3",
					FSType:     "ext2",
					Type:       "part",
					Size:       62913494528,
					Model:      "",
					Vendor:     "",
					Serial:     "",
					Rotational: true,
					ReadOnly:   false,
					Removable:  true,
					State:      "running",
					PartLabel:  "",
				},
			},
		},
		{
			label:             "Case 3: empty lsblk output",
			lsblkOutput:       `{"blockdevices": []}`,
			totalBlockDevices: 0,
			totalBadRows:      0,
			expected:          []BlockDevice{},
		},
		{
			label:             "Case 4: lsblk output with white space",
			lsblkOutput:       `{"blockdevices": [{"name":"sda","model":"VBOX HARDDISK   ","vendor":"ATA   "}]}`,
			totalBlockDevices: 1,
			totalBadRows:      0,
			expected: []BlockDevice{
				{
					Name:   "sda",
					Model:  "VBOX HARDDISK",
					Vendor: "ATA",
				},
			},
		},
	}

	for _, tc := range testcases {
		lsblkOut = tc.lsblkOutput
		blkidOut = tc.blkIDOutput
		ExecCommand = &mockCmdExec{stdout: []string{blkidOut, lsblkOut}}
		blockDevices, badRows, err := ListBlockDevices([]string{})
		assert.NoError(t, err, "[%q: Device]: invalid json", tc.label)
		assert.Equalf(t, tc.totalBadRows, len(badRows), "[%s] total bad rows list didn't match", tc.label)
		assert.Equalf(t, tc.totalBlockDevices, len(blockDevices), "[%s] total block device list didn't match", tc.label)
		for i := 0; i < len(blockDevices); i++ {
			assert.Equalf(t, tc.expected[i].Name, blockDevices[i].Name, "[%q: Device: %d]: invalid block device name", tc.label, i+1)
			assert.Equalf(t, tc.expected[i].Type, blockDevices[i].Type, "[%q: Device: %d]: invalid block device type", tc.label, i+1)
			assert.Equalf(t, tc.expected[i].FSType, blockDevices[i].FSType, "[%q: Device: %d]: invalid block device file system", tc.label, i+1)
			assert.Equalf(t, tc.expected[i].Size, blockDevices[i].Size, "[%q: Device: %d]: invalid block device size", tc.label, i+1)
			assert.Equalf(t, tc.expected[i].Vendor, blockDevices[i].Vendor, "[%q: Device: %d]: invalid block device vendor", tc.label, i+1)
			assert.Equalf(t, tc.expected[i].Model, blockDevices[i].Model, "[%q: Device: %d]: invalid block device Model", tc.label, i+1)
			assert.Equalf(t, tc.expected[i].Serial, blockDevices[i].Serial, "[%q: Device: %d]: invalid block device serial", tc.label, i+1)
			assert.Equalf(t, tc.expected[i].Rotational, blockDevices[i].Rotational, "[%q: Device: %d]: invalid block device rotational property", tc.label, i+1)
			assert.Equalf(t, tc.expected[i].ReadOnly, blockDevices[i].ReadOnly, "[%q: Device: %d]: invalid block device read only value", tc.label, i+1)
			assert.Equalf(t, tc.expected[i].PartLabel, blockDevices[i].PartLabel, "[%q: Device: %d]: invalid block device PartLabel value", tc.label, i+1)
		}
	}
}

func TestGetPathByID(t *testing.T) {
	testcases := []struct {
		label               string
		blockDevice         BlockDevice
		existingDeviceId    string
		fakeGlobfunc        func(string) ([]string, error)
		fakeEvalSymlinkfunc func(string) (string, error)
		expected            string
	}{
		{
			label:       "Case 1: pathByID is already available",
			blockDevice: BlockDevice{Name: "sdb", KName: "sdb", PathByID: "/dev/disk/by-id/sdb"},
			fakeGlobfunc: func(path string) ([]string, error) {
				return []string{"/dev/disk/by-id/dm-home", "/dev/disk/by-id/dm-uuid-LVM-6p00g8KptCD", "/dev/disk/by-id/sdb"}, nil
			},
			fakeEvalSymlinkfunc: func(path string) (string, error) {
				return "/dev/disk/by-id/sdb", nil
			},
			expected: "/dev/disk/by-id/sdb",
		},

		{
			label:       "Case 2: pathByID is not available",
			blockDevice: BlockDevice{Name: "sdb", KName: "sdb", PathByID: ""},
			fakeGlobfunc: func(path string) ([]string, error) {
				return []string{"/dev/disk/by-id/sdb"}, nil
			},
			fakeEvalSymlinkfunc: func(path string) (string, error) {
				return "/dev/disk/by-id/sdb", nil
			},
			expected: "/dev/disk/by-id/sdb",
		},
		{
			label:       "Prefer wwn-paths if available",
			blockDevice: BlockDevice{Name: "sdb", KName: "sdb", PathByID: ""},
			fakeGlobfunc: func(path string) ([]string, error) {
				return []string{"/dev/disk/by-id/abcde", "/dev/disk/by-id/wwn-abcde"}, nil
			},
			fakeEvalSymlinkfunc: func(string) (string, error) {
				return "/dev/sdb", nil
			},
			expected: "/dev/disk/by-id/wwn-abcde",
		},
		{
			label:       "Prefer wwn path even if scsi is if available",
			blockDevice: BlockDevice{Name: "sdb", KName: "sdb", PathByID: ""},
			fakeGlobfunc: func(path string) ([]string, error) {
				return []string{"/dev/disk/by-id/abcde", "/dev/disk/by-id/wwn-abcde", "/dev/disk/by-id/scsi-abcde"}, nil
			},
			fakeEvalSymlinkfunc: func(string) (string, error) {
				return "/dev/sdb", nil
			},
			expected: "/dev/disk/by-id/wwn-abcde",
		},
		{
			label:            "Prefer supplied path over anything else",
			blockDevice:      BlockDevice{Name: "sdb", KName: "sdb", PathByID: ""},
			existingDeviceId: "scsi-abcde",
			fakeGlobfunc: func(path string) ([]string, error) {
				return []string{"/dev/disk/by-id/abcde", "/dev/disk/by-id/wwn-abcde", "/dev/disk/by-id/scsi-abcde"}, nil
			},
			fakeEvalSymlinkfunc: func(string) (string, error) {
				return "/dev/sdb", nil
			},
			expected: "/dev/disk/by-id/scsi-abcde",
		},
	}

	for _, tc := range testcases {
		FilePathGlob = tc.fakeGlobfunc
		FilePathEvalSymLinks = tc.fakeEvalSymlinkfunc
		defer func() {
			FilePathGlob = filepath.Glob
			FilePathEvalSymLinks = filepath.EvalSymlinks
		}()

		actual, err := tc.blockDevice.GetPathByID(tc.existingDeviceId)
		assert.NoError(t, err)
		assert.Equalf(t, tc.expected, actual, "[%s] failed to get device path by ID", tc.label)
	}
}

func TestGetPathByIDFail(t *testing.T) {
	testcases := []struct {
		label               string
		blockDevice         BlockDevice
		fakeGlobfunc        func(string) ([]string, error)
		fakeEvalSymlinkfunc func(string) (string, error)
		expected            string
	}{
		{
			label:       "Case 1: filepath.Glob command failure",
			blockDevice: BlockDevice{KName: "sdb"},
			fakeGlobfunc: func(name string) ([]string, error) {
				return []string{}, fmt.Errorf("failed to list matching files")
			},
			fakeEvalSymlinkfunc: func(path string) (string, error) {
				return "/dev/disk/by-id/sdb", nil
			},
			expected: "",
		},

		{
			label:       "Case 2: filepath.EvalSymlinks command failure",
			blockDevice: BlockDevice{KName: "sdb", PathByID: ""},
			fakeGlobfunc: func(name string) ([]string, error) {
				return []string{"/dev/disk/by-id/sdb"}, nil
			},
			fakeEvalSymlinkfunc: func(path string) (string, error) {
				return "", fmt.Errorf("failed to evaluate symlink")
			},
			expected: "",
		},
	}

	for _, tc := range testcases {
		FilePathGlob = tc.fakeGlobfunc
		FilePathEvalSymLinks = tc.fakeEvalSymlinkfunc
		defer func() {
			FilePathGlob = filepath.Glob
			FilePathEvalSymLinks = filepath.EvalSymlinks
		}()

		actual, err := tc.blockDevice.GetPathByID("" /*existing symlinkpath */)
		assert.Error(t, err)
		assert.Equalf(t, tc.expected, actual, "[%s] failed to get device path by ID", tc.label)
	}
}
