package auth

import (
	"fmt"
	"os"
	"regexp"

	"../../lib/prometheus"
	"../../utils/schemes"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v0"
)

var re *regexp.Regexp
var endpoint string

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "date", logrus.FieldKeyLevel: "type"}})
	logrus.SetOutput(os.Stdout)
	re = regexp.MustCompile(`Bearer\s(\S+)`)
	client, err := api.NewClient(api.DefaultConfig())
	catalog := client.Catalog()
	response, _, err := catalog.Service("tokenzer", "", &api.QueryOptions{Datacenter: "dc1"})
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "VerifyAuth"}).Panic("Consul недоступен")
		panic(err)
	}
	if len(response) == 0 {
		logrus.WithFields(logrus.Fields{"module": "VerifyAuth"}).Panic("Tokenzer не зарегистрирован")
		panic(err)
	} else {
		endpoint = fmt.Sprintf("http://%s:%d/verify", response[0].ServiceAddress, response[0].ServicePort)
	}
}

func VerifyAuth(c *gin.Context) {
	var tokenzer httpschemes.AuthResponse
	var str = c.Request.Header.Get("Authorization")
	var result = re.FindStringSubmatch(str)
	if len(result) != 2 {
		c.JSON(400, gin.H{"message": "Заголовок отсутсвует"})
		c.Abort()
		return
	}
	response, err := resty.R().SetBody(httpschemes.AuthRequest{Token: result[1]}).SetResult(&tokenzer).Post(endpoint)
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "VerifyAuth"}).Error("Ошибка при запросе к Tokenzer")
		c.JSON(500, gin.H{"message": "Невозможно проверить сессию"})
		c.Abort()
		return
	}
	if response.StatusCode() == 200 {
		metrics.CountVerifyAuth.Inc()
		c.Set("email", tokenzer.Sub)
		c.Next()
	} else if response.StatusCode() == 500 {
		logrus.WithFields(logrus.Fields{"module": "VerifyAuth"}).Error("Ошибка Tokenzer'a")
		c.Data(response.StatusCode(), "application/json", response.Body())
		c.Abort()
		return
	}
}
