/*
Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	api "github.com/gardener/machine-controller-manager-provider-alicloud/pkg/alicloud/apis"
	"github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/codes"
	"github.com/gardener/machine-controller-manager/pkg/util/provider/machinecodes/status"
	"strings"
)

const (
	// AlicloudAccessKeyID is a constant for a key name that is part of the Alibaba cloud credentials.
	AlicloudAccessKeyID string = "alicloudAccessKeyID"
	// AlicloudAccessKeySecret is a constant for a key name that is part of the Alibaba cloud credentials.
	AlicloudAccessKeySecret string = "alicloudAccessKeySecret"
	// AlicloudUserData is a constant for user data
	AlicloudUserData string = "userData"
	// alicloudDriverName is the name of the CSI driver for Alibaba Cloud
	AlicloudDriverName = "diskplugin.csi.alibabacloud.com"
)

func decodeProviderSpec(machineClass *v1alpha1.MachineClass) (*api.ProviderSpec, error) {
	var providerSpec *api.ProviderSpec
	err := json.Unmarshal(machineClass.ProviderSpec.Raw, &providerSpec)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return providerSpec, nil
}

func toInstanceTags(tags map[string]string) ([]ecs.RunInstancesTag, error) {
	result := []ecs.RunInstancesTag{{}, {}}
	hasCluster := false
	hasRole := false

	for k, v := range tags {
		if strings.Contains(k, "kubernetes.io/cluster/") {
			hasCluster = true
			result[0].Key = k
			result[0].Value = v
		} else if strings.Contains(k, "kubernetes.io/role/") {
			hasRole = true
			result[1].Key = k
			result[1].Value = v
		} else {
			result = append(result, ecs.RunInstancesTag{Key: k, Value: v})
		}
	}

	if !hasCluster || !hasRole {
		err := fmt.Errorf("Tags should at least contains 2 keys, which are prefixed with kubernetes.io/cluster and kubernetes.io/role")
		return nil, err
	}

	return result, nil
}

func encodeProviderID(region, instanceID string) string {
	return fmt.Sprintf("%s.%s", region, instanceID)
}

func decodeProviderID(providerID string) string {
	splitProviderID := strings.Split(providerID, ".")
	return splitProviderID[len(splitProviderID)-1]
}

// Host name in Alicloud has relationship with Instance ID
// i-uf69zddmom11ci7est12 => izuf69zddmom11ci7est12z
func instanceIDToName(instanceID string) string {
	return strings.Replace(instanceID, "-", "z", 1) + "z"
}