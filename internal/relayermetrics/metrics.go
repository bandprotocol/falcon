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
	PacketsRelayedSuccess      *prometheus.CounterVec
	UnrelayedPackets           *prometheus.GaugeVec
	TasksCount                 *prometheus.CounterVec
	TaskExecutionTime          *prometheus.SummaryVec
	TunnelsPerDestinationChain *prometheus.CounterVec
	ActiveTargetContractsCount prometheus.Gauge
	TxsCount                   *prometheus.CounterVec
	TxProcessTime              *prometheus.SummaryVec
	GasUsed                    *prometheus.SummaryVec
}

func updateMetrics(updateFn func()) {
	if globalTelemetryEnabled {
		updateFn()
	}
}

// IncPacketsRelayedSuccess increments the count of successfully relayed packets for a specific tunnel.
func IncPacketsRelayedSuccess(tunnelID uint64) {
	updateMetrics(func() {
		metrics.PacketsRelayedSuccess.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Inc()
	})
}

// SetUnrelayedPackets sets the number of unrelayed packets for a specific tunnel.
func SetUnrelayedPackets(tunnelID uint64, unrelayedPackets float64) {
	updateMetrics(func() {
		metrics.UnrelayedPackets.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Set(unrelayedPackets)
	})
}

// IncTasksCount increments the total tasks count for a specific tunnel.
func IncTasksCount(tunnelID uint64) {
	updateMetrics(func() {
		metrics.TasksCount.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Inc()
	})
}

// ObserveTaskExecutionTime records the execution time of a task for a specific tunnel.
func ObserveTaskExecutionTime(tunnelID uint64, taskExecutionTime float64) {
	updateMetrics(func() {
		metrics.TaskExecutionTime.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Observe(taskExecutionTime)
	})
}

// IncTunnelsPerDestinationChain increments the total number of tunnels per specific destination chain.
func IncTunnelsPerDestinationChain(destinationChain string) {
	updateMetrics(func() {
		metrics.TunnelsPerDestinationChain.WithLabelValues(destinationChain).Inc()
	})
}

// IncActiveTargetContractsCount increases the count of active target contracts.
func IncActiveTargetContractsCount() {
	updateMetrics(func() {
		metrics.ActiveTargetContractsCount.Inc()
	})
}

// DecActiveTargetContractsCount decreases the count of active target contracts.
func DecActiveTargetContractsCount() {
	updateMetrics(func() {
		metrics.ActiveTargetContractsCount.Dec()
	})
}

// IncTxsCount increments the transactions count metric for a specific tunnel.
func IncTxsCount(tunnelID uint64) {
	updateMetrics(func() {
		metrics.TxsCount.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Inc()
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

func InitPrometheusMetrics() {
	packetLabels := []string{"tunnel_id"}
	taskLabels := []string{"tunnel_id"}
	tunnelPerDestinationChainLabels := []string{"destination_chain"}
	txLabels := []string{"tunnel_id"}
	gasUsedLabels := []string{"tunnel_id"}

	metrics = &PrometheusMetrics{
		PacketsRelayedSuccess: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "falcon_packets_relayed_success",
			Help: "Total number of packets successfully relayed from BandChain",
		}, packetLabels),
		UnrelayedPackets: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "falcon_unrelayed_packets",
			Help: "Number of unrelayed packets (the difference between total packets from BandChain and received packets from the target chain)",
		}, packetLabels),
		TasksCount: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "falcon_tasks_count",
			Help: "Total number of successfully executed tasks",
		}, taskLabels),
		TaskExecutionTime: promauto.NewSummaryVec(prometheus.SummaryOpts{
			Name: "falcon_task_execution_time",
			Help: "Task execution time in milliseconds",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		}, taskLabels),
		TunnelsPerDestinationChain: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "falcon_tunnels_per_destination_chain",
			Help: "Total number of destination chains",
		}, tunnelPerDestinationChainLabels),
		ActiveTargetContractsCount: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "falcon_active_target_contracts_count",
			Help: "Number of active target chain contracts",
		}),
		TxsCount: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "falcon_txs_count",
			Help: "Total number of transactions per tunnel",
		}, txLabels),
		TxProcessTime: promauto.NewSummaryVec(prometheus.SummaryOpts{
			Name: "falcon_tx_process_time",
			Help: "Transaction processing time in milliseconds",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		}, txLabels),
		GasUsed: promauto.NewSummaryVec(prometheus.SummaryOpts{
			Name: "falcon_gas_used",
			Help: "Amount of gas used per transaction",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		}, gasUsedLabels),
	}
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

	// allow for the global telemetry enabled state to be set.
	globalTelemetryEnabled = true

	// initialize Prometheus metrics
	InitPrometheusMetrics()

	// set up new mux identical to the default mux configuration in net/http/pprof.
	mux := http.NewServeMux()

	// serve prometheus metrics
	mux.Handle("/metrics", promhttp.Handler())

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
