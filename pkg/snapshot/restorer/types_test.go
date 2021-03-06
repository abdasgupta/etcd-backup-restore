// Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
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

package restorer_test

import (
	"fmt"
	"net/url"
	"time"

	"github.com/coreos/etcd/pkg/types"
	"github.com/gardener/etcd-backup-restore/pkg/snapstore"
	. "github.com/gardener/etcd-backup-restore/pkg/snapshot/restorer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("restorer types", func() {
	var (
		makeRestorationConfig = func(s string, b bool, i int) *RestorationConfig {
			return &RestorationConfig{
				InitialCluster: s,
				InitialClusterToken: s,
				RestoreDataDir: s,
				InitialAdvertisePeerURLs: []string{ s, s },
				Name: s,
				SkipHashCheck: b,
				MaxFetchers: uint(i),
				EmbeddedEtcdQuotaBytes: int64(i),
			}
		}
		makeSnap = func(s string, i int, t time.Time, b bool) snapstore.Snapshot {
			return snapstore.Snapshot{
				Kind: s,
				StartRevision: int64(i),
				LastRevision: int64(i),
				CreatedOn: t,
				SnapDir: s,
				SnapName: s,
				IsChunk: b,
			}
		}
		makeSnapList = func(s string, i int, t time.Time, b bool) snapstore.SnapList {
			var s1, s2 = makeSnap(s, i, t, b), makeSnap(s, i, t, b)
			return snapstore.SnapList{ &s1, &s2 }
		}
		makeURL = func(s string, b bool) url.URL {
			return url.URL{
				Scheme: s,
				Opaque: s,
				User: url.UserPassword(s, s),
				Host: s,
				Path: s,
				RawPath: s,
				ForceQuery: b,
				Fragment: s,
			}
		}
		makeURLs = func(s string, b bool) types.URLs {
			return types.URLs{ makeURL(s, b), makeURL(s, b) }
		}
		makeURLsMap = func(s string, b bool) types.URLsMap {
			var out = types.URLsMap{}
			for _, v := range []int{1, 2} {
				out[fmt.Sprintf("%s-%d", s, v)] = makeURLs(s, b)
			}
			return out
		}
		makeRestoreOptions = func(s string, i int, t time.Time, b bool) *RestoreOptions {
			return &RestoreOptions{
				Config: makeRestorationConfig(s, b, i),
				ClusterURLs: makeURLsMap(s, b),
				PeerURLs: makeURLs(s, b),
				BaseSnapshot: makeSnap(s, i, t, b),
				DeltaSnapList: makeSnapList(s, i, t, b),
			}
		}
	)

	Describe("RestorationConfig", func() {
		var (
			makeA = func() *RestorationConfig { return makeRestorationConfig("a", false, 1) }
			makeB = func() *RestorationConfig { return makeRestorationConfig("b", true, 2) }
		)
		Describe("DeepCopyInto", func() {
			It("new out", func() {
				var a, in, out = makeA(), makeA(), new(RestorationConfig)
				in.DeepCopyInto(out)
				Expect(out).To(Equal(in))
				Expect(out).ToNot(BeIdenticalTo(in))
				Expect(in).To(Equal(a))
				Expect(in).ToNot(BeIdenticalTo(a))
			})
			It("existing out", func() {
				var a, in, b, out = makeA(), makeA(), makeB(), makeB()
				in.DeepCopyInto(out)
				Expect(out).To(Equal(in))
				Expect(out).ToNot(BeIdenticalTo(in))
				Expect(in).To(Equal(a))
				Expect(in).ToNot(BeIdenticalTo(a))
				Expect(out).ToNot(Equal(b))
			})
		})
		Describe("DeepCopy", func() {
			It("out", func() {
				var a, in = makeA(), makeA()
				var out = in.DeepCopy()
				Expect(out).ToNot(BeNil())
				Expect(out).To(Equal(in))
				Expect(out).ToNot(BeIdenticalTo(in))
				Expect(in).To(Equal(a))
				Expect(in).ToNot(BeIdenticalTo(a))
			})
		})
	})

	Describe("SnapList", func() {
		var (
			now = time.Now()
			makeA = func() snapstore.SnapList { return makeSnapList("a", 1, now, false) }
		)
		Describe("DeepCopySnapList", func() {
			It("out", func() {
				var a, in = makeA(), makeA()
				var out = DeepCopySnapList(in)
				Expect(out).ToNot(BeNil())
				Expect(out).To(Equal(in))
				Expect(out).ToNot(BeIdenticalTo(in))
				Expect(in).To(Equal(a))
				Expect(in).ToNot(BeIdenticalTo(a))
			})
		})
	})

	Describe("URL", func() {
		var (
			makeA = func() *url.URL { var u = makeURL("a", false); return &u }
		)
		Describe("DeepCopyURL", func() {
			It("out", func() {
				var a, in = makeA(), makeA()
				var out = DeepCopyURL(in)
				Expect(out).ToNot(BeNil())
				Expect(out).To(Equal(in))
				Expect(out).ToNot(BeIdenticalTo(in))
				Expect(in).To(Equal(a))
				Expect(in).ToNot(BeIdenticalTo(a))
			})
		})
	})

	Describe("URLs", func() {
		var (
			makeA = func() types.URLs { return makeURLs("a", false) }
		)
		Describe("DeepCopyURLs", func() {
			It("out", func() {
				var a, in = makeA(), makeA()
				var out = DeepCopyURLs(in)
				Expect(out).ToNot(BeNil())
				Expect(out).To(Equal(in))
				Expect(out).ToNot(BeIdenticalTo(in))
				Expect(in).To(Equal(a))
				Expect(in).ToNot(BeIdenticalTo(a))
			})
		})
	})

	Describe("RestoreOptions", func() {
		var (
			now = time.Now()
			makeA = func() *RestoreOptions { return makeRestoreOptions("a", 1, now, false) }
			makeB = func() *RestoreOptions { return makeRestoreOptions("b", 2, now.Add(-1*time.Hour), true) }
		)
		Describe("DeepCopyInto", func() {
			It("new out", func() {
				var a, in, out = makeA(), makeA(), new(RestoreOptions)
				in.DeepCopyInto(out)
				Expect(out).To(Equal(in))
				Expect(out).ToNot(BeIdenticalTo(in))
				Expect(in).To(Equal(a))
				Expect(in).ToNot(BeIdenticalTo(a))
			})
			It("existing out", func() {
				var a, in, b, out = makeA(), makeA(), makeB(), makeB()
				in.DeepCopyInto(out)
				Expect(out).To(Equal(in))
				Expect(out).ToNot(BeIdenticalTo(in))
				Expect(in).To(Equal(a))
				Expect(in).ToNot(BeIdenticalTo(a))
				Expect(out).ToNot(Equal(b))
			})
		})
		Describe("DeepCopy", func() {
			It("out", func() {
				var a, in = makeA(), makeA()
				var out = in.DeepCopy()
				Expect(out).ToNot(BeNil())
				Expect(out).To(Equal(in))
				Expect(out).ToNot(BeIdenticalTo(in))
				Expect(in).To(Equal(a))
				Expect(in).ToNot(BeIdenticalTo(a))
			})
		})
	})

})