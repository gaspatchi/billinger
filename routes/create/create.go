package create

import (
	"fmt"
	"strconv"

	"../../lib/prometheus"
	"../../utils/schemes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/resty.v0"
)

func CreateToken(c *gin.Context) {
	email, exists := c.Get("email")
	if exists != true {
		c.JSON(500, gin.H{"message": "Невозможно проверить сессию"})
		c.Abort()
		return
	}
	var body httpschemes.CreateRequest
	if c.BindJSON(&body) != nil {
		c.JSON(400, gin.H{"message": "Ошибка при валидации JSON"})
		c.Abort()
		return
	}
	var request httpschemes.FondyCreate
	request.Request.Amount = strconv.FormatFloat(body.Amount*100, 'f', -1, 64)
	request.Request.Currency = "RUB"
	request.Request.MerchantID = "1400963"
	request.Request.OrderDesc = fmt.Sprintf("Пополнение аккаунта на %.2f рублей", body.Amount)
	request.Request.ProductID = email.(string)
	request.Request.ServerCallbackURL = "https://gt.gaspatchi.ru/billing/verify"
	request.Request.SenderEmail = email.(string)
	token, err := request.CreateSignature()
	if err != nil {
		c.JSON(500, gin.H{"message": "Невозможно создать форму оплаты"})
		return
	}
	request.Request.Signature = token

	var response httpschemes.FondyResponse
	_, err = resty.R().SetBody(request).SetResult(&response).Post("https://api.fondy.eu/api/checkout/url/")
	if err != nil {
		logrus.WithFields(logrus.Fields{"module": "CreateToken"}).Error(" ")
		c.JSON(500, gin.H{"message": "Не удалось получить ответ от платёжной системы"})
		return
	}
	if response.Response.ErrorCode == 0 {
		metrics.CountCreates.Inc()
		c.JSON(200, gin.H{"url": response.Response.CheckoutURL, "id": response.Response.PaymentID})
	} else {
		c.JSON(500, gin.H{"message": "Произошла ошибка при обработке запроса"})
	}
}
