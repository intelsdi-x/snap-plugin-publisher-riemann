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
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/amir/raidman"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
)

const (
	PluginName    = "riemann"
	PluginVersion = 8
	PluginType    = plugin.PublisherPluginType
)

// Meta returns the metadata details for the Riemann Publisher Plugin
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(PluginName, PluginVersion, PluginType, []string{plugin.SnapGOBContentType}, []string{plugin.SnapGOBContentType})
}

type riemannPublisher struct{}

// NewRiemannPublisher does something cool
func NewRiemannPublisher() *riemannPublisher {
	return &riemannPublisher{}
}

// GetConfigPolicy returns the config policy for the Riemann Publisher Plugin
func (r *riemannPublisher) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cp := cpolicy.New()
	config := cpolicy.NewPolicyNode()

	// Riemann server to publish event to
	r1, err := cpolicy.NewStringRule("broker", true)
	handleErr(err)
	r1.Description = "Broker in the format of broker-ip:port (ex: 192.168.1.1:5555)"

	config.Add(r1)
	cp.Add([]string{""}, config)
	return cp, nil
}

// Publish serializes the data and calls publish to send events to Riemann
func (r *riemannPublisher) Publish(contentType string, content []byte, config map[string]ctypes.ConfigValue) error {
	logger := log.New()
	//err := r.publish(event, broker)
	//return err
	logger.Println("Riemann Publishing Started")
	var metrics []plugin.MetricType
	switch contentType {
	case plugin.SnapGOBContentType:
		dec := gob.NewDecoder(bytes.NewBuffer(content))
		if err := dec.Decode(&metrics); err != nil {
			logger.Printf("Error decoding: error=%v content=%v", err, content)
			return err
		}
	default:
		logger.Printf("Error unknown content type '%v'", contentType)
		return errors.New(fmt.Sprintf("Unknown content type '%s'", contentType))
	}
	logger.Printf("publishing %v to %v", metrics, config)
	for _, m := range metrics {
		e := createEvent(m, config)
		if err := r.publish(e, config["broker"].(ctypes.ConfigValueStr).Value); err != nil {
			logger.Println(err)
			return err
		}
	}
	return nil
}

// publish sends events to riemann
func (r *riemannPublisher) publish(event *raidman.Event, broker string) error {
	c, err := raidman.Dial("tcp", broker)
	defer c.Close()
	if err != nil {
		return err
	}
	return c.Send(event)
}

func createEvent(m plugin.MetricType, config map[string]ctypes.ConfigValue) *raidman.Event {
	return &raidman.Event{
		Host:    m.Tags()[core.STD_TAG_PLUGIN_RUNNING_ON],
		Service: m.Namespace().String(),
		Metric:  m.Data(),
	}
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}
