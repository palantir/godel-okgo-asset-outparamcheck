// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package legacy

import (
	"encoding/json"

	v0 "github.com/palantir/godel-okgo-asset-outparamcheck/outparamcheck/config/internal/v0"
	"github.com/palantir/godel/v2/pkg/versionedconfig"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	versionedconfig.ConfigWithLegacy `yaml:",inline"`
	Args                             []string `yaml:"args"`
}

func UpgradeConfig(cfgBytes []byte) ([]byte, error) {
	var legacyCfg Config
	if err := yaml.UnmarshalStrict(cfgBytes, &legacyCfg); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal outparamcheck-asset legacy configuration")
	}
	if len(legacyCfg.Args) == 0 {
		return nil, nil
	}

	var upgradedCfg v0.Config
	if legacyCfg.Args[0] != "-config" {
		return nil, errors.Errorf(`outparamcheck-asset only supports legacy configuration if the first element in "args" is "-config"`)
	}

	if len(legacyCfg.Args) != 2 {
		return nil, errors.Errorf(`outparamcheck-asset only supports legacy configuration if "args" has exactly one element after "-config"`)
	}

	var jsonMapCfg map[string][]int
	if err := json.Unmarshal([]byte(legacyCfg.Args[1]), &jsonMapCfg); err != nil {
		return nil, errors.Wrapf(err, `failed to unmarshal second element of "args" in outparamcheck-asset legacy configuration as JSON map`)
	}

	if len(jsonMapCfg) > 0 {
		upgradedCfg.OutParamFuncs = make(map[string][]int)
	}
	for k, v := range jsonMapCfg {
		upgradedCfg.OutParamFuncs[k] = v
	}

	upgradedCfgBytes, err := yaml.Marshal(upgradedCfg)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal outparamcheck-asset v0 configuration")
	}
	return upgradedCfgBytes, nil
}
