package raw

import (
	"context"

	"google.golang.org/grpc"

	billingv1 "github.com/modelrelay/xai-go/gen/xai/management_api/v1"
	"github.com/modelrelay/xai-go/gen/xai/shared/analytics"
)

// BillingClient wraps the billing UISvc service.
type BillingClient struct {
	stub billingv1.UISvcClient
}

// NewBillingClient creates a new billing client.
func NewBillingClient(stub billingv1.UISvcClient) *BillingClient {
	return &BillingClient{stub: stub}
}

func (c *BillingClient) SetBillingInfo(ctx context.Context, req *billingv1.SetBillingInfoReq, opts ...grpc.CallOption) (*billingv1.SetBillingInfoResp, error) {
	return c.stub.SetBillingInfo(ctx, req, opts...)
}

func (c *BillingClient) GetBillingInfo(ctx context.Context, req *billingv1.GetBillingInfoReq, opts ...grpc.CallOption) (*billingv1.GetBillingInfoResp, error) {
	return c.stub.GetBillingInfo(ctx, req, opts...)
}

func (c *BillingClient) ListPaymentMethods(ctx context.Context, req *billingv1.ListPaymentMethodsReq, opts ...grpc.CallOption) (*billingv1.ListPaymentMethodsResp, error) {
	return c.stub.ListPaymentMethods(ctx, req, opts...)
}

func (c *BillingClient) SetDefaultPaymentMethod(ctx context.Context, req *billingv1.SetDefaultPaymentMethodReq, opts ...grpc.CallOption) (*billingv1.SetDefaultPaymentMethodResp, error) {
	return c.stub.SetDefaultPaymentMethod(ctx, req, opts...)
}

func (c *BillingClient) GetAmountToPay(ctx context.Context, req *billingv1.GetAmountToPayReq, opts ...grpc.CallOption) (*billingv1.GetAmountToPayResp, error) {
	return c.stub.GetAmountToPay(ctx, req, opts...)
}

func (c *BillingClient) AnalyzeBillingItems(ctx context.Context, req *billingv1.AnalyzeBillingItemsRequest, opts ...grpc.CallOption) (*analytics.AnalyticsResponse, error) {
	return c.stub.AnalyzeBillingItems(ctx, req, opts...)
}

func (c *BillingClient) ListInvoices(ctx context.Context, req *billingv1.ListInvoicesReq, opts ...grpc.CallOption) (*billingv1.ListInvoicesResp, error) {
	return c.stub.ListInvoices(ctx, req, opts...)
}

func (c *BillingClient) ListPrepaidBalanceChanges(ctx context.Context, req *billingv1.ListPrepaidBalanceChangesReq, opts ...grpc.CallOption) (*billingv1.ListPrepaidBalanceChangesResp, error) {
	return c.stub.ListPrepaidBalanceChanges(ctx, req, opts...)
}

func (c *BillingClient) TopUpOrGetExistingPendingChange(ctx context.Context, req *billingv1.TopUpOrGetExistingPendingChangeReq, opts ...grpc.CallOption) (*billingv1.TopUpOrGetExistingPendingChangeResp, error) {
	return c.stub.TopUpOrGetExistingPendingChange(ctx, req, opts...)
}

func (c *BillingClient) GetSpendingLimits(ctx context.Context, req *billingv1.GetSpendingLimitsReq, opts ...grpc.CallOption) (*billingv1.GetSpendingLimitsResp, error) {
	return c.stub.GetSpendingLimits(ctx, req, opts...)
}

func (c *BillingClient) SetSoftSpendingLimit(ctx context.Context, req *billingv1.SetSoftSpendingLimitReq, opts ...grpc.CallOption) (*billingv1.SetSoftSpendingLimitResp, error) {
	return c.stub.SetSoftSpendingLimit(ctx, req, opts...)
}
