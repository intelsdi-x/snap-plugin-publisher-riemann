//
// +build integration

package riemann

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/intelsdi-x/pulse/control/plugin"
	"github.com/intelsdi-x/pulse/control/plugin/cpolicy"
	"github.com/intelsdi-x/pulse/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/amir/raidman"
)

// integration test
func TestRiemannPublish(t *testing.T) {
	// This integration test requires a Riemann Server
	broker := os.Getenv("PULSE_TEST_RIEMANN")
	if broker == "" {
		fmt.Println("Skipping integration tests")
		return
	}

	var buf bytes.Buffer
	buf.Reset()
	r := NewRiemannPublisher()
	config := make(map[string]ctypes.ConfigValue)
	config["broker"] = ctypes.ConfigValueStr{Value: broker}
	cp := r.GetConfigPolicy()
	cfg, _ := cp.Get([]string{""}).Process(config)
	metrics := []plugin.PluginMetricType{
		*plugin.NewPluginMetricType([]string{"intel", "cpu", "temp"}, time.Now(), "bacon-powered", 100),
	}
	enc := gob.NewEncoder(&buf)
	enc.Encode(metrics)
	err := r.Publish(plugin.PulseGOBContentType, buf.Bytes(), *cfg)
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
