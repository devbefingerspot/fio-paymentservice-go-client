package paymentclient

import (
	"context"
	"time"

	paymentv1 "github.com/devbefingerspot/fio-paymentservice-go-client/internal/pb/payment/v1"
)

// OwnerQuotaStatus mirrors the OwnerQuotaStatus proto enum.
type OwnerQuotaStatus int32

const (
	OwnerQuotaStatusUnspecified      OwnerQuotaStatus = 0
	OwnerQuotaStatusAssigned         OwnerQuotaStatus = 1
	OwnerQuotaStatusReleased         OwnerQuotaStatus = 2
	OwnerQuotaStatusAlreadyAssigned  OwnerQuotaStatus = 3
	OwnerQuotaStatusNotFound         OwnerQuotaStatus = 4
	OwnerQuotaStatusCapacityExceeded OwnerQuotaStatus = 5
	OwnerQuotaStatusNotEntitled      OwnerQuotaStatus = 6
	OwnerQuotaStatusError            OwnerQuotaStatus = 7
)

// OwnerAssignment represents a single owner holding an app-owner seat.
type OwnerAssignment struct {
	OwnerID    string
	AssignedAt time.Time
	ExpiresAt  time.Time
}

// OwnerQuotaInfo is the public result of GetOwnerQuota.
type OwnerQuotaInfo struct {
	ItemKey        string
	Entitled       bool
	TotalLimit     *int32
	UsedCount      int32
	AvailableCount *int32
	AssignedOwners []OwnerAssignment
}

// OwnerQuotaResult is the public result of AssignOwnerQuota / ReleaseOwnerQuota.
type OwnerQuotaResult struct {
	Status  OwnerQuotaStatus
	Message string
}

// GetOwnerQuota returns the app-owner seat entitlement of a company, including
// the owners currently holding a seat.
func (c *Client) GetOwnerQuota(ctx context.Context, companyID, productCode, itemKey string) (*OwnerQuotaInfo, error) {
	if err := c.ensureDialed(); err != nil {
		return nil, err
	}

	resp, err := c.ownerFeatureClient.GetOwnerQuota(
		c.outboundContext(ctx),
		&paymentv1.GetOwnerQuotaRequest{
			CompanyId:   companyID,
			ProductCode: productCode,
			ItemKey:     itemKey,
		},
	)
	if err != nil {
		return nil, err
	}

	info := &OwnerQuotaInfo{
		ItemKey:        resp.ItemKey,
		Entitled:       resp.Entitled,
		TotalLimit:     resp.TotalLimit,
		UsedCount:      resp.UsedCount,
		AvailableCount: resp.AvailableCount,
	}

	assigned := make([]OwnerAssignment, 0, len(resp.AssignedOwners))
	for _, a := range resp.AssignedOwners {
		oa := OwnerAssignment{OwnerID: a.OwnerId}
		if a.AssignedAt != nil {
			oa.AssignedAt = a.AssignedAt.AsTime()
		}
		if a.ExpiresAt != nil {
			oa.ExpiresAt = a.ExpiresAt.AsTime()
		}
		assigned = append(assigned, oa)
	}
	info.AssignedOwners = assigned

	return info, nil
}

// AssignOwnerQuota assigns an app-owner seat to an owner user.
func (c *Client) AssignOwnerQuota(ctx context.Context, companyID, productCode, itemKey, ownerID string) (*OwnerQuotaResult, error) {
	if err := c.ensureDialed(); err != nil {
		return nil, err
	}

	resp, err := c.ownerFeatureClient.AssignOwnerQuota(
		c.outboundContext(ctx),
		&paymentv1.AssignOwnerQuotaRequest{
			CompanyId:   companyID,
			ProductCode: productCode,
			ItemKey:     itemKey,
			OwnerId:     ownerID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &OwnerQuotaResult{
		Status:  OwnerQuotaStatus(resp.Status),
		Message: resp.Message,
	}, nil
}

// ReleaseOwnerQuota revokes the app-owner seat held by an owner user.
func (c *Client) ReleaseOwnerQuota(ctx context.Context, companyID, productCode, itemKey, ownerID string) (*OwnerQuotaResult, error) {
	if err := c.ensureDialed(); err != nil {
		return nil, err
	}

	resp, err := c.ownerFeatureClient.ReleaseOwnerQuota(
		c.outboundContext(ctx),
		&paymentv1.ReleaseOwnerQuotaRequest{
			CompanyId:   companyID,
			ProductCode: productCode,
			ItemKey:     itemKey,
			OwnerId:     ownerID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &OwnerQuotaResult{
		Status:  OwnerQuotaStatus(resp.Status),
		Message: resp.Message,
	}, nil
}
