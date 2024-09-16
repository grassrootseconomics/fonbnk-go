package fonbnk

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

type (
	ValidateOption func(b *ValidateBody)

	ValidateBody struct {
		OfferID        string            `json:"offerId"`
		RequiredFields map[string]string `json:"requiredFields"`
	}

	CreateOrderBody struct {
		OfferID        string             `json:"offerId"`
		RequiredFields map[string]string  `json:"requiredFields"`
		PaymentType    OffRampPaymentType `json:"paymentType"`
		Network        Network            `json:"network"`
		Asset          OffRampAsset       `json:"asset"`
		Address        string             `json:"address"`
		Currency       OffRampCurrency    `json:"currency"`
		Amount         float64            `json:"amount"`
		IP             string             `json:"ip"`
		OrderParams    string             `json:"orderParams"`
	}

	ConfirmOrderBody struct {
		OrderID string `json:"orderId"`
		Hash    string `json:"hash"`
	}

	OrderLimitsQuery struct {
		Type    OffRampType `url:"type"`
		Country Country     `url:"country"`
	}

	BestOfferQuery struct {
		Amount   float64         `url:"amount"`
		Currency OffRampCurrency `url:"currency"`
		Type     OffRampType     `url:"type"`
		Country  Country         `url:"country"`
	}

	OrderLimitsResponse struct {
		MinUSD           float64 `json:"minUsd"`
		MaxUSD           float64 `json:"maxUsd"`
		MinLocalCurrency float64 `json:"minLocalCurrency"`
		MaxLocalCurrency float64 `json:"maxLocalCurrency"`
	}

	// BestOfferResponse is a minimized reponse struct with enough information to proceed with off-ramping.
	BestOfferResponse struct {
		Offer struct {
			ID           string  `json:"_id"`
			ExchangeRate float64 `json:"exchangeRate"`
		} `json:"offer"`
		Cashout struct {
			LocalCurrencyAmount           float64 `json:"localCurrencyAmount"`
			UsdAmount                     float64 `json:"usdAmount"`
			FeeAmountUsd                  float64 `json:"feeAmountUsd"`
			FeeAmountUsdFonbnk            float64 `json:"feeAmountUsdFonbnk"`
			FeeAmountUsdPartner           float64 `json:"feeAmountUsdPartner"`
			FeeAmountLocalCurrency        float64 `json:"feeAmountLocalCurrency"`
			FeeAmountLocalCurrencyFonbnk  float64 `json:"feeAmountLocalCurrencyFonbnk"`
			FeeAmountLocalCurrencyPartner float64 `json:"feeAmountLocalCurrencyPartner"`
		} `json:"cashout"`
	}

	ValidateResponse struct {
		Details []struct {
			Label string `json:"label"`
			Value string `json:"value"`
		} `json:"details"`
	}

	OrderResponse struct {
		ID           string             `json:"_id"`
		OfferID      string             `json:"offerId"`
		PaymentType  OffRampPaymentType `json:"paymentType"`
		Network      Network            `json:"network"`
		Asset        OffRampAsset       `json:"asset"`
		ExchangeRate float64            `json:"exchangeRate"`
		FromAddress  string             `json:"fromAddress"`
		ToAddress    string             `json:"toAddress"`
		Status       string             `json:"status"`
		CreatedAt    time.Time          `json:"createdAt"`
		ExpiresAt    time.Time          `json:"expiresAt"`
		Hash         string             `json:"hash"`
		OrderParams  string             `json:"orderParams"`
		Cashout      struct {
			LocalCurrencyAmount           float64 `json:"localCurrencyAmount"`
			UsdAmount                     float64 `json:"usdAmount"`
			FeeAmountUsd                  float64 `json:"feeAmountUsd"`
			FeeAmountUsdFonbnk            float64 `json:"feeAmountUsdFonbnk"`
			FeeAmountUsdPartner           float64 `json:"feeAmountUsdPartner"`
			FeeAmountLocalCurrency        float64 `json:"feeAmountLocalCurrency"`
			FeeAmountLocalCurrencyFonbnk  float64 `json:"feeAmountLocalCurrencyFonbnk"`
			FeeAmountLocalCurrencyPartner float64 `json:"feeAmountLocalCurrencyPartner"`
		} `json:"cashout"`
	}
)

