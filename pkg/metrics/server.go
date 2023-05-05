package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/nint8835/discord-vc-stats/pkg/config"
)

func ServeMetrics() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	log.Info().Str("bind_addr", config.Instance.BindAddr).Msg("Starting metrics server...")

	err := http.ListenAndServe(config.Instance.BindAddr, mux)

	if err != http.ErrServerClosed {
		return err
	}

	return nil
}
