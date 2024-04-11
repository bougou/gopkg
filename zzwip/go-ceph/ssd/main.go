package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func inSSDs(osdID int) bool {
	ssdOSDs := []int{
		22, 23, 24, 25, 26, 27, 28, 29,
		30, 51, 52, 53, 54, 55, 56, 57,
		62, 73, 75, 84, 85, 86, 87, 89,
		58, 59, 60, 61, 63, 64, 65, 66,
		67, 68, 69, 70, 71, 72, 74, 76,
		77, 78, 79, 80, 81, 82, 83, 88,
	}
	for _, ssd := range ssdOSDs {
		if ssd == osdID {
			return true
		}
	}

	return false
}

func main() {

	pgdumpJson, err := ioutil.ReadFile("pgdump.json")
	if err != nil {
		panic(err)
	}

	pgdump := &PgDump{}
	if err := json.Unmarshal(pgdumpJson, pgdump); err != nil {
		panic(err)
	}

	for _, pgStat := range pgdump.PgMap.PgStats {
		fmt.Println()
		pgid := pgStat.Pgid
		if !strings.HasPrefix(pgid, "3.") {
			return
		}

		fmt.Printf("pgid: %s, act: %10v, up: %10v\n", pgid, pgStat.Acting, pgStat.Up)

		str := ""
		for _, osd := range pgStat.Acting {
			if inSSDs(osd) {
				str += fmt.Sprintf("%4d: SSD,", osd)
			} else {
				str += fmt.Sprintf("%4d: HDD,", osd)
			}
		}

		primary := pgStat.ActingPrimary
		if !inSSDs(primary) {
			str += "NOTPASSED"
		}

		for _, osd := range pgStat.Acting {
			if osd != primary {
				if inSSDs(osd) {
					str += "NOTPASSED"
					break
				}
			}
		}
		fmt.Println(str)

	}

}