const offRampPath = "/api/offramp/"

func (fc *FonbnkClient) OrderLimits(ctx context.Context, input OrderLimitsQuery) (OrderLimitsResponse, error) {
	orderLimitsResp := OrderLimitsResponse{}

	v, err := query.Values(input)
	if err != nil {
		return orderLimitsResp, err
	}

	resp, err := fc.requestWithCtx(ctx, http.MethodGet, fc.endpoint+offRampPath+"limits?"+v.Encode(), nil)
	if err != nil {
		return orderLimitsResp, err
	}

	if err := parseResponse(resp, &orderLimitsResp); err != nil {
		return orderLimitsResp, err
	}

	return orderLimitsResp, nil
}

func (fc *FonbnkClient) BestOffer(ctx context.Context, input BestOfferQuery) (BestOfferResponse, error) {
	bestOfferResp := BestOfferResponse{}

	v, err := query.Values(input)
	if err != nil {
		return bestOfferResp, err
	}

	resp, err := fc.requestWithCtx(ctx, http.MethodGet, fc.endpoint+offRampPath+"best-offer?"+v.Encode(), nil)
	if err != nil {
		return bestOfferResp, err
	}

	if err := parseResponse(resp, &bestOfferResp); err != nil {
		return bestOfferResp, err
	}

	return bestOfferResp, nil
}

func (fc *FonbnkClient) Validate(ctx context.Context, input ValidateBody, validateOption ValidateOption) (ValidateResponse, error) {
	validateResp := ValidateResponse{}
	validateOption(&input)

	jsonRequestBody, err := json.Marshal(&input)
	if err != nil {
		return validateResp, err
	}

	resp, err := fc.requestWithCtx(ctx, http.MethodPost, fc.endpoint+offRampPath+"validate-fields", bytes.NewBuffer(jsonRequestBody))
	if err != nil {
		return validateResp, err
	}

	if err := parseResponse(resp, &validateResp); err != nil {
		return validateResp, err
	}

	return validateResp, nil
}

func (fc *FonbnkClient) CreateOrder(ctx context.Context, input CreateOrderBody, validateInput ValidateBody, validateOption ValidateOption) (OrderResponse, error) {
	createOrderResp := OrderResponse{}
	validateOption(&validateInput)
	input.RequiredFields = validateInput.RequiredFields

	jsonRequestBody, err := json.Marshal(&input)
	if err != nil {
		return createOrderResp, err
	}

	log.Printf("%+v\n", string(jsonRequestBody))

	resp, err := fc.requestWithCtx(ctx, http.MethodPost, fc.endpoint+offRampPath+"create-order", bytes.NewBuffer(jsonRequestBody))
	if err != nil {
		return createOrderResp, err
	}

	if err := parseResponse(resp, &createOrderResp); err != nil {
		return createOrderResp, err
	}

	return createOrderResp, nil
}

func (fc *FonbnkClient) ConfirmOrder(ctx context.Context, input ConfirmOrderBody) (OrderResponse, error) {
	confirmOrderResp := OrderResponse{}

	jsonRequestBody, err := json.Marshal(&input)
	if err != nil {
		return confirmOrderResp, err
	}

	log.Printf("%+v\n", string(jsonRequestBody))

	resp, err := fc.requestWithCtx(ctx, http.MethodPost, fc.endpoint+offRampPath+"confirm-order", bytes.NewBuffer(jsonRequestBody))
	if err != nil {
		return confirmOrderResp, err
	}

	if err := parseResponse(resp, &confirmOrderResp); err != nil {
		return confirmOrderResp, err
	}

	return confirmOrderResp, nil
}

func ValidatePhoneNumber(phoneNumber string) ValidateOption {
	return func(b *ValidateBody) {
		b.RequiredFields = map[string]string{
			"phoneNumber": phoneNumber,
		}
	}
}

func ValidatePayBill(shortCode string, accountNumber string) ValidateOption {
	return func(b *ValidateBody) {
		b.RequiredFields = map[string]string{
			"shortCode":     shortCode,
			"accountNumber": accountNumber,
		}
	}
}
