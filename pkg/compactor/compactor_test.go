package compactor_test

import (
	"fmt"
	"os"

	"github.com/gardener/etcd-backup-restore/pkg/compactor"
	"github.com/gardener/etcd-backup-restore/pkg/miscellaneous"
	"github.com/gardener/etcd-backup-restore/pkg/snapstore"
	brtypes "github.com/gardener/etcd-backup-restore/pkg/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.etcd.io/etcd/pkg/types"
)

var _ = Describe("Running Compactor", func() {
	var (
		dir             string
		store           brtypes.SnapStore
		cptr            *compactor.Compactor
		restorePeerURLs []string
		clusterUrlsMap  types.URLsMap
		peerUrls        types.URLs
		// deltaSnapshotPeriod time.Duration
	)
	const (
		restoreName            string = "default"
		restoreClusterToken    string = "etcd-cluster"
		restoreCluster         string = "default=http://localhost:2380"
		skipHashCheck          bool   = false
		maxFetchers            uint   = 6
		maxCallSendMsgSize            = 2 * 1024 * 1024 //2Mib
		maxRequestBytes               = 2 * 1024 * 1024 //2Mib
		maxTxnOps                     = 2 * 1024
		embeddedEtcdQuotaBytes int64  = 8 * 1024 * 1024 * 1024
	)

	BeforeEach(func() {
		//wg = &sync.WaitGroup{}
		restorePeerURLs = []string{"http://localhost:2380"}
		clusterUrlsMap, err = types.NewURLsMap(restoreCluster)
		Expect(err).ShouldNot(HaveOccurred())
		peerUrls, err = types.NewURLs(restorePeerURLs)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("Compact while a etcd server is running", func() {
		var restoreOpts *brtypes.RestoreOptions

		BeforeEach(func() {
			dir = fmt.Sprintf("%s/etcd/snapshotter.bkp", testSuitDir)

			store, err = snapstore.GetSnapstore(&brtypes.SnapstoreConfig{Container: dir, Provider: "Local"})
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Println("The store where compaction will save snapshot is: ", store)

			cptr = compactor.NewCompactor(store, logger)
			restoreOpts = &brtypes.RestoreOptions{
				Config: &brtypes.RestorationConfig{
					InitialClusterToken:      restoreClusterToken,
					InitialCluster:           restoreCluster,
					Name:                     restoreName,
					InitialAdvertisePeerURLs: restorePeerURLs,
					SkipHashCheck:            skipHashCheck,
					MaxFetchers:              maxFetchers,
					MaxCallSendMsgSize:       maxCallSendMsgSize,
					MaxRequestBytes:          maxRequestBytes,
					MaxTxnOps:                maxTxnOps,
					EmbeddedEtcdQuotaBytes:   embeddedEtcdQuotaBytes,
				},
				ClusterURLs: clusterUrlsMap,
				PeerURLs:    peerUrls,
			}
		})

		AfterEach(func() {
		})

		Context("with defragmention allowed", func() {
			It("should create a snapshot", func() {
				restoreOpts.Config.MaxFetchers = 4

				// Fetch latest snapshots
				baseSnapshot, deltaSnapList, err := miscellaneous.GetLatestFullSnapshotAndDeltaSnapList(store)
				Expect(err).ShouldNot(HaveOccurred())

				restoreOpts.BaseSnapshot = baseSnapshot
				restoreOpts.DeltaSnapList = deltaSnapList

				// Take the compacted full snapshot with defragmnetation allowed
				res, err := cptr.Compact(restoreOpts, true)
				Expect(err).ShouldNot(HaveOccurred())

				// Check if the compacted full snapshot is really present
				// fi, err := os.Stat(filepath.Join(dir, res.Snapshot.SnapDir, res.Snapshot.SnapName))
				fi, err := os.Stat(res.Path)
				Expect(err).ShouldNot(HaveOccurred())

				size := fi.Size()
				Expect(size).ShouldNot(BeZero())

				err = os.Remove(res.Path)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("with defragmention not allowed", func() {
			It("should create a snapshot", func() {
				restoreOpts.Config.MaxFetchers = 4

				// Fetch latest snapshots
				baseSnapshot, deltaSnapList, err := miscellaneous.GetLatestFullSnapshotAndDeltaSnapList(store)
				Expect(err).ShouldNot(HaveOccurred())

				restoreOpts.BaseSnapshot = baseSnapshot
				restoreOpts.DeltaSnapList = deltaSnapList

				// Take the compacted full snapshot with defragmnetation not allowed
				res, err := cptr.Compact(restoreOpts, false)
				Expect(err).ShouldNot(HaveOccurred())

				// Check if the compacted full snapshot is really present
				// fi, err := os.Stat(filepath.Join(dir, res.Snapshot.SnapDir, res.Snapshot.SnapName))
				fi, err := os.Stat(res.Path)
				Expect(err).ShouldNot(HaveOccurred())

				size := fi.Size()
				Expect(size).ShouldNot(BeZero())

				err = os.Remove(res.Path)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		/*Context("with just one full compacted snapshot in backup store", func() {
			It("should create a compacted snapshot", func() {
				restoreOpts.Config.MaxFetchers = 4

				// Fetch latest snapshots
				baseSnapshot, deltaSnapList, err := miscellaneous.GetLatestFullSnapshotAndDeltaSnapList(store)
				Expect(err).ShouldNot(HaveOccurred())

				restoreOpts.BaseSnapshot = baseSnapshot
				restoreOpts.DeltaSnapList = deltaSnapList

				// Take the first compacted full snapshot
				res, err := cptr.Compact(restoreOpts, true)
				Expect(err).ShouldNot(HaveOccurred())

				// Check if the compacted full snapshot is really present
				fi, err := os.Stat(filepath.Join(dir, res.Snapshot.SnapDir, res.Snapshot.SnapName))
				Expect(err).ShouldNot(HaveOccurred())

				size := fi.Size()
				Expect(size).ShouldNot(BeZero())

				// Fetch latest snapshots which includes only a compacted full snapshot
				baseSnapshot, deltaSnapList, err = miscellaneous.GetLatestFullSnapshotAndDeltaSnapList(store)
				Expect(err).ShouldNot(HaveOccurred())
				logger.Infof("Base snapshotpath is %v", restoreOpts.BaseSnapshot.SnapDir)

				restoreOpts.BaseSnapshot = baseSnapshot
				restoreOpts.DeltaSnapList = deltaSnapList

				// Take the compacted full snapshot again
				res, err = cptr.Compact(restoreOpts, true)
				Expect(err).ShouldNot(HaveOccurred())

				// Check if the compacted full snapshot is really present
				fi, err = os.Stat(filepath.Join(dir, res.Snapshot.SnapDir, res.Snapshot.SnapName))

				Expect(err).ShouldNot(HaveOccurred())

				size = fi.Size()
				Expect(size).ShouldNot(BeZero())
			})
		})
		Context("with one full compacted snapshot and subsequent delta snapshots in backup store", func() {
			It("should create a compacted snapshot", func() {
				restoreOpts.Config.MaxFetchers = 4

				// Fetch latest snapshots
				baseSnapshot, deltaSnapList, err := miscellaneous.GetLatestFullSnapshotAndDeltaSnapList(store)
				Expect(err).ShouldNot(HaveOccurred())

				restoreOpts.BaseSnapshot = baseSnapshot
				restoreOpts.DeltaSnapList = deltaSnapList

				// Take the first compacted full snapshot
				res, err := cptr.Compact(restoreOpts, true)
				Expect(err).ShouldNot(HaveOccurred())

				// Check if the compacted full snapshot is really present
				fi, err := os.Stat(filepath.Join(dir, res.Snapshot.SnapDir, res.Snapshot.SnapName))
				Expect(err).ShouldNot(HaveOccurred())

				size := fi.Size()
				Expect(size).ShouldNot(BeZero())

				// Then populate the etcd with some more data to add the subsequent delta snapshots
				populatorCtx, cancelPopulator := context.WithTimeout(testCtx, 15*time.Second)
				defer cancelPopulator()
				resp := &utils.EtcdDataPopulationResponse{}
				wg := &sync.WaitGroup{}
				wg.Add(1)
				go utils.PopulateEtcdWithWaitGroup(populatorCtx, wg, logger, endpoints, resp)
				Expect(resp.Err).ShouldNot(HaveOccurred())

				// Start the Snapshotter with compressionPolicy = "gzip" to take delta snapshots
				compressionConfig := compressor.NewCompressorConfig()
				compressionConfig.Enabled = true
				compressionConfig.CompressionPolicy = "gzip"
				ctx, cancel := context.WithTimeout(testCtx, 6*time.Second)
				defer cancel()
				deltaSnapshotPeriod = time.Second
				err = utils.RunSnapshotter(logger, dir, deltaSnapshotPeriod, endpoints, ctx.Done(), false, compressionConfig)
				Expect(err).ShouldNot(HaveOccurred())

				// Wait unitil the populator finishes with populating ETCD
				wg.Wait()

				// Fetch the latest snapshots which are actually one compacted full snapshot and subsequent delta snapshots
				baseSnapshot, deltaSnapList, err = miscellaneous.GetLatestFullSnapshotAndDeltaSnapList(store)
				Expect(err).ShouldNot(HaveOccurred())

				restoreOpts.BaseSnapshot = baseSnapshot
				restoreOpts.DeltaSnapList = deltaSnapList

				// Take the compacted full snapshot
				res, err = cptr.Compact(restoreOpts, true)
				Expect(err).ShouldNot(HaveOccurred())

				// Check if the compacted full snapshot is really present
				fi, err = os.Stat(filepath.Join(dir, res.Snapshot.SnapDir, res.Snapshot.SnapName))
				Expect(err).ShouldNot(HaveOccurred())

				size = fi.Size()
				Expect(size).ShouldNot(BeZero())
			})
		})*/
		Context("with no basesnapshot in backup store", func() {
			It("should not run compaction", func() {
				restoreOpts.Config.MaxFetchers = 4

				// Fetch the latest snapshots which are one compacted full snapshot and subsequent delta snapshots
				_, deltaSnapList, err := miscellaneous.GetLatestFullSnapshotAndDeltaSnapList(store)
				Expect(err).ShouldNot(HaveOccurred())

				// But set the BaseSnapshot as nil
				restoreOpts.BaseSnapshot = nil
				restoreOpts.DeltaSnapList = deltaSnapList

				// Try capturing the compacted full snapshot
				_, err = cptr.Compact(restoreOpts, false)
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
