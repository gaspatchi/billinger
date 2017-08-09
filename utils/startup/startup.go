package startup

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "date", logrus.FieldKeyLevel: "type"}})
	logrus.SetOutput(os.Stdout)
}

func register(agent *api.Agent, name string, address string, port int) {
	service := api.AgentServiceRegistration{Name: name, Address: address, Port: port}
	err := agent.ServiceRegister(&service)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "startup"}).Panic("Ошибка при регистрации сервиса")
		panic(err)
	}
}

func deregister(agent *api.Agent, name string) {
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for {
			s := <-signal_chan
			if len(s.String()) > 1 {
				err := agent.ServiceDeregister(name)
				if err != nil {
					logrus.WithFields(logrus.Fields{"module": "startup"}).Panic("Ошибка при дерегистрации сервиса")
				}
				os.Exit(0)
			}
		}
	}()
}

func InitService() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Print(err)
	}
	consulconfig := api.DefaultConfig()
	consulconfig.Address = fmt.Sprintf("%s:%d", viper.GetString("consul.address"), viper.GetInt("consul.port"))
	client, err := api.NewClient(consulconfig)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "startup"}).Panic("Ошибка при подключении к Consul")
		panic(err)
	}
	agent := client.Agent()
	register(agent, viper.GetString("server.name"), viper.GetString("server.address"), viper.GetInt("server.address"))
	deregister(agent, viper.GetString("server.name"))
}
