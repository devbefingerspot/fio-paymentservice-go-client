# fio-payment-go-client

Go client library for the Fingerspot payment backend gRPC API.

## Usage

```go
import paymentclient "github.com/devbefingerspot/fio-payment-go-client"

client := paymentclient.New("payment-grpc:50051", os.Getenv("PAYMENT_GRPC_API_KEY"))
defer client.Close()

// For local dev without TLS:
// client := paymentclient.New("localhost:50051", "...").WithInsecure()

catalog, err := client.GetCatalog(ctx, "fio-web")
sub, err := client.GetCompanySubscription(ctx, companyID, "fio-web")
features, err := client.GetEmployeeFeatures(ctx, companyID, "fio-web", employeeIDs)
syncResult, err := client.SyncEmployeeFeatures(ctx, companyID, "fio-web", entries)
```

## Services

| Service | Method | Description |
|---|---|---|
| `CatalogService` | `GetCatalog` | Returns the plan catalog for a product |
| `CompanySubscriptionService` | `GetCompanySubscription` | Returns the active subscription for a company |
| `EmployeeFeatureService` | `GetEmployeeFeatures` | Returns features assigned to employees |
| `EmployeeFeatureService` | `SyncEmployeeFeatures` | Applies declarative delta to employee features |

## Updating Proto Files

When the proto definitions in `team-fingerspot-io-payment-backend/proto/payment/v1/` change:

1. Copy the generated files:
   ```
   cp team-fingerspot-io-payment-backend/src/gen/payment/v1/*.go \
      team-payment-service-go-client/internal/pb/payment/v1/
   ```
2. Add wrapper methods for any new RPCs
3. Run `go mod tidy && go build ./...`
