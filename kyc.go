package fonbnk

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type (
	KYCStateResponse struct {
		KYCURL               string    `json:"kycUrl"`
		OfframpKYCRequired   bool      `json:"offrampKycRequired"`
		OnrampKYCRequired    bool      `json:"onrampKycRequired"`
		PassedKYC            bool      `json:"passedKyc"`
		KYCStatus            KYCStatus `json:"kycStatus"`
		KYCStatusDescription string    `json:"kycStatusDescription"`
		ReachedKYCLimit      bool      `json:"reachedKycLimit"`
	}

	KYCSubmitBody struct {
		PhoneNumber string    `json:"phoneNumber"`
		IDType      KYCIDType `json:"idType"`
		UserFields  struct {
			FirstName string    `json:"first_name"`
			LastName  string    `json:"last_name"`
			Dob       time.Time `json:"dob"`
			Email     string    `json:"email"`
			IDNumber  string    `json:"id_number"`
		} `json:"userFields"`
	}

	KYCSubmitResponse struct {
		Success bool `json:"success"`
	}
)

const kycPath = "/api/kyc/"

func (fc *FonbnkClient) KYCState(ctx context.Context, phoneNumber string) (KYCStateResponse, error) {
	kycStateResp := KYCStateResponse{}

	resp, err := fc.requestWithCtx(ctx, http.MethodGet, fc.endpoint+kycPath+"state?phoneNumber="+phoneNumber, nil)
	if err != nil {
		return kycStateResp, err
	}

	if err := parseResponse(resp, &kycStateResp); err != nil {
		return kycStateResp, err
	}

	return kycStateResp, nil
}

func (fc *FonbnkClient) KYCSubmit(ctx context.Context, input KYCSubmitBody) (KYCSubmitResponse, error) {
	kycSubmitResp := KYCSubmitResponse{}

	jsonRequestBody, err := json.Marshal(&input)
	if err != nil {
		return kycSubmitResp, err
	}

	resp, err := fc.requestWithCtx(ctx, http.MethodPost, fc.endpoint+kycPath+"submit", bytes.NewBuffer(jsonRequestBody))
	if err != nil {
		return kycSubmitResp, err
	}

	if err := parseResponse(resp, &kycSubmitResp); err != nil {
		return kycSubmitResp, err
	}

	return kycSubmitResp, nil
}
