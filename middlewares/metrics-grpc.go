package middlewares

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/JohnSalazar/microservices-go-common/metrics"
	"github.com/JohnSalazar/microservices-go-common/services"
	"google.golang.org/grpc"
)

func StreamServerInterceptorMetrics(service services.Metrics) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		status := http.StatusOK

		err := handler(srv, stream)
		if err != nil {
			log.Println(err)
		}

		appMetric := metrics.NewHttpMetrics(info.FullMethod, "POST")
		appMetric.Started()
		appMetric.Finished()
		appMetric.StatusCode = strconv.Itoa(status)
		service.SaveHttp(appMetric)

		return err
	}
}

func UnaryServerInterceptorMetrics(service services.Metrics) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		status := http.StatusOK

		resp, err = handler(ctx, req)
		if err != nil {
			log.Println(err)
		}

		appMetric := metrics.NewHttpMetrics(info.FullMethod, "POST")
		appMetric.Started()
		appMetric.Finished()
		appMetric.StatusCode = strconv.Itoa(status)
		service.SaveHttp(appMetric)

		return resp, err
	}
}
