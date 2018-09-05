// Copyright © 2017 Aqua Security Software Ltd. <info@aquasec.com>
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

package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/aquasecurity/bench-common/check"
)

var (
	in       []byte
	controls *check.Controls
)

func init() {
	var err error
	in, err = ioutil.ReadFile("data")
	if err != nil {
		panic("Failed reading test data: " + err.Error())
	}

	// substitute variables in data file
	user := os.Getenv("USER")
	s := strings.Replace(string(in), "$user", user, -1)

	controls, err = check.NewControls([]byte(s))
	if err != nil {
		panic("Failed creating test controls: " + err.Error())
	}
}

func TestTestExecute(t *testing.T) {
	cases := []struct {
		*check.Check
		cmdStr string
		expect bool
	}{
		{
			controls.Groups[0].Checks[0],
			"2:45 ../kubernetes/kube-apiserver --allow-privileged=false --option1=20,30,40",
			true,
		},
		{
			controls.Groups[0].Checks[1],
			"2:45 ../kubernetes/kube-apiserver --allow-privileged=false",
			true,
		},
		{
			controls.Groups[0].Checks[2],
			"niinai   13617  2635 99 19:26 pts/20   00:03:08 ./kube-apiserver --insecure-port=0 --anonymous-auth",
			true,
		},
		{
			controls.Groups[0].Checks[3],
			"2:45 ../kubernetes/kube-apiserver --secure-port=0 --audit-log-maxage=40 --option",
			true,
		},
		{
			controls.Groups[0].Checks[4],
			"2:45 ../kubernetes/kube-apiserver --max-backlog=20 --secure-port=0 --audit-log-maxage=40 --option",
			true,
		},
		{
			controls.Groups[0].Checks[5],
			"2:45 ../kubernetes/kube-apiserver --option --admission-control=WebHook,RBAC ---audit-log-maxage=40",
			true,
		},
		{
			controls.Groups[0].Checks[6],
			"2:45 .. --kubelet-clientkey=foo --kubelet-client-certificate=bar --admission-control=Webhook,RBAC",
			true,
		},
		{
			controls.Groups[0].Checks[7],
			"2:45 ..  --secure-port=0 --kubelet-client-certificate=bar --admission-control=Webhook,RBAC",
			true,
		},
		{
			controls.Groups[0].Checks[8],
			"644",
			true,
		},
		{
			controls.Groups[0].Checks[9],
			"640",
			true,
		},
		{
			controls.Groups[0].Checks[9],
			"600",
			true,
		},
		{
			controls.Groups[0].Checks[10],
			"2:45 ../kubernetes/kube-apiserver --option --admission-control=WebHook,RBAC ---audit-log-maxage=40",
			true,
		},
		{
			controls.Groups[0].Checks[11],
			"2:45 ../kubernetes/kube-apiserver --option --admission-control=WebHook,RBAC ---audit-log-maxage=40",
			true,
		},
		{
			controls.Groups[0].Checks[12],
			"2:45 ../kubernetes/kube-apiserver --option --admission-control=WebHook,Something,RBAC ---audit-log-maxage=40",
			true,
		},
		{
			controls.Groups[0].Checks[13],
			"2:45 ../kubernetes/kube-apiserver --option --admission-control=Something ---audit-log-maxage=40",
			true,
		},
		{
			controls.Groups[0].Checks[14],
			"2:45 ../kubernetes/kube-apiserver --option --admission-control=Something ---audit-log-maxage=40",
			false,
		},
		{
			controls.Groups[0].Checks[15],
			"2:45 ../kubernetes/kube-apiserver --option --admission-control-config=Something ---audit-log-maxage=40",
			false,
		},
	}

	for _, c := range cases {
		res := c.Tests.Execute(c.cmdStr)
		if res.TestResult != c.expect {
			t.Errorf("id:%s %s, expected:%v, got:%v\n", c.ID, c.Description, c.expect, res.TestResult)
		}
	}
}
