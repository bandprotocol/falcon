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
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// metrics stores the Prometheus metrics instance.
var metrics *PrometheusMetrics

// globalTelemetryEnabled indicates whether telemetry is enabled globally.
// It is set on initialization and does not change for the lifetime of the program.
var globalTelemetryEnabled bool

// Task statuses used as labels
const (
	FinishedTaskStatus = "finished"
	ErrorTaskStatus    = "error"
	SkippedTaskStatus  = "skipped"
)

type PrometheusMetrics struct {
	PacketsRelayedSuccess      *prometheus.CounterVec
	UnrelayedPackets           *prometheus.GaugeVec
	TasksCount                 *prometheus.CounterVec
	FinishedTaskExecutionTime  *prometheus.SummaryVec
	TunnelsPerDestinationChain *prometheus.CounterVec
	ActiveTargetContractsCount *prometheus.GaugeVec
	TxsCount                   *prometheus.CounterVec
	TxProcessTime              *prometheus.SummaryVec
	GasUsed                    *prometheus.SummaryVec
}

func updateMetrics(updateFn func()) {
	if globalTelemetryEnabled {
		updateFn()
	}
}

// IncPacketsRelayedSuccess increments the count of successfully relayed packets.
func IncPacketsRelayedSuccess(tunnelID uint64) {
	updateMetrics(func() {
		metrics.PacketsRelayedSuccess.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Inc()
	})
}

// SetUnrelayedPackets sets the number of unrelayed packets.
func SetUnrelayedPackets(tunnelID uint64, unrelayedPackets uint64) {
	updateMetrics(func() {
		metrics.UnrelayedPackets.WithLabelValues(fmt.Sprintf("%d", tunnelID)).Set(float64(unrelayedPackets))
	})
}

// IncTasksCount increments the total count of executed tasks.
func IncTasksCount(tunnelID uint64, destinationChain string, taskStatus string) {
	updateMetrics(func() {
		metrics.TasksCount.WithLabelValues(fmt.Sprintf("%d", tunnelID), destinationChain, taskStatus).Inc()
	})
}

// ObserveFinishedTaskExecutionTime records the execution time (ms) of a finished task.
func ObserveFinishedTaskExecutionTime(
	tunnelID uint64,
	destinationChain string,
	finishedTaskExecutionTime int64,
) {
	updateMetrics(func() {
		metrics.FinishedTaskExecutionTime.WithLabelValues(fmt.Sprintf("%d", tunnelID), destinationChain).
			Observe(float64(finishedTaskExecutionTime))
	})
}

// IncTunnelsPerDestinationChain increments the count of tunnels per destination chain.
func IncTunnelsPerDestinationChain(destinationChain string) {
	updateMetrics(func() {
		metrics.TunnelsPerDestinationChain.WithLabelValues(destinationChain).Inc()
	})
}

// IncActiveTargetContractsCount increases the count of active target contracts.
func IncActiveTargetContractsCount(destinationChain string) {
	updateMetrics(func() {
		metrics.ActiveTargetContractsCount.WithLabelValues(destinationChain).Inc()
	})
}

// DecActiveTargetContractsCount decreases the count of active target contracts.
func DecActiveTargetContractsCount(destinationChain string) {
	updateMetrics(func() {
		metrics.ActiveTargetContractsCount.WithLabelValues(destinationChain).Dec()
	})
}

// IncTxsCount increments the transactions count per tunnel, categorized by transaction status.
func IncTxsCount(tunnelID uint64, destinationChain string, txStatus string) {
	updateMetrics(func() {
		metrics.TxsCount.WithLabelValues(fmt.Sprintf("%d", tunnelID), destinationChain, txStatus).Inc()
	})
}

// ObserveTxProcessTime records the processing time (ms) for each transaction.
func ObserveTxProcessTime(tunnelID uint64, destinationChain string, txStatus string, txProcessTime int64) {
	updateMetrics(func() {
		metrics.TxProcessTime.WithLabelValues(fmt.Sprintf("%d", tunnelID), destinationChain, txStatus).
			Observe(float64(txProcessTime))
	})
}

// ObserveGasUsed tracks the amount of gas used for each transaction.
func ObserveGasUsed(tunnelID uint64, destinationChain string, txStatus string, gasUsed decimal.NullDecimal) {
	updateMetrics(func() {
		metrics.GasUsed.WithLabelValues(fmt.Sprintf("%d", tunnelID), destinationChain, txStatus).
			Observe(gasUsed.Decimal.InexactFloat64())
	})
}

func InitPrometheusMetrics() {
	packetLabels := []string{"tunnel_id"}
	tasksCountLabels := []string{"tunnel_id", "destination_chain", "task_status"}
	finishedTaskExecutionTimeLabels := []string{"tunnel_id", "destination_chain"}
	tunnelPerDestinationChainLabels := []string{"destination_chain"}
	activeTargetContractsLabels := []string{"destination_chain"}
	txsCountLabels := []string{"tunnel_id", "destination_chain", "tx_status"}
	txProcessTimeLabels := []string{"tunnel_id", "destination_chain", "tx_status"}

	gasUsedLabels := []string{"tunnel_id", "destination_chain", "tx_status"}

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
			Help: "Total number of executed tasks",
		}, tasksCountLabels),
		FinishedTaskExecutionTime: promauto.NewSummaryVec(prometheus.SummaryOpts{
			Name: "falcon_finished_task_execution_time",
			Help: "Execution time (ms) for finished tasks",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		}, finishedTaskExecutionTimeLabels),
		TunnelsPerDestinationChain: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "falcon_tunnels_per_destination_chain",
			Help: "Total number of tunnels per destination chain",
		}, tunnelPerDestinationChainLabels),
		ActiveTargetContractsCount: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "falcon_active_target_contracts_count",
			Help: "Number of active target chain contracts",
		}, activeTargetContractsLabels),
		TxsCount: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "falcon_txs_count",
			Help: "Total number of transactions",
		}, txsCountLabels),
		TxProcessTime: promauto.NewSummaryVec(prometheus.SummaryOpts{
			Name: "falcon_tx_process_time",
			Help: "Processing time (ms) for transaction",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		}, txProcessTimeLabels),
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
			"Failed to start metrics server you can change the address and port using metrics-listen-addr config setting or --metrics-listen-addr flag",
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
