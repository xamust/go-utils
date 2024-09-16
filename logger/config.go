package logger

type Config struct {
	Level      string `json:"loglevel" yaml:"loglevel" mapstructure:"loglevel"`
	MaxMsgSize int    `json:"max_message_size" yaml:"max_message_size" mapstructure:"max_message_size"`
}
