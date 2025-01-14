package amazonpay

import (
	"context"
	"fmt"
	"net/http"
)

type CreateChargeRequest struct {
	ChargePermissionID            string            `json:"chargePermissionId,omitempty"`
	ChargeAmount                  *Price            `json:"chargeAmount,omitempty"`
	CaptureNow                    *bool             `json:"captureNow,omitempty"`
	SoftDescriptor                string            `json:"softDescriptor,omitempty"`
	CanHandlePendingAuthorization *bool             `json:"canHandlePendingAuthorization,omitempty"`
	MerchantMetadata              *MerchantMetadata `json:"merchantMetadata,omitempty"`
	ProviderMetadata              *ProviderMetadata `json:"providerMetadata,omitempty"`
}

type CreateChargeResponse struct {
	ErrorResponse
	ChargeID            string            `json:"chargeId,omitempty"`
	ChargePermissionID  string            `json:"chargePermissionId,omitempty"`
	ChargeAmount        *Price            `json:"chargeAmount,omitempty"`
	CaptureAmount       *Price            `json:"captureAmount,omitempty"`
	RefundedAmount      *Price            `json:"refundedAmount,omitempty"`
	ConvertedAmount     string            `json:"convertedAmount,omitempty"`
	ConversionRate      string            `json:"conversionRate,omitempty"`
	SoftDescriptor      string            `json:"softDescriptor,omitempty"`
	MerchantMetadata    *MerchantMetadata `json:"merchantMetadata,omitempty"`
	ProviderMetadata    *ProviderMetadata `json:"providerMetadata,omitempty"`
	StatusDetails       *StatusDetails    `json:"statusDetails,omitempty"`
	CreationTimestamp   string            `json:"creationTimestamp,omitempty"`
	ExpirationTimestamp string            `json:"expirationTimestamp,omitempty"`
	ReleaseEnvironment  string            `json:"releaseEnvironment,omitempty"`
}

func (c *Client) CreateCharge(ctx context.Context, req *CreateChargeRequest) (*CreateChargeResponse, *http.Response, error) {
	path := fmt.Sprintf("%s/charges", APIVersion)
	httpReq, err := c.NewRequest(http.MethodPost, path, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(CreateChargeResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

type GetChargeResponse struct {
	ErrorResponse
	ChargeID            string            `json:"chargeId"`
	ChargePermissionID  string            `json:"chargePermissionId"`
	ChargeAmount        *Price            `json:"chargeAmount"`
	CaptureAmount       *Price            `json:"captureAmount"`
	RefundedAmount      *Price            `json:"refundedAmount"`
	ConvertedAmount     string            `json:"convertedAmount"`
	ConversionRate      string            `json:"conversionRate"`
	SoftDescriptor      string            `json:"softDescriptor"`
	MerchantMetadata    *MerchantMetadata `json:"merchantMetadata"`
	ProviderMetadata    *ProviderMetadata `json:"providerMetadata"`
	StatusDetails       *StatusDetails    `json:"statusDetails"`
	CreationTimestamp   string            `json:"creationTimestamp"`
	ExpirationTimestamp string            `json:"expirationTimestamp"`
	ReleaseEnvironment  string            `json:"releaseEnvironment"`
}

func (c *Client) GetCharge(ctx context.Context, chargeID string) (*GetChargeResponse, *http.Response, error) {
	path := fmt.Sprintf("%s/charges/%s", APIVersion, chargeID)
	httpReq, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	resp := new(GetChargeResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

type CaptureChargeRequest struct {
	CaptureAmount    *Price            `json:"captureAmount,omitempty"`
	SoftDescriptor   string            `json:"softDescriptor,omitempty"`
	MerchantMetadata *MerchantMetadata `json:"merchantMetadata,omitempty"`
	ProviderMetadata *ProviderMetadata `json:"providerMetadata,omitempty"`
}

type CaptureChargeResponse struct {
	ErrorResponse
	ChargeID            string            `json:"chargeId,omitempty"`
	ChargePermissionID  string            `json:"chargePermissionId,omitempty"`
	ChargeAmount        *Price            `json:"chargeAmount,omitempty"`
	CaptureAmount       *Price            `json:"captureAmount,omitempty"`
	RefundedAmount      *Price            `json:"refundedAmount,omitempty"`
	ConvertedAmount     string            `json:"convertedAmount,omitempty"`
	ConversionRate      string            `json:"conversionRate,omitempty"`
	SoftDescriptor      string            `json:"softDescriptor,omitempty"`
	MerchantMetadata    *MerchantMetadata `json:"merchantMetadata,omitempty"`
	ProviderMetadata    *ProviderMetadata `json:"providerMetadata,omitempty"`
	StatusDetails       *StatusDetails    `json:"statusDetails,omitempty"`
	CreationTimestamp   string            `json:"creationTimestamp,omitempty"`
	ExpirationTimestamp string            `json:"expirationTimestamp,omitempty"`
	ReleaseEnvironment  string            `json:"releaseEnvironment,omitempty"`
}

func (c *Client) CaptureCharge(ctx context.Context, chargeID string, req *CaptureChargeRequest) (*CaptureChargeResponse, *http.Response, error) {
	path := fmt.Sprintf("%s/charges/%s/capture", APIVersion, chargeID)
	httpReq, err := c.NewRequest(http.MethodPost, path, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(CaptureChargeResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

type CancelChargeRequest struct {
	CancellationReason string `json:"cancellationReason,omitempty"`
}

type CancelChargeResponse struct {
	ErrorResponse
	ChargeID            string            `json:"chargeId,omitempty"`
	ChargePermissionID  string            `json:"chargePermissionId,omitempty"`
	ChargeAmount        *Price            `json:"chargeAmount,omitempty"`
	CaptureAmount       *Price            `json:"captureAmount,omitempty"`
	RefundedAmount      *Price            `json:"refundedAmount,omitempty"`
	ConvertedAmount     string            `json:"convertedAmount,omitempty"`
	ConversionRate      string            `json:"conversionRate,omitempty"`
	SoftDescriptor      string            `json:"softDescriptor,omitempty"`
	MerchantMetadata    *MerchantMetadata `json:"merchantMetadata,omitempty"`
	ProviderMetadata    *ProviderMetadata `json:"providerMetadata,omitempty"`
	StatusDetails       *StatusDetails    `json:"statusDetails,omitempty"`
	CreationTimestamp   string            `json:"creationTimestamp,omitempty"`
	ExpirationTimestamp string            `json:"expirationTimestamp,omitempty"`
	ReleaseEnvironment  string            `json:"releaseEnvironment,omitempty"`
}

func (c *Client) CancelCharge(ctx context.Context, chargeID string, req *CancelChargeRequest) (*CancelChargeResponse, *http.Response, error) {
	path := fmt.Sprintf("%s/charges/%s/cancel", APIVersion, chargeID)
	httpReq, err := c.NewRequest(http.MethodDelete, path, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(CancelChargeResponse)
	httpResp, err := c.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}
