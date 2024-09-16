package kafka_logger

import (
	"github.com/xamust/go-utils/models_bcon"
	"time"
)

// ------- kafka logger -------

type Stage string

const (
	KafelTopic = "synapse_kf"

	TestStage Stage = "test"
	// todo @ need fill config
	//preprod stage = "preprod"
	//prod    stage = "prod"

	K8sPodName       = "k8s.pod_name"
	K8sContainerName = "k8s.container_name"
)

type KafkaLogConfig struct {
	Config Config `mapstructure:"config" json:"config" yaml:"config"`
	Topic  string `mapstructure:"topic" json:"topic" yaml:"topic"`
}

type Config struct {
	Addr     []string             `json:"addr" yaml:"addr" mapstructure:"addr"`
	TypeAuth string               `json:"type_auth" yaml:"type_auth" mapstructure:"type_auth"`
	Username string               `json:"username" yaml:"username" mapstructure:"username"`
	Password string               `json:"password" yaml:"password" mapstructure:"password"`
	Timeout  models_bcon.Duration `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Writer   Writer               `json:"writer" yaml:"writer" mapstructure:"writer"`
}

type Writer struct {
	MaxAttempts int   `json:"max_attempts" yaml:"max_attempts" mapstructure:"max_attempts"`
	BatchSize   int   `json:"batch_size" yaml:"batch_size" mapstructure:"batch_size"`
	BatchBytes  int64 `json:"batch_bytes" yaml:"batch_bytes" mapstructure:"batch_bytes"`
}

var configs = map[Stage]KafkaLogConfig{
	TestStage: {
		Config: Config{
			Addr:     []string{"10.4.113.72:9092", "10.4.113.73:9092", "10.4.113.74:9092"},
			TypeAuth: "plain",
			Username: "synapse_usr",
			Password: "FDJ0ybMy",
			Timeout: models_bcon.Duration{
				Duration: 10 * time.Second,
			},
		},
		Topic: KafelTopic,
	},
	// todo @ need fill config
	//preprod: {},
	//prod:    {},
}

func GetConfig(stage Stage) KafkaLogConfig {
	return configs[stage]
}
