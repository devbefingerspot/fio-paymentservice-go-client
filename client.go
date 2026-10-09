// Package paymentclient provides a gRPC client for the Fingerspot payment backend.
//
// Usage:
//
//	client := paymentclient.New("payment-grpc.example.com:50051", "my-api-key")
//	catalog, err := client.GetCatalog(ctx, "fio-web")
//	sub, err := client.GetCompanySubscription(ctx, "company_123", "fio-web")
//	features, err := client.GetEmployeeFeatures(ctx, "company_123", "fio-web", []string{"emp_1"})
package paymentclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"

	paymentv1 "github.com/devbefingerspot/fio-paymentservice-go-client/internal/pb/payment/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Client is the main client for interacting with the payment backend gRPC API.
type Client struct {
	target string
	apiKey string

	connOnce sync.Once
	conn     *grpc.ClientConn
	connErr  error

	catalogClient             paymentv1.CatalogServiceClient
	companySubscriptionClient paymentv1.CompanySubscriptionServiceClient
	employeeFeatureClient     paymentv1.EmployeeFeatureServiceClient
	deviceFeatureClient       paymentv1.DeviceFeatureServiceClient
	deviceRegistryClient      paymentv1.DeviceRegistryServiceClient
	masterDeviceClient        paymentv1.MasterDeviceServiceClient
}

// New creates a new payment gRPC client.
//
//   - target: the gRPC server address (e.g. "payment-grpc:50051")
//   - apiKey: the pre-shared API key for authentication
func New(target, apiKey string) *Client {
	return &Client{
		target: target,
		apiKey: apiKey,
	}
}

// WithInsecure disables TLS for the gRPC connection and uses plaintext instead.
// Use this only for local development or when connecting through a proxy that
// already handles encryption (e.g. a service mesh).
//
// Must be called before the first RPC method call.
func (c *Client) WithInsecure() *Client {
	c.connOnce.Do(func() {
		c.conn, c.connErr = grpc.NewClient(
			c.target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if c.connErr != nil {
			return
		}
		c.catalogClient = paymentv1.NewCatalogServiceClient(c.conn)
		c.companySubscriptionClient = paymentv1.NewCompanySubscriptionServiceClient(c.conn)
		c.employeeFeatureClient = paymentv1.NewEmployeeFeatureServiceClient(c.conn)
		c.deviceFeatureClient = paymentv1.NewDeviceFeatureServiceClient(c.conn)
		c.deviceRegistryClient = paymentv1.NewDeviceRegistryServiceClient(c.conn)
		c.masterDeviceClient = paymentv1.NewMasterDeviceServiceClient(c.conn)
	})
	return c
}

// Close releases the cached gRPC connection.
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// dialOnce lazily establishes the gRPC connection on first use (TLS by default).
func (c *Client) dialOnce() (*grpc.ClientConn, error) {
	c.connOnce.Do(func() {
		creds := credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})
		c.conn, c.connErr = grpc.NewClient(
			c.target,
			grpc.WithTransportCredentials(creds),
		)
		if c.connErr != nil {
			return
		}
		c.catalogClient = paymentv1.NewCatalogServiceClient(c.conn)
		c.companySubscriptionClient = paymentv1.NewCompanySubscriptionServiceClient(c.conn)
		c.employeeFeatureClient = paymentv1.NewEmployeeFeatureServiceClient(c.conn)
		c.deviceFeatureClient = paymentv1.NewDeviceFeatureServiceClient(c.conn)
		c.deviceRegistryClient = paymentv1.NewDeviceRegistryServiceClient(c.conn)
		c.masterDeviceClient = paymentv1.NewMasterDeviceServiceClient(c.conn)
	})
	return c.conn, c.connErr
}

// outboundContext returns ctx enriched with the authorization metadata header
// that the payment backend's auth interceptor expects.
func (c *Client) outboundContext(ctx context.Context) context.Context {
	if c.apiKey == "" {
		return ctx
	}
	return metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", c.apiKey))
}

// ensureDialed ensures the connection is established before an RPC call.
func (c *Client) ensureDialed() error {
	_, err := c.dialOnce()
	if err != nil {
		return fmt.Errorf("paymentclient: dial %s: %w", c.target, err)
	}
	return nil
}
