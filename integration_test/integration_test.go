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
	"github.com/palantir/godel/pkg/products"
	"github.com/palantir/okgo/okgotester"
	"github.com/stretchr/testify/require"
)

const (
	okgoPluginLocator  = "com.palantir.okgo:okgo-plugin:0.3.0"
	okgoPluginResolver = "https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"

	godelYML = `exclude:
  names:
    - "\\..+"
    - "vendor"
  paths:
    - "godel"
`
)

func TestOutparamcheck(t *testing.T) {
	assetPath, err := products.Bin("outparamcheck-asset")
	require.NoError(t, err)

	configFiles := map[string]string{
		"godel/config/godel.yml": godelYML,
		"godel/config/check.yml": `
checks:
  outparamcheck:
    config:
      out-param-fns:
        "github.com/org/repo/foo.LoadConfig": [1]
`,
	}

	okgotester.RunAssetCheckTest(t,
		okgoPluginLocator, okgoPluginResolver,
		assetPath, "outparamcheck",
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
`,
			},
		},
	)
}
