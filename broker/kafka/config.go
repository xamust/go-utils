package kafka

import "github.com/xamust/go-utils/models_bcon"

type Config struct {
	Addr     []string             `json:"addr" yaml:"addr" mapstructure:"addr"`
	TypeAuth string               `json:"type_auth" yaml:"type_auth" mapstructure:"type_auth"`
	Username string               `json:"username" yaml:"username" mapstructure:"username"`
	Password string               `json:"password" yaml:"password" mapstructure:"password"`
	Timeout  models_bcon.Duration `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Reader   Reader               `json:"reader" yaml:"reader" mapstructure:"reader"`
	Writer   Writer               `json:"writer" yaml:"writer" mapstructure:"writer"`
}

type Reader struct {
	Workers                int                  `json:"workers"`
	Group                  string               `json:"group"`
	QueueCapacity          int                  `json:"queue_capacity"`
	MinBytes               int                  `json:"min_bytes"`
	MaxBytes               int                  `json:"max_bytes"`
	MaxWait                models_bcon.Duration `json:"max_wait"`
	ReadLagInterval        models_bcon.Duration `json:"read_lag_interval"`
	HeartbeatInterval      models_bcon.Duration `json:"heartbeat_interval"`
	CommitInterval         models_bcon.Duration `json:"commit_interval"`
	PartitionWatchInterval models_bcon.Duration `json:"partition_watch_interval"`
	SessionTimeout         models_bcon.Duration `json:"session_timeout"`
	RebalanceTimeout       models_bcon.Duration `json:"rebalance_timeout"`
	JoinGroupBackoff       models_bcon.Duration `json:"join_group_backoff"`
	RetentionTime          models_bcon.Duration `json:"retention_time"`
	StartOffset            int64                `json:"offset"`
	ReadBackoffMin         models_bcon.Duration `json:"read_backoff_min"`
	ReadBackoffMax         models_bcon.Duration `json:"read_backoff_max"`
	MaxAttempts            int                  `json:"max_attempts"`
}

type Writer struct {
	MaxAttempts int   `json:"max_attempts" yaml:"max_attempts" mapstructure:"max_attempts"`
	BatchSize   int   `json:"batch_size" yaml:"batch_size" mapstructure:"batch_size"`
	BatchBytes  int64 `json:"batch_bytes" yaml:"batch_bytes" mapstructure:"batch_bytes"`
}
