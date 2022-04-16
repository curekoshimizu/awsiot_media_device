package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	mediaDevices "github.com/antonfisher/go-media-devices-state"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

var (
	configPath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "awsiot_media_device",
	Short: "AWS IoT Media Device State Publisher",
	Run: func(cmd *cobra.Command, args []string) {
		main()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&configPath, "config", "config.yaml", "config file.")
}

func main() {
	configYaml, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	var config Config
	if err := yaml.Unmarshal(configYaml, &config); err != nil {
		panic(err)
	}
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		panic(err)
	}

	tlsConfig, err := newTLSConfig(config)
	if err != nil {
		panic(err)
	}

	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("ssl://%s:%d", config.EndPoint, config.Port)).
		SetTLSConfig(tlsConfig)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("failed to connect broker: %v", token.Error()))
	}
	defer client.Disconnect(250)

	// if token := client.Publish(config.Topic, config.QoS, false, pubMessage); token.Wait() && token.Error() != nil {
	// 	fmt.Printf("failed to publish %s: %v\n", config.Topic, token.Error())
	// }

	if err := mainLoop(client, config); err != nil {
		panic(err)
	}
}

func mainLoop(client mqtt.Client, config Config) error {
	var isOn *bool

	for {
		isCameraOn, err := mediaDevices.IsCameraOn()
		if err != nil {
			return err
		}
		isMicrophoneOn, err := mediaDevices.IsMicrophoneOn()
		if err != nil {
			return err
		}

		requiredState := isCameraOn || isMicrophoneOn
		if isOn == nil || *isOn != requiredState {
			var topic string
			if requiredState {
				topic = config.Event.On
			} else {
				topic = config.Event.Off
			}
			if err := publish(topic, client, config); err != nil {
				fmt.Printf("failed to publish %s: %v\n", topic, err)
			} else {
				// Successfully updated
				isOn = &requiredState
			}
		}

		time.Sleep(time.Second * 1)
	}

}

func publish(event string, client mqtt.Client, config Config) error {
	if token := client.Publish(config.Topic, config.QoS, false, event); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func newTLSConfig(config Config) (*tls.Config, error) {
	rootCA, err := ioutil.ReadFile(config.RootCAFilePath)
	if err != nil {
		return nil, err
	}
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(rootCA)
	cert, err := tls.LoadX509KeyPair(config.CertificateFilePath, config.PrivateKeyFilePath)
	if err != nil {
		return nil, err
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		RootCAs:            certpool,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
		NextProtos:         []string{"x-amzn-mqtt-ca"},
	}, nil
}
