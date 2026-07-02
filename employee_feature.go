package paymentclient

import (
	"context"

	paymentv1 "github.com/devbefingerspot/fio-paymentservice-go-client/internal/pb/payment/v1"
)

// SyncFeatureEntry is a public type representing a single sync entry for
// employee feature assignment.
type SyncFeatureEntry struct {
	EmployeeID      string
	DesiredFeatures []string
}

// SyncFeatureStatus mirrors the SyncStatus proto enum.
type SyncFeatureStatus int32

const (
	SyncStatusUnspecified SyncFeatureStatus = 0
	SyncStatusAssigned    SyncFeatureStatus = 1
	SyncStatusRevoked     SyncFeatureStatus = 2
	SyncStatusUnchanged   SyncFeatureStatus = 3
	SyncStatusError       SyncFeatureStatus = 4
)

// FeatureSyncOutcome is the public result for a single feature within a sync.
type FeatureSyncOutcome struct {
	ItemKey string
	Status  SyncFeatureStatus
	Error   string
}

// EmployeeSyncOutcome is the public result for a single employee after sync.
type EmployeeSyncOutcome struct {
	EmployeeID     string
	FeatureResults []FeatureSyncOutcome
}

// SyncEmployeeFeaturesOutcome is the public response wrapper for
// SyncEmployeeFeatures, avoiding exposure of internal proto types.
type SyncEmployeeFeaturesOutcome struct {
	Results []EmployeeSyncOutcome
}

// SyncEmployeeFeatures applies a declarative delta to employee feature
// assignments and returns public outcome types.
func (c *Client) SyncEmployeeFeatures(ctx context.Context, companyID, productCode string, entries []SyncFeatureEntry) (*SyncEmployeeFeaturesOutcome, error) {
	if err := c.ensureDialed(); err != nil {
		return nil, err
	}

	protoEntries := make([]*paymentv1.SyncEntry, 0, len(entries))
	for _, e := range entries {
		desired := e.DesiredFeatures
		if desired == nil {
			desired = []string{}
		}
		protoEntries = append(protoEntries, &paymentv1.SyncEntry{
			EmployeeId:      e.EmployeeID,
			DesiredFeatures: desired,
		})
	}

	resp, err := c.employeeFeatureClient.SyncEmployeeFeatures(
		c.outboundContext(ctx),
		&paymentv1.SyncEmployeeFeaturesRequest{
			CompanyId:   companyID,
			ProductCode: productCode,
			SyncEntries: protoEntries,
		},
	)
	if err != nil {
		return nil, err
	}

	return convertSyncResponse(resp), nil
}

func convertSyncResponse(resp *paymentv1.SyncEmployeeFeaturesResponse) *SyncEmployeeFeaturesOutcome {
	if resp == nil {
		return &SyncEmployeeFeaturesOutcome{}
	}
	results := make([]EmployeeSyncOutcome, 0, len(resp.Results))
	for _, r := range resp.Results {
		frs := make([]FeatureSyncOutcome, 0, len(r.FeatureResults))
		for _, fr := range r.FeatureResults {
			frs = append(frs, FeatureSyncOutcome{
				ItemKey: fr.ItemKey,
				Status:  SyncFeatureStatus(fr.Status),
				Error:   fr.Error,
			})
		}
		results = append(results, EmployeeSyncOutcome{
			EmployeeID:     r.EmployeeId,
			FeatureResults: frs,
		})
	}
	return &SyncEmployeeFeaturesOutcome{Results: results}
}

// GetEmployeeFeatures returns the features assigned to the requested employees.
// For flag features, presence in the response means the employee has access at
// the company level. For seat features, only features with an active assignment
// to the given employee are included.
func (c *Client) GetEmployeeFeatures(ctx context.Context, companyID, productCode string, employeeIDs []string) (*paymentv1.GetEmployeeFeaturesResponse, error) {
	if err := c.ensureDialed(); err != nil {
		return nil, err
	}
	return c.employeeFeatureClient.GetEmployeeFeatures(
		c.outboundContext(ctx),
		&paymentv1.GetEmployeeFeaturesRequest{
			CompanyId:   companyID,
			ProductCode: productCode,
			EmployeeIds: employeeIDs,
		},
	)
}
