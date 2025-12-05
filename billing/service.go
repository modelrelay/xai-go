package billing

import (
	"context"

	"google.golang.org/grpc"

	"github.com/modelrelay/xai-go/gen/xai/management_api/v1"
	"github.com/modelrelay/xai-go/gen/xai/shared/analytics"
	"github.com/modelrelay/xai-go/internal/raw"
)

// Service exposes billing management RPCs.
type Service struct {
	raw *raw.BillingClient
}

// NewService constructs a billing service.
func NewService(rawClient *raw.BillingClient) Service {
	return Service{raw: rawClient}
}

func (s Service) SetBillingInfo(ctx context.Context, req *v1.SetBillingInfoReq, opts ...grpc.CallOption) (*v1.SetBillingInfoResp, error) {
	return s.raw.SetBillingInfo(ctx, req, opts...)
}

func (s Service) GetBillingInfo(ctx context.Context, req *v1.GetBillingInfoReq, opts ...grpc.CallOption) (*v1.GetBillingInfoResp, error) {
	return s.raw.GetBillingInfo(ctx, req, opts...)
}

func (s Service) ListPaymentMethods(ctx context.Context, req *v1.ListPaymentMethodsReq, opts ...grpc.CallOption) (*v1.ListPaymentMethodsResp, error) {
	return s.raw.ListPaymentMethods(ctx, req, opts...)
}

func (s Service) SetDefaultPaymentMethod(ctx context.Context, req *v1.SetDefaultPaymentMethodReq, opts ...grpc.CallOption) (*v1.SetDefaultPaymentMethodResp, error) {
	return s.raw.SetDefaultPaymentMethod(ctx, req, opts...)
}

func (s Service) GetAmountToPay(ctx context.Context, req *v1.GetAmountToPayReq, opts ...grpc.CallOption) (*v1.GetAmountToPayResp, error) {
	return s.raw.GetAmountToPay(ctx, req, opts...)
}

func (s Service) AnalyzeBillingItems(ctx context.Context, req *v1.AnalyzeBillingItemsRequest, opts ...grpc.CallOption) (*analytics.AnalyticsResponse, error) {
	return s.raw.AnalyzeBillingItems(ctx, req, opts...)
}

func (s Service) ListInvoices(ctx context.Context, req *v1.ListInvoicesReq, opts ...grpc.CallOption) (*v1.ListInvoicesResp, error) {
	return s.raw.ListInvoices(ctx, req, opts...)
}

func (s Service) ListPrepaidBalanceChanges(ctx context.Context, req *v1.ListPrepaidBalanceChangesReq, opts ...grpc.CallOption) (*v1.ListPrepaidBalanceChangesResp, error) {
	return s.raw.ListPrepaidBalanceChanges(ctx, req, opts...)
}

func (s Service) TopUpOrGetExistingPendingChange(ctx context.Context, req *v1.TopUpOrGetExistingPendingChangeReq, opts ...grpc.CallOption) (*v1.TopUpOrGetExistingPendingChangeResp, error) {
	return s.raw.TopUpOrGetExistingPendingChange(ctx, req, opts...)
}

func (s Service) GetSpendingLimits(ctx context.Context, req *v1.GetSpendingLimitsReq, opts ...grpc.CallOption) (*v1.GetSpendingLimitsResp, error) {
	return s.raw.GetSpendingLimits(ctx, req, opts...)
}

func (s Service) SetSoftSpendingLimit(ctx context.Context, req *v1.SetSoftSpendingLimitReq, opts ...grpc.CallOption) (*v1.SetSoftSpendingLimitResp, error) {
	return s.raw.SetSoftSpendingLimit(ctx, req, opts...)
}
