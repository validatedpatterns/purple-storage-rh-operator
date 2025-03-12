import {
    K8sGroupVersionKind,
    K8sResourceCommon,
} from '@openshift-console/dynamic-plugin-sdk';
// import { K8sModel } from '@openshift-console/dynamic-plugin-sdk/lib/api/common-types';

export const LocalVolumeDiscoveryResultKind: K8sGroupVersionKind = {
    version: 'v1alpha1',
    group: 'purple.purplestorage.com',
    kind: 'LocalVolumeDiscoveryResult',
};

export enum DeviceType {
    disk = "disk",
    mpath = "mpath"
}

export type Device = {
    size: number // 53687091200
    path: string // /dev/sda
    fstype: string // ''
    vendor: string // 'QEMU    '
    model: string // QEMU HARDDISK
    WWN: string // '0x5000c50015ea7599'
    deviceID: string // /dev/disk/by-id/wwn-0x5000c50015ea7599
    status: {
        state: string //Available
    }
    serial: string // seconddisk
    property: string //Rotational
    type: string //disk
}

export type LocalVolumeDiscoveryResultSpec = {
    spec?: {
        nodeName: string,
    }
    status?: {
        discoveredDevices?: Device[]
    };
} & K8sResourceCommon;
