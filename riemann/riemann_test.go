//
// +build unit

package riemann

import (
	"testing"

	"github.com/intelsdi-x/pulse/control/plugin"
	"github.com/intelsdi-x/pulse/control/plugin/cpolicy"
	"github.com/intelsdi-x/pulse/core/ctypes"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRiemannPlugin(t *testing.T) {
	Convey("Meta returns proper metadata", t, func() {
		meta := Meta()
		So(meta.Name, ShouldResemble, PluginName)
		So(meta.Version, ShouldResemble, PluginVersion)
		So(meta.Type, ShouldResemble, plugin.PublisherPluginType)
	})

	Convey("Create Riemann Publisher", t, func() {
		rp := NewRiemannPublisher()
		Convey("So Riemann Publisher should not be nil", func() {
			So(rp, ShouldNotBeNil)
		})
		Convey("So Riemann Publisher shoud be of type riemannPublisher", func() {
			So(rp, ShouldHaveSameTypeAs, &riemannPublisher{})
		})
		configPolicy, err := rp.GetConfigPolicy()
		Convey("GetConfigPolicy() should return a config policy", func() {
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})
			Convey("So GetConfigPolicy() should not return an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("So config policy should be of cpolicy.ConfigPolicy type", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})
			testConfig := make(map[string]ctypes.ConfigValue)
			testConfig["broker"] = ctypes.ConfigValueStr{Value: "127.0.0.1:5555"}
			cfg, errs := configPolicy.Get([]string{""}).Process(testConfig)
			Convey("So config policy should process testConfig and return a config", func() {
				So(cfg, ShouldNotBeNil)
			})
			Convey("So testConfig processing should return no errors", func() {
				So(errs.HasErrors(), ShouldBeFalse)
			})
		})
	})
}
