package paymentclient

import (
	"context"

	paymentv1 "github.com/devbefingerspot/fio-payment-go-client/internal/pb/payment/v1"
)

// GetCatalog returns the full subscription plan catalog for a product,
// including business plan metadata, plans with their feature compositions.
func (c *Client) GetCatalog(ctx context.Context, productCode string) (*paymentv1.GetCatalogResponse, error) {
	if err := c.ensureDialed(); err != nil {
		return nil, err
	}
	return c.catalogClient.GetCatalog(
		c.outboundContext(ctx),
		&paymentv1.GetCatalogRequest{ProductCode: productCode},
	)
}
