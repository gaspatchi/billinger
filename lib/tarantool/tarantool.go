package tarantool

import (
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tarantool/go-tarantool"
)

var user string
var endpoint string
var trntl *tarantool.Connection

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "date", logrus.FieldKeyLevel: "type"}})
	logrus.SetOutput(os.Stdout)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Tarantool"}).Panic("Ошибка при чтении конфига")
		fmt.Print(err)
	}
	client, err := api.NewClient(api.DefaultConfig())
	catalog := client.Catalog()
	response, _, err := catalog.Service("tarantool", "", &api.QueryOptions{Datacenter: "dc1"})
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Tarantool"}).Panic("Consul недоступен")
		panic(err)
	}
	if len(response) == 0 {
		logrus.WithFields(logrus.Fields{"module": "Tarantool"}).Panic("Tarantool не зарегистрирован")
		panic(err)
	} else {
		endpoint = fmt.Sprintf("%s:%d", response[0].ServiceAddress, response[0].ServicePort)
	}
	setup(viper.GetString("tarantool.user"), viper.GetString("tarantool.password"))
}

func setup(user string, password string) {
	opts := tarantool.Opts{User: user, Pass: password, MaxReconnects: 5, Reconnect: 1}
	conn, err := tarantool.Connect(endpoint, opts)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "Tarantool"}).Panic("Ошибка при подключении к Tarantool")
		panic(err)
	}
	trntl = conn
}

func GetTarantool() *tarantool.Connection {
	return trntl
}
