package paymentclient

import (
	"context"

	paymentv1 "github.com/devbefingerspot/fio-paymentservice-go-client/internal/pb/payment/v1"
)

// DeviceType represents a physical device type's specifications and capacities.
type DeviceType struct {
	ID                 int64  `json:"id"`
	Name               string `json:"name"`
	Img                string `json:"img"`
	Brand              string `json:"brand"`
	AccurateCode       string `json:"accurate_code"`
	BasePrice          int64  `json:"base_price"`
	IsActive           bool   `json:"is_active"`
	FirmwareType       int32  `json:"firmware_type"`
	CapacityScan       int32  `json:"capacity_scan"`
	CapacityEmployee   int32  `json:"capacity_employee"`
	CapacityFinger     int32  `json:"capacity_finger"`
	CapacityPassword   int32  `json:"capacity_password"`
	CapacityRFID       int32  `json:"capacity_rfid"`
	CapacityFace       int32  `json:"capacity_face"`
	CapacityAdmin      int32  `json:"capacity_admin"`
	CapacityVein       int32  `json:"capacity_vein"`
	OnlineVerification bool   `json:"online_verification"`
	ReadyFIO           bool   `json:"ready_fio"`
	ReadyIOT           bool   `json:"ready_iot"`
	CanAutoRestart     bool   `json:"can_auto_restart"`
	DefaultAutoRestart bool   `json:"default_auto_restart"`
	IsIOT              bool   `json:"is_iot"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

// ListDeviceTypes retrieves all master device types from the payment service.
func (c *Client) ListDeviceTypes(ctx context.Context) ([]DeviceType, error) {
	if err := c.ensureDialed(); err != nil {
		return nil, err
	}

	resp, err := c.masterDeviceClient.ListDeviceTypes(
		c.outboundContext(ctx),
		&paymentv1.ListDeviceTypesRequest{},
	)
	if err != nil {
		return nil, err
	}

	result := make([]DeviceType, len(resp.GetDeviceTypes()))
	for i, item := range resp.GetDeviceTypes() {
		result[i] = DeviceType{
			ID:                 item.GetId(),
			Name:               item.GetName(),
			Img:                item.GetImg(),
			Brand:              item.GetBrand(),
			AccurateCode:       item.GetAccurateCode(),
			BasePrice:          item.GetBasePrice(),
			IsActive:           item.GetIsActive(),
			FirmwareType:       item.GetFirmwareType(),
			CapacityScan:       item.GetCapacityScan(),
			CapacityEmployee:   item.GetCapacityEmployee(),
			CapacityFinger:     item.GetCapacityFinger(),
			CapacityPassword:   item.GetCapacityPassword(),
			CapacityRFID:       item.GetCapacityRfid(),
			CapacityFace:       item.GetCapacityFace(),
			CapacityAdmin:      item.GetCapacityAdmin(),
			CapacityVein:       item.GetCapacityVein(),
			OnlineVerification: item.GetOnlineVerification(),
			ReadyFIO:           item.GetReadyFio(),
			ReadyIOT:           item.GetReadyIot(),
			CanAutoRestart:     item.GetCanAutoRestart(),
			DefaultAutoRestart: item.GetDefaultAutoRestart(),
			IsIOT:              item.GetIsIot(),
			CreatedAt:          item.GetCreatedAt(),
			UpdatedAt:          item.GetUpdatedAt(),
		}
	}

	return result, nil
}
