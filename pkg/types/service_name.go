package golibtypes

type ServiceName string

const (
	ServiceNameAuth                   = ServiceName("auth")
)

type ServiceBaseUrl string

const (
	ServiceAuthBaseUrl           = ServiceBaseUrl("paas-auth-grpc-svc")
)

func GetBaseUrl(serviceName ServiceName) ServiceBaseUrl {
	var baseUrl ServiceBaseUrl
	switch serviceName {
	case ServiceNameAuth:
		baseUrl = ServiceAuthBaseUrl
	}
	return baseUrl
}
