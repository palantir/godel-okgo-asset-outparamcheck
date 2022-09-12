// Copyright 2022 Palantir Technologies, Inc.
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

// Copyright 2016 Palantir Technologies, Inc. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package amalgomated

import (
	"github.com/palantir/godel-okgo-asset-outparamcheck/internal/generated_src/internal/amalgomated_flag"
	"fmt"
	"os"
	"runtime"

	"github.com/palantir/godel-okgo-asset-outparamcheck/internal/generated_src/internal/github.com/palantir/outparamcheck/outparamcheck"
)

func AmalgomatedMain() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cfgPath := ""
	fset := flag.CommandLine
	fset.StringVar(&cfgPath, "config", "", "JSON configuration or '@' followed by path to a configuration file (@pathToJsonFile)")
	flag.Parse()

	err := outparamcheck.Run(cfgPath, flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
