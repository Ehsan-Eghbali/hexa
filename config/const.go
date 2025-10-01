package config

import "time"

const (
	QueryCacheTtl    = 10 * time.Minute
	QueryCacheIsMust = false
	NatsWorkerTopic  = "hexa_nats_topic"
)
