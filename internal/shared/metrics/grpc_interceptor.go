package metrics

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func GRPCUnaryServerInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		grpcType := "unary"
		method := info.FullMethod

		GRPCServerStartedTotal.WithLabelValues(serviceName, method, grpcType).Inc()
		GRPCServerMsgReceivedTotal.WithLabelValues(serviceName, method, grpcType).Inc()

		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start).Seconds()
		code := status.Code(err).String()

		if err == nil {
			GRPCServerMsgSentTotal.WithLabelValues(serviceName, method, grpcType).Inc()
		}

		GRPCServerHandledTotal.WithLabelValues(serviceName, method, grpcType, code).Inc()
		GRPCServerHandlingSeconds.WithLabelValues(serviceName, method, grpcType).Observe(duration)

		return resp, err
	}
}

func GRPCStreamServerInterceptor(serviceName string) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		grpcType := grpcStreamType(info)
		method := info.FullMethod

		GRPCServerStartedTotal.WithLabelValues(serviceName, method, grpcType).Inc()

		start := time.Now()
		err := handler(srv, &wrappedServerStream{
			ServerStream: ss,
			serviceName:  serviceName,
			method:       method,
			grpcType:     grpcType,
		})
		duration := time.Since(start).Seconds()
		code := status.Code(err).String()

		GRPCServerHandledTotal.WithLabelValues(serviceName, method, grpcType, code).Inc()
		GRPCServerHandlingSeconds.WithLabelValues(serviceName, method, grpcType).Observe(duration)

		return err
	}
}

type wrappedServerStream struct {
	grpc.ServerStream
	serviceName string
	method      string
	grpcType    string
}

func (w *wrappedServerStream) SendMsg(m interface{}) error {
	err := w.ServerStream.SendMsg(m)
	if err == nil {
		GRPCServerMsgSentTotal.WithLabelValues(
			w.serviceName,
			w.method,
			w.grpcType,
		).Inc()
	}

	return err
}

func (w *wrappedServerStream) RecvMsg(m interface{}) error {
	err := w.ServerStream.RecvMsg(m)
	if err == nil {
		GRPCServerMsgReceivedTotal.WithLabelValues(
			w.serviceName,
			w.method,
			w.grpcType,
		).Inc()
	}

	return err
}

func grpcStreamType(info *grpc.StreamServerInfo) string {
	if info.IsClientStream && info.IsServerStream {
		return "bidi_stream"
	}
	if info.IsClientStream {
		return "client_stream"
	}
	if info.IsServerStream {
		return "server_stream"
	}

	return "stream"
}