type PgDump struct {
	PgReady bool `json:"pg_ready"`
	PgMap   struct {
		Version         int    `json:"version"`
		Stamp           string `json:"stamp"`
		LastOsdmapEpoch int    `json:"last_osdmap_epoch"`
		LastPgScan      int    `json:"last_pg_scan"`
		PgStatsSum      struct {
			StatSum struct {
				NumBytes                   int64 `json:"num_bytes"`
				NumObjects                 int   `json:"num_objects"`
				NumObjectClones            int   `json:"num_object_clones"`
				NumObjectCopies            int   `json:"num_object_copies"`
				NumObjectsMissingOnPrimary int   `json:"num_objects_missing_on_primary"`
				NumObjectsMissing          int   `json:"num_objects_missing"`
				NumObjectsDegraded         int   `json:"num_objects_degraded"`
				NumObjectsMisplaced        int   `json:"num_objects_misplaced"`
				NumObjectsUnfound          int   `json:"num_objects_unfound"`
				NumObjectsDirty            int   `json:"num_objects_dirty"`
				NumWhiteouts               int   `json:"num_whiteouts"`
				NumRead                    int   `json:"num_read"`
				NumReadKb                  int   `json:"num_read_kb"`
				NumWrite                   int   `json:"num_write"`
				NumWriteKb                 int   `json:"num_write_kb"`
				NumScrubErrors             int   `json:"num_scrub_errors"`
				NumShallowScrubErrors      int   `json:"num_shallow_scrub_errors"`
				NumDeepScrubErrors         int   `json:"num_deep_scrub_errors"`
				NumObjectsRecovered        int   `json:"num_objects_recovered"`
				NumBytesRecovered          int   `json:"num_bytes_recovered"`
				NumKeysRecovered           int   `json:"num_keys_recovered"`
				NumObjectsOmap             int   `json:"num_objects_omap"`
				NumObjectsHitSetArchive    int   `json:"num_objects_hit_set_archive"`
				NumBytesHitSetArchive      int   `json:"num_bytes_hit_set_archive"`
				NumFlush                   int   `json:"num_flush"`
				NumFlushKb                 int   `json:"num_flush_kb"`
				NumEvict                   int   `json:"num_evict"`
				NumEvictKb                 int   `json:"num_evict_kb"`
				NumPromote                 int   `json:"num_promote"`
				NumFlushModeHigh           int   `json:"num_flush_mode_high"`
				NumFlushModeLow            int   `json:"num_flush_mode_low"`
				NumEvictModeSome           int   `json:"num_evict_mode_some"`
				NumEvictModeFull           int   `json:"num_evict_mode_full"`
				NumObjectsPinned           int   `json:"num_objects_pinned"`
				NumLegacySnapsets          int   `json:"num_legacy_snapsets"`
				NumLargeOmapObjects        int   `json:"num_large_omap_objects"`
				NumObjectsManifest         int   `json:"num_objects_manifest"`
				NumOmapBytes               int   `json:"num_omap_bytes"`
				NumOmapKeys                int   `json:"num_omap_keys"`
				NumObjectsRepaired         int   `json:"num_objects_repaired"`
			} `json:"stat_sum"`
			StoreStats struct {
				Total                   int `json:"total"`
				Available               int `json:"available"`
				InternallyReserved      int `json:"internally_reserved"`
				Allocated               int `json:"allocated"`
				DataStored              int `json:"data_stored"`
				DataCompressed          int `json:"data_compressed"`
				DataCompressedAllocated int `json:"data_compressed_allocated"`
				DataCompressedOriginal  int `json:"data_compressed_original"`
				OmapAllocated           int `json:"omap_allocated"`
				InternalMetadata        int `json:"internal_metadata"`
			} `json:"store_stats"`
			LogSize       int `json:"log_size"`
			OndiskLogSize int `json:"ondisk_log_size"`
			Up            int `json:"up"`
			Acting        int `json:"acting"`
			NumStoreStats int `json:"num_store_stats"`
		} `json:"pg_stats_sum"`
		OsdStatsSum struct {
			UpFrom             int `json:"up_from"`
			Seq                int `json:"seq"`
			NumPgs             int `json:"num_pgs"`
			NumOsds            int `json:"num_osds"`
			NumPerPoolOsds     int `json:"num_per_pool_osds"`
			NumPerPoolOmapOsds int `json:"num_per_pool_omap_osds"`
			Kb                 int `json:"kb"`
			KbUsed             int `json:"kb_used"`
			KbUsedData         int `json:"kb_used_data"`
			KbUsedOmap         int `json:"kb_used_omap"`
			KbUsedMeta         int `json:"kb_used_meta"`
			KbAvail            int `json:"kb_avail"`
			Statfs             struct {
				Total                   int64 `json:"total"`
				Available               int64 `json:"available"`
				InternallyReserved      int   `json:"internally_reserved"`
				Allocated               int64 `json:"allocated"`
				DataStored              int64 `json:"data_stored"`
				DataCompressed          int   `json:"data_compressed"`
				DataCompressedAllocated int   `json:"data_compressed_allocated"`
				DataCompressedOriginal  int   `json:"data_compressed_original"`
				OmapAllocated           int   `json:"omap_allocated"`
				InternalMetadata        int   `json:"internal_metadata"`
			} `json:"statfs"`
			HbPeers           []interface{} `json:"hb_peers"`
			SnapTrimQueueLen  int           `json:"snap_trim_queue_len"`
			NumSnapTrimming   int           `json:"num_snap_trimming"`
			NumShardsRepaired int           `json:"num_shards_repaired"`
			OpQueueAgeHist    struct {
				Histogram  []int `json:"histogram"`
				UpperBound int   `json:"upper_bound"`
			} `json:"op_queue_age_hist"`
			PerfStat struct {
				CommitLatencyMs int `json:"commit_latency_ms"`
				ApplyLatencyMs  int `json:"apply_latency_ms"`
				CommitLatencyNs int `json:"commit_latency_ns"`
				ApplyLatencyNs  int `json:"apply_latency_ns"`
			} `json:"perf_stat"`
			Alerts           []interface{} `json:"alerts"`
			NetworkPingTimes []interface{} `json:"network_ping_times"`
		} `json:"osd_stats_sum"`
		PgStatsDelta struct {
			StatSum struct {
				NumBytes                   int `json:"num_bytes"`
				NumObjects                 int `json:"num_objects"`
				NumObjectClones            int `json:"num_object_clones"`
				NumObjectCopies            int `json:"num_object_copies"`
				NumObjectsMissingOnPrimary int `json:"num_objects_missing_on_primary"`
				NumObjectsMissing          int `json:"num_objects_missing"`
				NumObjectsDegraded         int `json:"num_objects_degraded"`
				NumObjectsMisplaced        int `json:"num_objects_misplaced"`
				NumObjectsUnfound          int `json:"num_objects_unfound"`
				NumObjectsDirty            int `json:"num_objects_dirty"`
				NumWhiteouts               int `json:"num_whiteouts"`
				NumRead                    int `json:"num_read"`
				NumReadKb                  int `json:"num_read_kb"`
				NumWrite                   int `json:"num_write"`
				NumWriteKb                 int `json:"num_write_kb"`
				NumScrubErrors             int `json:"num_scrub_errors"`
				NumShallowScrubErrors      int `json:"num_shallow_scrub_errors"`
				NumDeepScrubErrors         int `json:"num_deep_scrub_errors"`
				NumObjectsRecovered        int `json:"num_objects_recovered"`
				NumBytesRecovered          int `json:"num_bytes_recovered"`
				NumKeysRecovered           int `json:"num_keys_recovered"`
				NumObjectsOmap             int `json:"num_objects_omap"`
				NumObjectsHitSetArchive    int `json:"num_objects_hit_set_archive"`
				NumBytesHitSetArchive      int `json:"num_bytes_hit_set_archive"`
				NumFlush                   int `json:"num_flush"`
				NumFlushKb                 int `json:"num_flush_kb"`
				NumEvict                   int `json:"num_evict"`
				NumEvictKb                 int `json:"num_evict_kb"`
				NumPromote                 int `json:"num_promote"`
				NumFlushModeHigh           int `json:"num_flush_mode_high"`
				NumFlushModeLow            int `json:"num_flush_mode_low"`
				NumEvictModeSome           int `json:"num_evict_mode_some"`
				NumEvictModeFull           int `json:"num_evict_mode_full"`
				NumObjectsPinned           int `json:"num_objects_pinned"`
				NumLegacySnapsets          int `json:"num_legacy_snapsets"`
				NumLargeOmapObjects        int `json:"num_large_omap_objects"`
				NumObjectsManifest         int `json:"num_objects_manifest"`
				NumOmapBytes               int `json:"num_omap_bytes"`
				NumOmapKeys                int `json:"num_omap_keys"`
				NumObjectsRepaired         int `json:"num_objects_repaired"`
			} `json:"stat_sum"`
			StoreStats struct {
				Total                   int `json:"total"`
				Available               int `json:"available"`
				InternallyReserved      int `json:"internally_reserved"`
				Allocated               int `json:"allocated"`
				DataStored              int `json:"data_stored"`
				DataCompressed          int `json:"data_compressed"`
				DataCompressedAllocated int `json:"data_compressed_allocated"`
				DataCompressedOriginal  int `json:"data_compressed_original"`
				OmapAllocated           int `json:"omap_allocated"`
				InternalMetadata        int `json:"internal_metadata"`
			} `json:"store_stats"`
			LogSize       int    `json:"log_size"`
			OndiskLogSize int    `json:"ondisk_log_size"`
			Up            int    `json:"up"`
			Acting        int    `json:"acting"`
			NumStoreStats int    `json:"num_store_stats"`
			StampDelta    string `json:"stamp_delta"`
		} `json:"pg_stats_delta"`
		PgStats []struct {
			Pgid                    string `json:"pgid"`
			Version                 string `json:"version"`
			ReportedSeq             string `json:"reported_seq"`
			ReportedEpoch           string `json:"reported_epoch"`
			State                   string `json:"state"`
			LastFresh               string `json:"last_fresh"`
			LastChange              string `json:"last_change"`
			LastActive              string `json:"last_active"`
			LastPeered              string `json:"last_peered"`
			LastClean               string `json:"last_clean"`
			LastBecameActive        string `json:"last_became_active"`
			LastBecamePeered        string `json:"last_became_peered"`
			LastUnstale             string `json:"last_unstale"`
			LastUndegraded          string `json:"last_undegraded"`
			LastFullsized           string `json:"last_fullsized"`
			MappingEpoch            int    `json:"mapping_epoch"`
			LogStart                string `json:"log_start"`
			OndiskLogStart          string `json:"ondisk_log_start"`
			Created                 int    `json:"created"`
			LastEpochClean          int    `json:"last_epoch_clean"`
			Parent                  string `json:"parent"`
			ParentSplitBits         int    `json:"parent_split_bits"`
			LastScrub               string `json:"last_scrub"`
			LastScrubStamp          string `json:"last_scrub_stamp"`
			LastDeepScrub           string `json:"last_deep_scrub"`
			LastDeepScrubStamp      string `json:"last_deep_scrub_stamp"`
			LastCleanScrubStamp     string `json:"last_clean_scrub_stamp"`
			LogSize                 int    `json:"log_size"`
			OndiskLogSize           int    `json:"ondisk_log_size"`
			StatsInvalid            bool   `json:"stats_invalid"`
			DirtyStatsInvalid       bool   `json:"dirty_stats_invalid"`
			OmapStatsInvalid        bool   `json:"omap_stats_invalid"`
			HitsetStatsInvalid      bool   `json:"hitset_stats_invalid"`
			HitsetBytesStatsInvalid bool   `json:"hitset_bytes_stats_invalid"`
			PinStatsInvalid         bool   `json:"pin_stats_invalid"`
			ManifestStatsInvalid    bool   `json:"manifest_stats_invalid"`
			SnaptrimqLen            int    `json:"snaptrimq_len"`
			StatSum                 struct {
				NumBytes                   int `json:"num_bytes"`
				NumObjects                 int `json:"num_objects"`
				NumObjectClones            int `json:"num_object_clones"`
				NumObjectCopies            int `json:"num_object_copies"`
				NumObjectsMissingOnPrimary int `json:"num_objects_missing_on_primary"`
				NumObjectsMissing          int `json:"num_objects_missing"`
				NumObjectsDegraded         int `json:"num_objects_degraded"`
				NumObjectsMisplaced        int `json:"num_objects_misplaced"`
				NumObjectsUnfound          int `json:"num_objects_unfound"`
				NumObjectsDirty            int `json:"num_objects_dirty"`
				NumWhiteouts               int `json:"num_whiteouts"`
				NumRead                    int `json:"num_read"`
				NumReadKb                  int `json:"num_read_kb"`
				NumWrite                   int `json:"num_write"`
				NumWriteKb                 int `json:"num_write_kb"`
				NumScrubErrors             int `json:"num_scrub_errors"`
				NumShallowScrubErrors      int `json:"num_shallow_scrub_errors"`
				NumDeepScrubErrors         int `json:"num_deep_scrub_errors"`
				NumObjectsRecovered        int `json:"num_objects_recovered"`
				NumBytesRecovered          int `json:"num_bytes_recovered"`
				NumKeysRecovered           int `json:"num_keys_recovered"`
				NumObjectsOmap             int `json:"num_objects_omap"`
				NumObjectsHitSetArchive    int `json:"num_objects_hit_set_archive"`
				NumBytesHitSetArchive      int `json:"num_bytes_hit_set_archive"`
				NumFlush                   int `json:"num_flush"`
				NumFlushKb                 int `json:"num_flush_kb"`
				NumEvict                   int `json:"num_evict"`
				NumEvictKb                 int `json:"num_evict_kb"`
				NumPromote                 int `json:"num_promote"`
				NumFlushModeHigh           int `json:"num_flush_mode_high"`
				NumFlushModeLow            int `json:"num_flush_mode_low"`
				NumEvictModeSome           int `json:"num_evict_mode_some"`
				NumEvictModeFull           int `json:"num_evict_mode_full"`
				NumObjectsPinned           int `json:"num_objects_pinned"`
				NumLegacySnapsets          int `json:"num_legacy_snapsets"`
				NumLargeOmapObjects        int `json:"num_large_omap_objects"`
				NumObjectsManifest         int `json:"num_objects_manifest"`
				NumOmapBytes               int `json:"num_omap_bytes"`
				NumOmapKeys                int `json:"num_omap_keys"`
				NumObjectsRepaired         int `json:"num_objects_repaired"`
			} `json:"stat_sum"`
			Up                   []int         `json:"up"`
			Acting               []int         `json:"acting"`
			AvailNoMissing       []interface{} `json:"avail_no_missing"`
			ObjectLocationCounts []interface{} `json:"object_location_counts"`
			BlockedBy            []interface{} `json:"blocked_by"`
			UpPrimary            int           `json:"up_primary"`
			ActingPrimary        int           `json:"acting_primary"`
			PurgedSnaps          []interface{} `json:"purged_snaps"`
		} `json:"pg_stats"`
		PoolStats []struct {
			Poolid  int `json:"poolid"`
			NumPg   int `json:"num_pg"`
			StatSum struct {
				NumBytes                   int `json:"num_bytes"`
				NumObjects                 int `json:"num_objects"`
				NumObjectClones            int `json:"num_object_clones"`
				NumObjectCopies            int `json:"num_object_copies"`
				NumObjectsMissingOnPrimary int `json:"num_objects_missing_on_primary"`
				NumObjectsMissing          int `json:"num_objects_missing"`
				NumObjectsDegraded         int `json:"num_objects_degraded"`
				NumObjectsMisplaced        int `json:"num_objects_misplaced"`
				NumObjectsUnfound          int `json:"num_objects_unfound"`
				NumObjectsDirty            int `json:"num_objects_dirty"`
				NumWhiteouts               int `json:"num_whiteouts"`
				NumRead                    int `json:"num_read"`
				NumReadKb                  int `json:"num_read_kb"`
				NumWrite                   int `json:"num_write"`
				NumWriteKb                 int `json:"num_write_kb"`
				NumScrubErrors             int `json:"num_scrub_errors"`
				NumShallowScrubErrors      int `json:"num_shallow_scrub_errors"`
				NumDeepScrubErrors         int `json:"num_deep_scrub_errors"`
				NumObjectsRecovered        int `json:"num_objects_recovered"`
				NumBytesRecovered          int `json:"num_bytes_recovered"`
				NumKeysRecovered           int `json:"num_keys_recovered"`
				NumObjectsOmap             int `json:"num_objects_omap"`
				NumObjectsHitSetArchive    int `json:"num_objects_hit_set_archive"`
				NumBytesHitSetArchive      int `json:"num_bytes_hit_set_archive"`
				NumFlush                   int `json:"num_flush"`
				NumFlushKb                 int `json:"num_flush_kb"`
				NumEvict                   int `json:"num_evict"`
				NumEvictKb                 int `json:"num_evict_kb"`
				NumPromote                 int `json:"num_promote"`
				NumFlushModeHigh           int `json:"num_flush_mode_high"`
				NumFlushModeLow            int `json:"num_flush_mode_low"`
				NumEvictModeSome           int `json:"num_evict_mode_some"`
				NumEvictModeFull           int `json:"num_evict_mode_full"`
				NumObjectsPinned           int `json:"num_objects_pinned"`
				NumLegacySnapsets          int `json:"num_legacy_snapsets"`
				NumLargeOmapObjects        int `json:"num_large_omap_objects"`
				NumObjectsManifest         int `json:"num_objects_manifest"`
				NumOmapBytes               int `json:"num_omap_bytes"`
				NumOmapKeys                int `json:"num_omap_keys"`
				NumObjectsRepaired         int `json:"num_objects_repaired"`
			} `json:"stat_sum"`
			StoreStats struct {
				Total                   int `json:"total"`
				Available               int `json:"available"`
				InternallyReserved      int `json:"internally_reserved"`
				Allocated               int `json:"allocated"`
				DataStored              int `json:"data_stored"`
				DataCompressed          int `json:"data_compressed"`
				DataCompressedAllocated int `json:"data_compressed_allocated"`
				DataCompressedOriginal  int `json:"data_compressed_original"`
				OmapAllocated           int `json:"omap_allocated"`
				InternalMetadata        int `json:"internal_metadata"`
			} `json:"store_stats"`
			LogSize       int `json:"log_size"`
			OndiskLogSize int `json:"ondisk_log_size"`
			Up            int `json:"up"`
			Acting        int `json:"acting"`
			NumStoreStats int `json:"num_store_stats"`
		} `json:"pool_stats"`
		OsdStats []struct {
			Osd                int   `json:"osd"`
			UpFrom             int   `json:"up_from"`
			Seq                int64 `json:"seq"`
			NumPgs             int   `json:"num_pgs"`
			NumOsds            int   `json:"num_osds"`
			NumPerPoolOsds     int   `json:"num_per_pool_osds"`
			NumPerPoolOmapOsds int   `json:"num_per_pool_omap_osds"`
			Kb                 int   `json:"kb"`
			KbUsed             int   `json:"kb_used"`
			KbUsedData         int   `json:"kb_used_data"`
			KbUsedOmap         int   `json:"kb_used_omap"`
			KbUsedMeta         int   `json:"kb_used_meta"`
			KbAvail            int   `json:"kb_avail"`
			Statfs             struct {
				Total                   int64 `json:"total"`
				Available               int64 `json:"available"`
				InternallyReserved      int   `json:"internally_reserved"`
				Allocated               int64 `json:"allocated"`
				DataStored              int64 `json:"data_stored"`
				DataCompressed          int   `json:"data_compressed"`
				DataCompressedAllocated int   `json:"data_compressed_allocated"`
				DataCompressedOriginal  int   `json:"data_compressed_original"`
				OmapAllocated           int   `json:"omap_allocated"`
				InternalMetadata        int   `json:"internal_metadata"`
			} `json:"statfs"`
			HbPeers           []interface{} `json:"hb_peers"`
			SnapTrimQueueLen  int           `json:"snap_trim_queue_len"`
			NumSnapTrimming   int           `json:"num_snap_trimming"`
			NumShardsRepaired int           `json:"num_shards_repaired"`
			OpQueueAgeHist    struct {
				Histogram  []int `json:"histogram"`
				UpperBound int   `json:"upper_bound"`
			} `json:"op_queue_age_hist"`
			PerfStat struct {
				CommitLatencyMs int `json:"commit_latency_ms"`
				ApplyLatencyMs  int `json:"apply_latency_ms"`
				CommitLatencyNs int `json:"commit_latency_ns"`
				ApplyLatencyNs  int `json:"apply_latency_ns"`
			} `json:"perf_stat"`
			Alerts           []interface{} `json:"alerts"`
			NetworkPingTimes []interface{} `json:"network_ping_times"`
		} `json:"osd_stats"`
	} `json:"pg_map"`
}
