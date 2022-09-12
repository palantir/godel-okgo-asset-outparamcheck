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

// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flag_test

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
)

func ExampleFunc() {
	fs := flag.NewFlagSet("ExampleFunc", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	var ip net.IP
	fs.Func("ip", "`IP address` to parse", func(s string) error {
		ip = net.ParseIP(s)
		if ip == nil {
			return errors.New("could not parse IP")
		}
		return nil
	})
	fs.Parse([]string{"-ip", "127.0.0.1"})
	fmt.Printf("{ip: %v, loopback: %t}\n\n", ip, ip.IsLoopback())

	// 256 is not a valid IPv4 component
	fs.Parse([]string{"-ip", "256.0.0.1"})
	fmt.Printf("{ip: %v, loopback: %t}\n\n", ip, ip.IsLoopback())

	// Output:
	// {ip: 127.0.0.1, loopback: true}
	//
	// invalid value "256.0.0.1" for flag -ip: could not parse IP
	// Usage of ExampleFunc:
	//   -ip IP address
	//     	IP address to parse
	// {ip: <nil>, loopback: false}
}
