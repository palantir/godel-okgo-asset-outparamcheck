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

package integration_test

import (
	"testing"

	"github.com/nmiyake/pkg/gofiles"
	"github.com/palantir/godel/framework/pluginapitester"
	"github.com/palantir/godel/pkg/products"
	"github.com/palantir/okgo/okgotester"
	"github.com/stretchr/testify/require"
)

const (
	okgoPluginLocator  = "com.palantir.okgo:check-plugin:1.0.0"
	okgoPluginResolver = "https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"
)

func TestCheck(t *testing.T) {
	const godelYML = `exclude:
  names:
    - "\\..+"
    - "vendor"
  paths:
    - "godel"
`

	assetPath, err := products.Bin("outparamcheck-asset")
	require.NoError(t, err)

	configFiles := map[string]string{
		"godel/config/godel.yml": godelYML,
		"godel/config/check-plugin.yml": `
checks:
  outparamcheck:
    config:
      out-param-funcs:
        "github.com/org/repo/foo.LoadConfig": [1]
`,
	}

	pluginProvider, err := pluginapitester.NewPluginProviderFromLocator(okgoPluginLocator, okgoPluginResolver)
	require.NoError(t, err)

	okgotester.RunAssetCheckTest(t,
		pluginProvider,
		pluginapitester.NewAssetProvider(assetPath),
		"outparamcheck",
		".",
		[]okgotester.AssetTestCase{
			{
				Name: "output parameter not used properly",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo.go",
						Src: `package foo

import "encoding/json"

func Foo() {
	var out string
	_ = json.Unmarshal(nil, out)
}
`,
					},
				},
				ConfigFiles: configFiles,
				WantError:   true,
				WantOutput: `Running outparamcheck...
foo.go:7:26: _ = json.Unmarshal(nil, out)  // 2nd argument of 'Unmarshal' requires '&'
Finished outparamcheck
Check(s) produced output: [outparamcheck]
`,
			},
			{
				Name: "output parameter specified by function not used properly",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo.go",
						Src: `package foo

import "github.com/org/repo/foo"

func Foo() {
	var out string
	foo.LoadConfig(nil, out)
}
`,
					},
					{
						RelPath: "vendor/github.com/org/repo/foo/foo.go",
						Src: `package foo

func LoadConfig(in []byte, out interface{}) {}
`,
					},
				},
				ConfigFiles: configFiles,
				WantError:   true,
				WantOutput: `Running outparamcheck...
foo.go:7:22: foo.LoadConfig(nil, out)  // 2nd argument of 'LoadConfig' requires '&'
Finished outparamcheck
Check(s) produced output: [outparamcheck]
`,
			},
			{
				Name: "output parameter not used properly in file from inner directory",
				Specs: []gofiles.GoFileSpec{
					{
						RelPath: "foo.go",
						Src: `package foo

import "encoding/json"

func Foo() {
	var out string
	_ = json.Unmarshal(nil, out)
}
`,
					},
					{
						RelPath: "inner/bar",
					},
				},
				ConfigFiles: configFiles,
				Wd:          "inner",
				WantError:   true,
				WantOutput: `Running outparamcheck...
../foo.go:7:26: _ = json.Unmarshal(nil, out)  // 2nd argument of 'Unmarshal' requires '&'
Finished outparamcheck
Check(s) produced output: [outparamcheck]
`,
			},
		},
	)
}

