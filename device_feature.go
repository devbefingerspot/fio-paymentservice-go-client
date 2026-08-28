package paymentclient

import (
	"context"
	"time"

	paymentv1 "github.com/devbefingerspot/fio-paymentservice-go-client/internal/pb/payment/v1"
)

// DeviceQuotaStatus mirrors the DeviceQuotaStatus proto enum.
type DeviceQuotaStatus int32

const (
	DeviceQuotaStatusUnspecified      DeviceQuotaStatus = 0
	DeviceQuotaStatusAssigned         DeviceQuotaStatus = 1
	DeviceQuotaStatusReleased         DeviceQuotaStatus = 2
	DeviceQuotaStatusAlreadyAssigned  DeviceQuotaStatus = 3
	DeviceQuotaStatusNotFound         DeviceQuotaStatus = 4
	DeviceQuotaStatusCapacityExceeded DeviceQuotaStatus = 5
	DeviceQuotaStatusError            DeviceQuotaStatus = 6
)

// DeviceAssignment represents a single device that holds an active seat.
type DeviceAssignment struct {
	CloudID    string
	AssignedAt time.Time
	ExpiresAt  time.Time
}

// DeviceQuotaInfo is the public result of GetDeviceQuota.
type DeviceQuotaInfo struct {
	ItemKey         string
	TotalLimit      *int32
	UsedCount       int32
	AvailableCount  *int32
	AssignedDevices []DeviceAssignment
}

// DeviceQuotaResult is the public result of AssignDeviceQuota / ReleaseDeviceQuota.
type DeviceQuotaResult struct {
	Status  DeviceQuotaStatus
	Message string
}

// GetDeviceQuota returns quota usage for a device seat feature, including the
// list of currently assigned devices.
func (c *Client) GetDeviceQuota(ctx context.Context, companyID, productCode, itemKey string) (*DeviceQuotaInfo, error) {
	if err := c.ensureDialed(); err != nil {
		return nil, err
	}

	resp, err := c.deviceFeatureClient.GetDeviceQuota(
		c.outboundContext(ctx),
		&paymentv1.GetDeviceQuotaRequest{
			CompanyId:   companyID,
			ProductCode: productCode,
			ItemKey:     itemKey,
		},
	)
	if err != nil {
		return nil, err
	}

	info := &DeviceQuotaInfo{
		ItemKey:    resp.ItemKey,
		UsedCount:  resp.UsedCount,
		TotalLimit: resp.TotalLimit,
	}

	if resp.AvailableCount != nil {
		info.AvailableCount = resp.AvailableCount
	}

	assigned := make([]DeviceAssignment, 0, len(resp.AssignedDevices))
	for _, a := range resp.AssignedDevices {
		da := DeviceAssignment{
			CloudID: a.CloudId,
		}
		if a.AssignedAt != nil {
			da.AssignedAt = a.AssignedAt.AsTime()
		}
		if a.ExpiresAt != nil {
			da.ExpiresAt = a.ExpiresAt.AsTime()
		}
		assigned = append(assigned, da)
	}
	info.AssignedDevices = assigned

	return info, nil
}

// AssignDeviceQuota assigns a device seat to a specific device identified by
// its cloud ID.
func (c *Client) AssignDeviceQuota(ctx context.Context, companyID, productCode, itemKey, cloudID string) (*DeviceQuotaResult, error) {
	if err := c.ensureDialed(); err != nil {
		return nil, err
	}

	resp, err := c.deviceFeatureClient.AssignDeviceQuota(
		c.outboundContext(ctx),
		&paymentv1.AssignDeviceQuotaRequest{
			CompanyId:   companyID,
			ProductCode: productCode,
			ItemKey:     itemKey,
			CloudId:     cloudID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &DeviceQuotaResult{
		Status:  DeviceQuotaStatus(resp.Status),
		Message: resp.Message,
	}, nil
}

// ReleaseDeviceQuota revokes a device seat assignment from a device identified
// by its cloud ID.
func (c *Client) ReleaseDeviceQuota(ctx context.Context, companyID, productCode, itemKey, cloudID string) (*DeviceQuotaResult, error) {
	if err := c.ensureDialed(); err != nil {
		return nil, err
	}

	resp, err := c.deviceFeatureClient.ReleaseDeviceQuota(
		c.outboundContext(ctx),
		&paymentv1.ReleaseDeviceQuotaRequest{
			CompanyId:   companyID,
			ProductCode: productCode,
			ItemKey:     itemKey,
			CloudId:     cloudID,
		},
	)
	if err != nil {
		return nil, err
	}

	return &DeviceQuotaResult{
		Status:  DeviceQuotaStatus(resp.Status),
		Message: resp.Message,
	}, nil
}
