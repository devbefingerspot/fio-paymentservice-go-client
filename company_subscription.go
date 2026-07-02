package paymentclient

import (
	"context"

	paymentv1 "github.com/devbefingerspot/fio-paymentservice-go-client/internal/pb/payment/v1"
)

// GetCompanySubscription returns the active subscription for a company on a
// product, including full plan, feature, and addon details.
// Returns nil subscription when no active subscription exists.
func (c *Client) GetCompanySubscription(ctx context.Context, companyID, productCode string) (*paymentv1.CompanySubscription, error) {
	if err := c.ensureDialed(); err != nil {
		return nil, err
	}
	resp, err := c.companySubscriptionClient.GetCompanySubscription(
		c.outboundContext(ctx),
		&paymentv1.GetCompanySubscriptionRequest{
			CompanyId:   companyID,
			ProductCode: productCode,
		},
	)
	if err != nil {
		return nil, err
	}
	return resp.GetSubscription(), nil
}
