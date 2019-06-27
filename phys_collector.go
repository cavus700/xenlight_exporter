package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/xen-project/xen/tools/golang/xenlight"
)

/*
domain_cpu_time_total
domain_vcpu_time_total
domain_memory_max_bytes
domain_memory_current_bytes
domain_memory_outstanding_bytes
domain_cpu_count
domain_cpu_online_count

version_info
*/
var (
	physTopologyNodesDesc = prometheus.NewDesc(
		"xen_physical_topology_nodes_number",
		"Number of socket on the host",
		nil, nil,
	)
	physTopologyCoresDesc = prometheus.NewDesc(
		"xen_physical_topology_cores_per_socket",
		"Number of cores per socket on the host",
		nil, nil,
	)
	physTopologyThreadsDesc = prometheus.NewDesc(
		"xen_physical_topology_threads_per_core",
		"Number of threads per core on the host",
		nil, nil,
	)
	physMemoryTotalDesc = prometheus.NewDesc(
		"xen_physical_memory_total_bytes",
		"Total ammount of RAM on the host",
		nil, nil,
	)
	physMemoryFreeDesc = prometheus.NewDesc(
		"xen_physical_memory_free_bytes",
		"Total ammount of free RAM on the host",
		nil, nil,
	)
	physMemoryScrubDesc = prometheus.NewDesc(
		"xen_physical_memory_scrub_bytes",
		"Total ammount of scrub RAM on the host",
		nil, nil,
	)
	physMemoryOutstandingDesc = prometheus.NewDesc(
		"xen_physical_memory_outstanding_bytes",
		"Total ammount of outstanding RAM on the host",
		nil, nil,
	)
)

type PhysicalCollector struct{}

func init() {
	registerCollector("physical", defaultEnabled, NewPhysicalCollector)
}

func NewPhysicalCollector() prometheus.Collector {
	return &PhysicalCollector{}
}

func (collector PhysicalCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(collector, ch)
}

func (collector PhysicalCollector) Collect(ch chan<- prometheus.Metric) {
	xenlight.Ctx.Open()
	physinfo, err := xenlight.Ctx.GetPhysinfo()
	if err != nil {
		return
	}
	versinfo, err := xenlight.Ctx.GetVersionInfo()
	if err != nil {
		return
	}
	pageSize := versinfo.Pagesize
	ch <- prometheus.MustNewConstMetric(
		physTopologyNodesDesc,
		prometheus.GaugeValue,
		float64(physinfo.NrNodes),
	)
	ch <- prometheus.MustNewConstMetric(
		physTopologyCoresDesc,
		prometheus.GaugeValue,
		float64(physinfo.CoresPerSocket),
	)
	ch <- prometheus.MustNewConstMetric(
		physTopologyThreadsDesc,
		prometheus.GaugeValue,
		float64(physinfo.ThreadsPerCore),
	)
	ch <- prometheus.MustNewConstMetric(
		physMemoryTotalDesc,
		prometheus.GaugeValue,
		float64(physinfo.TotalPages*uint64(pageSize)),
	)
	ch <- prometheus.MustNewConstMetric(
		physMemoryFreeDesc,
		prometheus.GaugeValue,
		float64(physinfo.FreePages*uint64(pageSize)),
	)
	ch <- prometheus.MustNewConstMetric(
		physMemoryScrubDesc,
		prometheus.GaugeValue,
		float64(physinfo.ScrubPages*uint64(pageSize)),
	)
	ch <- prometheus.MustNewConstMetric(
		physMemoryOutstandingDesc,
		prometheus.GaugeValue,
		float64(physinfo.OutstandingPages*uint64(pageSize)),
	)
}
