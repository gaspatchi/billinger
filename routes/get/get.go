package get

import (
	"os"
	"time"

	"../../lib/prometheus"
	"../../lib/tarantool"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "date", logrus.FieldKeyLevel: "type"}})
	logrus.SetOutput(os.Stdout)
}

func GetStatement(c *gin.Context) {
	t := tarantool.GetTarantool()
	email, exists := c.Get("email")
	if exists != true {
		c.JSON(500, gin.H{"message": "Невозможно получить стэйтмент"})
	}
	start := time.Now()
	response, err := t.Call("selectBilling", []interface{}{email})
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "GetStatement"}).Error(err)
		c.JSON(500, gin.H{"message": "Невозможно получить стэйтмент"})
		c.Abort()
		return
	} else {
		stop := time.Since(start)
		metrics.TarantoolResponseTime.Set(float64(stop.Nanoseconds() / 10000))
		metrics.CountSelects.Inc()
		c.Data(200, "application/json; charset=utf-8", []byte(response.Data[1].([]interface{})[0].(string)))
	}
}
