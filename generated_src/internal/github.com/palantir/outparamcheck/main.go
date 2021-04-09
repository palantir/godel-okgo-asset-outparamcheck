// Copyright 2016 Palantir Technologies, Inc. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package amalgomated

import (
	"amalgomated_flag"
	"fmt"
	"os"
	"runtime"

	"./outparamcheck"
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
