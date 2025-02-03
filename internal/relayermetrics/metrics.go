package relayermetrics

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var metrics *PrometheusMetrics

// globalTelemetryEnabled is a private variable that stores the telemetry enabled state.
// It is set on initialization and does not change for the lifetime of the program.
var globalTelemetryEnabled bool

// IsTelemetryEnabled provides controlled access to check if telemetry is enabled.
func IsTelemetryEnabled() bool {
	return globalTelemetryEnabled
}

// EnableTelemetry allows for the global telemetry enabled state to be set.
func EnableTelemetry() {
	globalTelemetryEnabled = true
}

type PrometheusMetrics struct {
	Registry              *prometheus.Registry
	TunnelCount           prometheus.Counter
	PacketReceived        *prometheus.CounterVec
	UnrelayedPacket       *prometheus.GaugeVec
	TasksCount            *prometheus.CounterVec
	TaskExecutionTime     *prometheus.SummaryVec
	DestinationChainCount prometheus.Counter
	TargetContract        *prometheus.GaugeVec
	TxCount               *prometheus.CounterVec
	TxProcessTime         *prometheus.SummaryVec
	GasUsed               *prometheus.SummaryVec
}

func AddTunnellCount(count uint64) {
	metrics.TunnelCount.Add(float64(count))
}

func IncPacketlReceived(tunnelID uint64) {
	metrics.PacketReceived.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Inc()
}

func SetUnrelayedPacket(tunnelID uint64, unrelayedPacket float64) {
	metrics.UnrelayedPacket.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Set(unrelayedPacket)
}

func IncTasksCount(tunnelID uint64) {
	metrics.TasksCount.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Inc()
}

func ObserveTaskExecutionTime(tunnelID uint64, taskExecutionTime float64) {
	metrics.TaskExecutionTime.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Observe(taskExecutionTime)
}

func AddDestinationChainCount(count uint64) {
	metrics.DestinationChainCount.Add(float64(count))
}

func IncTargetContractCount(status string) {
	metrics.TargetContract.WithLabelValues(status).Inc()
}

func DecTargetContractCount(status string) {
	metrics.TargetContract.WithLabelValues(status).Dec()
}

func IncTxCount(tunnelID uint64) {
	metrics.TxCount.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Inc()
}

func ObserveTxProcessTime(chainName string, taskExecutionTime float64) {
	metrics.TxProcessTime.WithLabelValues(chainName).Observe(taskExecutionTime)
}

func ObserveGasUsed(tunnelID uint64, gasUsed uint64) {
	metrics.GasUsed.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Observe(float64(gasUsed))
}

func NewPrometheusMetrics() *PrometheusMetrics {
	tunnelLabels := []string{"tunnel_id"}
	targetChainContractLabels := []string{"status"}
	txCountLabels := []string{"tunnel_id"}
	txProcessTimeLabels := []string{"chain_name"}
	gasUsedLabels := []string{"tunnel_id"}

	registry := prometheus.NewRegistry()
	registerer := promauto.With(registry)
	metrics = &PrometheusMetrics{
		Registry: registry,
		TunnelCount: registerer.NewCounter(prometheus.CounterOpts{
			Name: "falcon_tunnel_count_total",
			Help: "Total number of observed tunnels",
		}),
		PacketReceived: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "falcon_packet_received_total",
			Help: "Total number of packets received",
		}, tunnelLabels),
		UnrelayedPacket: registerer.NewGaugeVec(prometheus.GaugeOpts{
			Name: "falcon_unrelayed_packet_count",
			Help: "Number of unrelayed packets",
		}, tunnelLabels),
		TasksCount: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "falcon_task_count_total",
			Help: "Total number of observed tasks",
		}, tunnelLabels),
		TaskExecutionTime: registerer.NewSummaryVec(prometheus.SummaryOpts{
			Name: "falcon_task_execution_time",
			Help: "Task execution time in milliseconds",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		}, tunnelLabels),
		DestinationChainCount: registerer.NewCounter(prometheus.CounterOpts{
			Name: "falcon_destination_chain_count_total",
			Help: "Total number of destination chains",
		}),
		TargetContract: registerer.NewGaugeVec(prometheus.GaugeOpts{
			Name: "falcon_target_chain_contract_count",
			Help: "Number of target chain contracts",
		}, targetChainContractLabels),
		TxCount: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "falcon_tx_count_total",
			Help: "Total number of transactions per tunnel",
		}, txCountLabels),
		TxProcessTime: registerer.NewSummaryVec(prometheus.SummaryOpts{
			Name: "falcon_tx_process_time",
			Help: "Transaction processing time in milliseconds",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		}, txProcessTimeLabels),
		GasUsed: registerer.NewSummaryVec(prometheus.SummaryOpts{
			Name: "falcon_gas_used_per_tx",
			Help: "Gas used per transaction",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		}, gasUsedLabels),
	}
	return metrics
}

// StartMetricsServer starts a metrics server in a background goroutine,
// accepting connections on the given listener.
// Any HTTP logging will be written at info level to the given logger.
// The server will be forcefully shut down when ctx finishes.
func StartMetricsServer(ctx context.Context, log *zap.Logger, ln net.Listener, registry *prometheus.Registry) {
	// Set up new mux identical to the default mux configuration in net/http/pprof.
	mux := http.NewServeMux()

	// Serve default prometheus metrics
	mux.Handle("/metrics", promhttp.Handler())

	// Serve relayer metrics
	mux.Handle("/relayer/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	srv := &http.Server{
		Handler:  mux,
		ErrorLog: zap.NewStdLog(log),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		_ = srv.Serve(ln)
	}()

	go func() {
		<-ctx.Done()
		srv.Close()
	}()
}
