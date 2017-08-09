package verify

import (
	"fmt"
	"os"

	"../../lib/prometheus"
	"../../lib/tarantool"
	"../../utils/schemes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{FieldMap: logrus.FieldMap{logrus.FieldKeyTime: "date", logrus.FieldKeyLevel: "type"}})
	logrus.SetOutput(os.Stdout)
}

func VerifyRefill(c *gin.Context) {
	var request httpschemes.FondyRefill
	if c.BindJSON(&request) != nil {
		c.JSON(400, gin.H{"message": "Ошибка при валидации JSON"})
		return
	}
	if request.ValidateSignature() != nil {
		logrus.WithFields(logrus.Fields{"module": "VerifyRefill"}).Error(request)
		c.JSON(400, gin.H{"message": "Токен не валиден"})
		return
	}
	t := tarantool.GetTarantool()
	res, err := t.Call("refillUser", []interface{}{request.ProductID, request.Amount, request.PaymentID})
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "VerifyRefill"}).Error(err)
		c.JSON(500, gin.H{"message": "Произошла ошибка при обработке платежа"})
		return
	}
	if res.Data[0].([]interface{})[0].(bool) == true {
		logrus.WithFields(logrus.Fields{"module": "VerifyRefill"}).Info(fmt.Sprintf("product_id - %s, amount - %s, payment_id - %d", request.ProductID, request.Amount, request.PaymentID))
		metrics.CountPayments.Inc()
		c.JSON(200, gin.H{"message": "Успешное зачисление"})
	} else {
		logrus.WithFields(logrus.Fields{"module": "VerifyRefill"}).Error("Ошибка при обработке платежа")
		c.JSON(500, gin.H{"message": "Произошла ошибка при обработке платежа"})
	}
}
