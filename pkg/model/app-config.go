package model

type AppConfig struct {
	Server  ServerConfig  `toml:"server"`
	Printer PrinterConfig `toml:"printer"`
}

type ServerConfig struct {
	Host string `toml:"host" default:"127.0.0.1"`
	Port int    `toml:"port" default:"8080"`
}

type PrinterConfig struct {
	Port     string `toml:"port" default:"/dev/ttyUSB0"`
	BaudRate int    `toml:"baud_rate" default:"19200"`
	DataBits int    `toml:"data_bits" default:"8"`
	StopBits int    `toml:"stop_bits" default:"1"`
	Parity   int    `toml:"parity" default:"0"`
}
