package paymentclient

import (
	"context"

	paymentv1 "github.com/devbefingerspot/fio-paymentservice-go-client/internal/pb/payment/v1"
)

// DeviceData mirrors the DeviceData proto message for device registry operations.
type DeviceData struct {
	ID           string
	Label        string
	CompanyID    string
	SerialNumber string
	CloudID      string
	IPAddress    string
	PortNumber   int32
	Username     string
	Password     string
	Brand        string
	Model        string
	DeviceTypeID int64
	Timezone     string
	ShortLabel   string
	Description  string
	Notes        string
	IsActive     bool
}

// DeviceLookupResult is returned by LookupDevice.
type DeviceLookupResult struct {
	Found  bool
	Device *DeviceData
}

// DeviceCreateResult is returned by CreateDevice.
type DeviceCreateResult struct {
	ID                string
	ConflictCompanyID string
}

// LookupDevice finds a device by serial_number or cloud_id in the payment registry.
func (c *Client) LookupDevice(ctx context.Context, serialNumber, cloudID string) (*DeviceLookupResult, error) {
	if err := c.ensureDialed(); err != nil {
		return nil, err
	}

	resp, err := c.deviceRegistryClient.LookupDevice(
		c.outboundContext(ctx),
		&paymentv1.LookupDeviceRequest{
			SerialNumber: serialNumber,
			CloudId:      cloudID,
		},
	)
	if err != nil {
		return nil, err
	}

	result := &DeviceLookupResult{Found: resp.Found}
	if resp.Device != nil {
		result.Device = protoDeviceToPublic(resp.Device)
	}
	return result, nil
}

// CreateDevice inserts a new device into the payment registry.
func (c *Client) CreateDevice(ctx context.Context, d *DeviceData) (*DeviceCreateResult, error) {
	if err := c.ensureDialed(); err != nil {
		return nil, err
	}

	resp, err := c.deviceRegistryClient.CreateDevice(
		c.outboundContext(ctx),
		&paymentv1.CreateDeviceRequest{Device: publicDeviceToProto(d)},
	)
	if err != nil {
		return nil, err
	}

	return &DeviceCreateResult{
		ID:                resp.Id,
		ConflictCompanyID: resp.ConflictCompanyId,
	}, nil
}

// UpdateDevice syncs mutable fields for a device in the payment registry.
func (c *Client) UpdateDevice(ctx context.Context, id string, req *UpdateDeviceFields) error {
	if err := c.ensureDialed(); err != nil {
		return err
	}

	protoReq := &paymentv1.UpdateDeviceRequest{Id: id}
	if req != nil {
		protoReq.Label = req.Label
		protoReq.IpAddress = req.IPAddress
		protoReq.PortNumber = req.PortNumber
		protoReq.Username = req.Username
		protoReq.Password = req.Password
		protoReq.Brand = req.Brand
		protoReq.Timezone = req.Timezone
		protoReq.DeviceTypeId = req.DeviceTypeID
		protoReq.ShortLabel = req.ShortLabel
		protoReq.Description = req.Description
		protoReq.Notes = req.Notes
		protoReq.IsActive = req.IsActive
	}

	_, err := c.deviceRegistryClient.UpdateDevice(c.outboundContext(ctx), protoReq)
	return err
}

// UpdateDeviceFields holds the mutable fields for UpdateDevice.
type UpdateDeviceFields struct {
	Label        string
	IPAddress    string
	PortNumber   int32
	Username     string
	Password     string
	Brand        string
	Timezone     string
	DeviceTypeID int64
	ShortLabel   string
	Description  string
	Notes        string
	IsActive     bool
}

// DeleteDevice removes a device from the payment registry (hard delete).
func (c *Client) DeleteDevice(ctx context.Context, id string) error {
	if err := c.ensureDialed(); err != nil {
		return err
	}

	_, err := c.deviceRegistryClient.DeleteDevice(
		c.outboundContext(ctx),
		&paymentv1.DeleteDeviceRequest{Id: id},
	)
	return err
}

// ---------------------------------------------------------------------------
// Proto converters
// ---------------------------------------------------------------------------

func protoDeviceToPublic(d *paymentv1.DeviceData) *DeviceData {
	if d == nil {
		return nil
	}
	return &DeviceData{
		ID:           d.Id,
		Label:        d.Label,
		CompanyID:    d.CompanyId,
		SerialNumber: d.SerialNumber,
		CloudID:      d.CloudId,
		IPAddress:    d.IpAddress,
		PortNumber:   d.PortNumber,
		Username:     d.Username,
		Password:     d.Password,
		Brand:        d.Brand,
		Model:        d.Model,
		DeviceTypeID: d.DeviceTypeId,
		Timezone:     d.Timezone,
		ShortLabel:   d.ShortLabel,
		Description:  d.Description,
		Notes:        d.Notes,
		IsActive:     d.IsActive,
	}
}

func publicDeviceToProto(d *DeviceData) *paymentv1.DeviceData {
	if d == nil {
		return nil
	}
	return &paymentv1.DeviceData{
		Id:           d.ID,
		Label:        d.Label,
		CompanyId:    d.CompanyID,
		SerialNumber: d.SerialNumber,
		CloudId:      d.CloudID,
		IpAddress:    d.IPAddress,
		PortNumber:   d.PortNumber,
		Username:     d.Username,
		Password:     d.Password,
		Brand:        d.Brand,
		Model:        d.Model,
		DeviceTypeId: d.DeviceTypeID,
		Timezone:     d.Timezone,
		ShortLabel:   d.ShortLabel,
		Description:  d.Description,
		Notes:        d.Notes,
		IsActive:     d.IsActive,
	}
}
