package httpschemes

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var token string

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Print(err)
	}
	token = viper.GetString("server.fondytoken")
}

type FondyRefill struct {
	Rrn                     string `json:"rrn"`
	MaskedCard              string `json:"masked_card"`
	SenderCellPhone         string `json:"sender_cell_phone"`
	ResponseSignatureString string `json:"response_signature_string"`
	ResponseStatus          string `json:"response_status"`
	SenderAccount           string `json:"sender_account"`
	Fee                     string `json:"fee"`
	RectokenLifetime        string `json:"rectoken_lifetime"`
	ReversalAmount          string `json:"reversal_amount"`
	SettlementAmount        string `json:"settlement_amount"`
	ActualAmount            string `json:"actual_amount"`
	OrderStatus             string `json:"order_status"`
	ResponseDescription     string `json:"response_description"`
	VerificationStatus      string `json:"verification_status"`
	OrderTime               string `json:"order_time"`
	ActualCurrency          string `json:"actual_currency"`
	OrderID                 string `json:"order_id"`
	ParentOrderID           string `json:"parent_order_id"`
	MerchantData            string `json:"merchant_data"`
	TranType                string `json:"tran_type"`
	Eci                     string `json:"eci"`
	SettlementDate          string `json:"settlement_date"`
	PaymentSystem           string `json:"payment_system"`
	Rectoken                string `json:"rectoken"`
	ApprovalCode            string `json:"approval_code"`
	MerchantID              int    `json:"merchant_id"`
	SettlementCurrency      string `json:"settlement_currency"`
	PaymentID               int    `json:"payment_id"`
	ProductID               string `json:"product_id"`
	Currency                string `json:"currency"`
	CardBin                 int    `json:"card_bin"`
	ResponseCode            string `json:"response_code"`
	CardType                string `json:"card_type"`
	Amount                  string `json:"amount"`
	SenderEmail             string `json:"sender_email"`
	Signature               string `json:"signature"`
}

type FondyCreate struct {
	Request struct {
		Amount            string `json:"amount"`
		Currency          string `json:"currency"`
		MerchantID        string `json:"merchant_id"`
		OrderDesc         string `json:"order_desc"`
		ProductID         string `json:"product_id"`
		ServerCallbackURL string `json:"server_callback_url"`
		SenderEmail       string `json:"sender_email"`
		Signature         string `json:"signature"`
	} `json:"request"`
}

type FondyResponse struct {
	Response struct {
		PaymentID      string `json:"payment_id"`
		ResponseStatus string `json:"response_status"`
		CheckoutURL    string `json:"checkout_url"`
		ErrorCode      int    `json:"error_code"`
	} `json:"response"`
}

type CreateRequest struct {
	Amount float64 `json:"amount"`
}

type AuthRequest struct {
	Token string `json:"token"`
}

type AuthResponse struct {
	Type      string `json:"type"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Iat       int    `json:"iat"`
	Exp       int    `json:"exp"`
	Sub       string `json:"sub"`
}

func (scheme *FondyRefill) ValidateSignature() error {
	response := strings.Replace(scheme.ResponseSignatureString, "**********", token, 1)
	hash := sha1.New()
	hash.Write([]byte(response))
	if hex.EncodeToString(hash.Sum(nil)) == scheme.Signature {
		return nil
	}
	return errors.New("Токен невалиден")
}

func (scheme *FondyCreate) CreateSignature() (string, error) {
	response := token + "|" + scheme.Request.Amount + "|" + scheme.Request.Currency + "|" + scheme.Request.MerchantID +
		"|" + scheme.Request.OrderDesc + "|" + scheme.Request.ProductID + "|" + scheme.Request.SenderEmail + "|" + scheme.Request.ServerCallbackURL
	hash := sha1.New()
	hash.Write([]byte(response))
	return hex.EncodeToString(hash.Sum(nil)), nil
}
