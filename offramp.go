package fonbnk

import (
	"context"
	"net/http"

	"github.com/google/go-querystring/query"
)

type (
	OrderLimitsQuery struct {
		Type    OffRampType `url:"type"`
		Country Country     `url:"country"`
	}

	OrderLimitsResponse struct {
		MinUSD           float64
		MaxUSD           uint64
		MinLocalCurrency uint64
		MaxLocalCurrency uint64
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