func TestUpgradeConfig(t *testing.T) {
	pluginProvider, err := pluginapitester.NewPluginProviderFromLocator(okgoPluginLocator, okgoPluginResolver)
	require.NoError(t, err)

	assetPath, err := products.Bin("outparamcheck-asset")
	require.NoError(t, err)
	assetProvider := pluginapitester.NewAssetProvider(assetPath)

	pluginapitester.RunUpgradeConfigTest(t,
		pluginProvider,
		[]pluginapitester.AssetProvider{assetProvider},
		[]pluginapitester.UpgradeConfigTestCase{
			{
				Name: `legacy configuration with empty "args" field is updated`,
				ConfigFiles: map[string]string{
					"godel/config/check.yml": `
checks:
  outparamcheck:
    filters:
      - value: "should have comment or be unexported"
      - type: name
        value: ".*.pb.go"
`,
				},
				Legacy: true,
				WantOutput: `Upgraded configuration for check-plugin.yml
`,
				WantFiles: map[string]string{
					"godel/config/check-plugin.yml": `checks:
  outparamcheck:
    filters:
    - value: should have comment or be unexported
    exclude:
      names:
      - .*.pb.go
`,
				},
			},
			{
				Name: `legacy configuration with "config" args is upgraded`,
				ConfigFiles: map[string]string{
					"godel/config/check.yml": `
checks:
  outparamcheck:
    args:
      - "-config"
      - |
        {
          "github.com/palantir/go-palantir/configloading.LoadConfig": [1],
          "github.com/palantir/go-palantir/safejson.Decode": [1]
        }
`,
				},
				Legacy: true,
				WantOutput: `Upgraded configuration for check-plugin.yml
`,
				WantFiles: map[string]string{
					"godel/config/check-plugin.yml": `checks:
  outparamcheck:
    config:
      out-param-funcs:
        github.com/palantir/go-palantir/configloading.LoadConfig:
        - 1
        github.com/palantir/go-palantir/safejson.Decode:
        - 1
`,
				},
			},
			{
				Name: `legacy configuration with args other than "config" fails`,
				ConfigFiles: map[string]string{
					"godel/config/check.yml": `
checks:
  outparamcheck:
    args:
      - "-help"
`,
				},
				Legacy:    true,
				WantError: true,
				WantOutput: `Failed to upgrade configuration:
	godel/config/check-plugin.yml: failed to upgrade configuration: failed to upgrade check "outparamcheck" legacy configuration: failed to upgrade asset configuration: outparamcheck-asset only supports legacy configuration if the first element in "args" is "-config"
`,
				WantFiles: map[string]string{
					"godel/config/check.yml": `
checks:
  outparamcheck:
    args:
      - "-help"
`,
				},
			},
			{
				Name: `legacy configuration with args that starts with "-config" but has more than 2 arguments fails"`,
				ConfigFiles: map[string]string{
					"godel/config/check.yml": `
checks:
  outparamcheck:
    args:
      - "-config"
      - a
      - b
`,
				},
				Legacy:    true,
				WantError: true,
				WantOutput: `Failed to upgrade configuration:
	godel/config/check-plugin.yml: failed to upgrade configuration: failed to upgrade check "outparamcheck" legacy configuration: failed to upgrade asset configuration: outparamcheck-asset only supports legacy configuration if "args" has exactly one element after "-config"
`,
				WantFiles: map[string]string{
					"godel/config/check.yml": `
checks:
  outparamcheck:
    args:
      - "-config"
      - a
      - b
`,
				},
			},
			{
				Name: `legacy configuration with "-config" argument that is not a valid JSON map fails`,
				ConfigFiles: map[string]string{
					"godel/config/check.yml": `
checks:
  outparamcheck:
    args:
      - "-config"
      - |
        {"foo":"bar",}
`,
				},
				Legacy:    true,
				WantError: true,
				WantOutput: `Failed to upgrade configuration:
	godel/config/check-plugin.yml: failed to upgrade configuration: failed to upgrade check "outparamcheck" legacy configuration: failed to upgrade asset configuration: failed to unmarshal second element of "args" in outparamcheck-asset legacy configuration as JSON map: invalid character '}' looking for beginning of object key string
`,
				WantFiles: map[string]string{
					"godel/config/check.yml": `
checks:
  outparamcheck:
    args:
      - "-config"
      - |
        {"foo":"bar",}
`,
				},
			},
			{
				Name: `valid v0 config works`,
				ConfigFiles: map[string]string{
					"godel/config/check-plugin.yml": `
checks:
  outparamcheck:
    config:
      # comment
      out-param-funcs:
        "github.com/palantir/go-palantir/configloading.LoadConfig": [1]
        "github.com/palantir/go-palantir/safejson.Decode": [1]
`,
				},
				WantOutput: ``,
				WantFiles: map[string]string{
					"godel/config/check-plugin.yml": `
checks:
  outparamcheck:
    config:
      # comment
      out-param-funcs:
        "github.com/palantir/go-palantir/configloading.LoadConfig": [1]
        "github.com/palantir/go-palantir/safejson.Decode": [1]
`,
				},
			},
		},
	)
}
