package prometheusMetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(BasicAuthAttempts, BasicAuthSuccess, JWTAuthAttempts, JWTAuthSuccess)
}

var BasicAuthAttempts = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "basic_auth_login_attempts",
		Help: "Number of basic authentication requests",
	},
)

var BasicAuthSuccess = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "basic_auth_login_success",
		Help: "Number of VALID basic authentication requests",
	},
)

var JWTAuthAttempts = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "jwt_auth_attempts",
		Help: "Number of jwt authentication requests",
	},
)

var JWTAuthSuccess = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "jwt_auth_success",
		Help: "Number of VALID jwt authentication requests",
	},
)
