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

package outparamcheck

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/palantir/okgo/checker"
	"github.com/palantir/okgo/okgo"
)

const (
	TypeName okgo.CheckerType     = "outparamcheck"
	Priority okgo.CheckerPriority = 0
	MultiCPU okgo.CheckerMultiCPU = true
)

type Checker struct {
	OutParamFns map[string][]int
}

func (c *Checker) MultiCPU() (okgo.CheckerMultiCPU, error) {
	return MultiCPU, nil
}

func (c *Checker) Type() (okgo.CheckerType, error) {
	return TypeName, nil
}

func (c *Checker) Priority() (okgo.CheckerPriority, error) {
	return Priority, nil
}

var lineRegexp = regexp.MustCompile(`(.+):(\d+):(\d+)\t(.+)`)

func (c *Checker) Check(pkgPaths []string, pkgDir string, stdout io.Writer) {
	cfgJSON, err := json.Marshal(c.OutParamFns)
	if err != nil {
		okgo.WriteErrorAsIssue(err, stdout)
		return
	}

	cmd, wd := checker.AmalgomatedCheckCmd(string(TypeName), append([]string{
		"--config",
		string(cfgJSON),
	}, pkgPaths...), stdout)
	if cmd == nil {
		return
	}
	checker.RunCommandAndStreamOutput(cmd, func(line string) okgo.Issue {
		if strings.HasSuffix(line, `; the parameters listed above require the use of '&', for example f(&x) instead of f(x)`) {
			// skip the summary line
			return okgo.Issue{}
		}
		if match := lineRegexp.FindStringSubmatch(line); match != nil {
			// outparamcheck uses tab rather than space to separate prefix from content: transform to use space instead
			line = fmt.Sprintf("%s:%s:%s: %s", match[1], match[2], match[3], match[4])
		}
		return okgo.NewIssueFromLine(line, wd)
	}, stdout)
}

func (c *Checker) RunCheckCmd(args []string, stdout io.Writer) {
	checker.AmalgomatedRunRawCheck(string(TypeName), args, stdout)
}
