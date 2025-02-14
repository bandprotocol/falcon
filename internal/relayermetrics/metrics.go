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

// metrics stores the Prometheus metrics instance.
var metrics *PrometheusMetrics

// globalTelemetryEnabled indicates whether telemetry is enabled globally.
// It is set on initialization and does not change for the lifetime of the program.
var globalTelemetryEnabled bool

type PrometheusMetrics struct {
	Registry                  *prometheus.Registry
	PacketReceived            *prometheus.CounterVec
	UnrelayedPacket           *prometheus.GaugeVec
	TaskCount                 *prometheus.CounterVec
	TaskExecutionTime         *prometheus.SummaryVec
	TunnelPerDestinationChain *prometheus.CounterVec
	ActiveTargetContract      prometheus.Gauge
	TxCount                   *prometheus.CounterVec
	TxProcessTime             *prometheus.SummaryVec
	GasUsed                   *prometheus.SummaryVec
}

func updateMetrics(updateFn func()) {
	if globalTelemetryEnabled {
		updateFn()
	}
}

// IncPacketlReceived increments the count of successfully relayed packets for a specific tunnel.
func IncPacketlReceived(tunnelID uint64) {
	updateMetrics(func() {
		metrics.PacketReceived.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Inc()
	})
}

// SetUnrelayedPacket sets the number of unrelayed packets for a specific tunnel.
func SetUnrelayedPacket(tunnelID uint64, unrelayedPacket float64) {
	updateMetrics(func() {
		metrics.UnrelayedPacket.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Set(unrelayedPacket)
	})
}

// IncTaskCount increments the total task count for a specific tunnel.
func IncTaskCount(tunnelID uint64) {
	updateMetrics(func() {
		metrics.TaskCount.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Inc()
	})
}

// ObserveTaskExecutionTime records the execution time of a task for a specific tunnel.
func ObserveTaskExecutionTime(tunnelID uint64, taskExecutionTime float64) {
	updateMetrics(func() {
		metrics.TaskExecutionTime.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Observe(taskExecutionTime)
	})
}

// IncTunnelPerDestinationChain increments the total number of tunnels per specify destination chain.
func IncTunnelPerDestinationChain(destinationChain string) {
	updateMetrics(func() {
		metrics.TunnelPerDestinationChain.WithLabelValues(destinationChain).Inc()
	})
}

// IncActiveTargetContractCount increases the count of active target contracts.
func IncActiveTargetContractCount() {
	updateMetrics(func() {
		metrics.ActiveTargetContract.Inc()
	})
}

// DecActiveTargetContractCount decreases the count of active target contracts.
func DecActiveTargetContractCount() {
	updateMetrics(func() {
		metrics.ActiveTargetContract.Dec()
	})
}

// IncTxCount increments the transaction count metric for a specific tunnel.
func IncTxCount(tunnelID uint64) {
	updateMetrics(func() {
		metrics.TxCount.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Inc()
	})
}

// ObserveTxProcessTime tracks transaction processing time in seconds with millisecond precision.
func ObserveTxProcessTime(destinationChain string, txProcessTime float64) {
	updateMetrics(func() {
		metrics.TxProcessTime.WithLabelValues(destinationChain).Observe(txProcessTime)
	})
}

// ObserveGasUsed tracks gas used for the each relayed transaction.
func ObserveGasUsed(tunnelID uint64, gasUsed uint64) {
	updateMetrics(func() {
		metrics.GasUsed.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Observe(float64(gasUsed))
	})
}

func NewPrometheusMetrics() *PrometheusMetrics {
	tunnelLabels := []string{"tunnel_id"}
	destinationChainLabels := []string{"destination_chain"}

	registry := prometheus.NewRegistry()
	registerer := promauto.With(registry)
	metrics = &PrometheusMetrics{
		Registry: registry,
		PacketReceived: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "falcon_packet_received_total",
			Help: "Total number of packets received",
		}, tunnelLabels),
		UnrelayedPacket: registerer.NewGaugeVec(prometheus.GaugeOpts{
			Name: "falcon_unrelayed_packet_count",
			Help: "Number of unrelayed packets",
		}, tunnelLabels),
		TaskCount: registerer.NewCounterVec(prometheus.CounterOpts{
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
		TunnelPerDestinationChain: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "falcon_tunnel_per_destination_chain",
			Help: "Total number of destination chains",
		}, destinationChainLabels),
		ActiveTargetContract: registerer.NewGauge(prometheus.GaugeOpts{
			Name: "falcon_active_target_chain_contract_count",
			Help: "Number of active target chain contracts",
		}),
		TxCount: registerer.NewCounterVec(prometheus.CounterOpts{
			Name: "falcon_tx_count_total",
			Help: "Total number of transactions per tunnel",
		}, tunnelLabels),
		TxProcessTime: registerer.NewSummaryVec(prometheus.SummaryOpts{
			Name: "falcon_tx_process_time",
			Help: "Transaction processing time in milliseconds",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		}, destinationChainLabels),
		GasUsed: registerer.NewSummaryVec(prometheus.SummaryOpts{
			Name: "falcon_gas_used_per_tx",
			Help: "Gas used per transaction",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		}, tunnelLabels),
	}
	return metrics
}

// StartMetricsServer starts a metrics server in a background goroutine,
// accepting connections on the given listener.
// Any HTTP logging will be written at info level to the given logger.
// The server will be forcefully shut down when ctx finishes.
func StartMetricsServer(ctx context.Context, log *zap.Logger, metricsListenAddr string) error {
	ln, err := net.Listen("tcp", metricsListenAddr)
	if err != nil {
		log.Error(
			"Failed to start metrics server you can change the address and port using metrics-listen-addr config setting or --metrics-listen-flag",
		)

		return fmt.Errorf("failed to listen on metrics address %q: %w", metricsListenAddr, err)
	}
	log = log.With(zap.String("sys", "metricshttp"))
	log.Info("Metrics server listening", zap.String("addr", metricsListenAddr))

	// Allow for the global telemetry enabled state to be set.
	globalTelemetryEnabled = true

	prometheusMetrics := NewPrometheusMetrics()

	// Set up new mux identical to the default mux configuration in net/http/pprof.
	mux := http.NewServeMux()

	// Serve default prometheus metrics
	mux.Handle("/metrics", promhttp.Handler())

	// Serve relayer metrics
	mux.Handle("/relayer/metrics", promhttp.HandlerFor(prometheusMetrics.Registry, promhttp.HandlerOpts{}))

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

	return nil
}
