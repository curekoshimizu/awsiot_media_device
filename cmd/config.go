package cmd

type Config struct {
	QoS                 byte   `yaml:"QoS" validate:"required"`
	Topic               string `yaml:"topic" validate:"required"`
	EndPoint            string `yaml:"endpoint" validate:"required"`
	Port                int    `yaml:"port" validate:"required"`
	RootCAFilePath      string `yaml:"root_ca" validate:"required"`
	PrivateKeyFilePath  string `yaml:"private_key" validate:"required"`
	CertificateFilePath string `yaml:"certificate" validate:"required"`
}
