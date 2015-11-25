//
// +build integration

/*
http://www.apache.org/licenses/LICENSE-2.0.txt

Copyright 2015 Intel Corporation

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

package riemann

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/amir/raidman"
)

// integration test
func TestRiemannPublish(t *testing.T) {
	// This integration test requires a Riemann Server
	broker := os.Getenv("SNAP_TEST_RIEMANN")
	if broker == "" {
		fmt.Println("Skipping integration tests")
		return
	}

	var buf bytes.Buffer
	buf.Reset()
	r := NewRiemannPublisher()
	config := make(map[string]ctypes.ConfigValue)
	config["broker"] = ctypes.ConfigValueStr{Value: broker}
	cp, _ := r.GetConfigPolicy()
	cfg, _ := cp.Get([]string{""}).Process(config)
	metrics := []plugin.PluginMetricType{
		*plugin.NewPluginMetricType([]string{"intel", "cpu", "temp"}, time.Now(), "bacon-powered", nil, nil, 100),
	}
	enc := gob.NewEncoder(&buf)
	enc.Encode(metrics)
	err := r.Publish(plugin.SnapGOBContentType, buf.Bytes(), *cfg)
	Convey("Publish metric to Riemann", t, func() {
		Convey("So err should not be returned", func() {
			So(err, ShouldBeNil)
		})
		Convey("So metric retrieved equals metric published", func() {
			c, _ := raidman.Dial("tcp", broker)
			events, _ := c.Query("host = \"bacon-powered\"")
			So(len(events), ShouldBeGreaterThan, 0)
		})
	})
}

func getProcessErrorStr(err *cpolicy.ProcessingErrors) string {
	errString := ""
	if err.HasErrors() {
		for _, e := range err.Errors() {
			errString += fmt.Sprintln(e.Error())
		}
	}
	return errString
}
