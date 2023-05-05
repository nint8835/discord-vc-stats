package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var MembersInVoice = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "discord_vc_stats_members_in_voice",
	},
	[]string{
		"channel",
		"user",
	},
)
