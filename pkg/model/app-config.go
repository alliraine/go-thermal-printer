package model

type AppConfig struct {
	Server   ServerConfig  `toml:"server"`
	Printer  PrinterConfig `toml:"printer"`
	TestMode bool          `toml:"test_mode" default:"false"`
	USBMode  bool          `toml:"usb_mode" default:"false"`
}

type ServerConfig struct {
	Host        string `toml:"host" default:"0.0.0.0"`
	Port        int    `toml:"port" default:"8080"`
	ApiKey      string `toml:"api_key" default:""`
	SwaggerHost string `toml:"swagger_host" default:"localhost:8080"`
}

type PrinterConfig struct {
	Port     string `toml:"port" default:"/dev/ttyUSB0"`
	BaudRate int    `toml:"baud_rate" default:"19200"`
	DataBits int    `toml:"data_bits" default:"8"`
	StopBits int    `toml:"stop_bits" default:"1"`
	Parity   int    `toml:"parity" default:"0"`
}
