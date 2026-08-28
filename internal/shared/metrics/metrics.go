package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	AppInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_info",
			Help: "Application info",
		},
		[]string{"service", "version", "go_version"},
	)
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"service", "method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"service", "method", "path", "status"},
	)

	HTTPRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
		[]string{"service"},
	)

	HTTPResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000, 500000, 1000000},
		},
		[]string{"service", "method", "path"},
	)
)

var (
	GRPCServerStartedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_server_started_total",
			Help: "Total number of gRPC calls started on the server",
		},
		[]string{"service", "grpc_method", "grpc_type"},
	)

	GRPCServerHandledTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_server_handled_total",
			Help: "Total number of gRPC calls completed on the server",
		},
		[]string{"service", "grpc_method", "grpc_type", "grpc_code"},
	)

	GRPCServerHandlingSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_server_handling_seconds",
			Help:    "gRPC call handling duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"service", "grpc_method", "grpc_type"},
	)

	GRPCServerMsgReceivedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_server_msg_received_total",
			Help: "Total number of gRPC messages received by the server",
		},
		[]string{"service", "grpc_method", "grpc_type"},
	)

	GRPCServerMsgSentTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_server_msg_sent_total",
			Help: "Total number of gRPC messages sent by the server",
		},
		[]string{"service", "grpc_method", "grpc_type"},
	)
)

var (
	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"service", "operation", "table"},
	)

	DBOpenConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "db_open_connections",
			Help: "Number of open database connections",
		},
		[]string{"service", "db_type"},
	)
)

var (
	RabbitMQMessagesPublished = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rabbitmq_messages_published_total",
			Help: "Total number of messages published to RabbitMQ",
		},
		[]string{"service", "exchange", "routing_key"},
	)

	RabbitMQMessagesConsumed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rabbitmq_messages_consumed_total",
			Help: "Total number of messages consumed from RabbitMQ",
		},
		[]string{"service", "queue", "status"},
	)

	RabbitMQMessageProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "rabbitmq_message_processing_seconds",
			Help:    "RabbitMQ message processing duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"service", "queue"},
	)
)
