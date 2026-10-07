/*
Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved.

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

// Package worker contains functions used at the worker controller
package worker

import (
	hcloud "github.com/hetznercloud/hcloud-go/v2/hcloud"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ActuatorHooks", func() {
	Describe("#checkServerTypeAvailableInLocation", func() {
		It("should succeed if the server type is available in the location", func() {
			serverType := &hcloud.ServerType{
				Name: "cpx42",
				Locations: []hcloud.ServerTypeLocation{
					{Location: &hcloud.Location{Name: "fsn1"}, Available: false},
					{Location: &hcloud.Location{Name: "nbg1"}, Available: true},
				},
			}

			Expect(checkServerTypeAvailableInLocation(serverType, "nbg1")).To(Succeed())
		})

		It("should fail if the server type is present but not available in the location", func() {
			serverType := &hcloud.ServerType{
				Name: "cpx42",
				Locations: []hcloud.ServerTypeLocation{
					{Location: &hcloud.Location{Name: "nbg1"}, Available: false},
				},
			}

			Expect(checkServerTypeAvailableInLocation(serverType, "nbg1")).To(MatchError(ContainSubstring("machine type cpx42 is currently not available in nbg1")))
		})

		It("should fail if the location is missing", func() {
			serverType := &hcloud.ServerType{
				Name: "cpx42",
				Locations: []hcloud.ServerTypeLocation{
					{Location: &hcloud.Location{Name: "fsn1"}, Available: true},
				},
			}

			Expect(checkServerTypeAvailableInLocation(serverType, "nbg1")).To(HaveOccurred())
		})

		It("should not panic on a nil location entry", func() {
			serverType := &hcloud.ServerType{
				Name: "cpx42",
				Locations: []hcloud.ServerTypeLocation{
					{Location: nil, Available: true},
					{Location: &hcloud.Location{Name: "nbg1"}, Available: true},
				},
			}

			Expect(func() {
				Expect(checkServerTypeAvailableInLocation(serverType, "nbg1")).To(Succeed())
			}).NotTo(Panic())
			Expect(checkServerTypeAvailableInLocation(&hcloud.ServerType{
				Name:      "cpx42",
				Locations: []hcloud.ServerTypeLocation{{Location: nil, Available: true}},
			}, "nbg1")).To(HaveOccurred())
		})
	})
})
