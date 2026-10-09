package paymentclient_test

import (
	"context"
	"net"
	"testing"

	paymentclient "github.com/devbefingerspot/fio-paymentservice-go-client"
	paymentv1 "github.com/devbefingerspot/fio-paymentservice-go-client/internal/pb/payment/v1"
	"google.golang.org/grpc"
)

type mockMasterDeviceServer struct {
	paymentv1.UnimplementedMasterDeviceServiceServer
}

func (s *mockMasterDeviceServer) ListDeviceTypes(
	_ context.Context,
	_ *paymentv1.ListDeviceTypesRequest,
) (*paymentv1.ListDeviceTypesResponse, error) {
	return &paymentv1.ListDeviceTypesResponse{
		DeviceTypes: []*paymentv1.DeviceTypeItem{
			{
				Id:           10,
				Name:         "Revo WFV",
				Brand:        "revo",
				AccurateCode: "ACC-01",
				BasePrice:    2000000,
				IsActive:     true,
				CapacityScan: 50000,
			},
		},
	}, nil
}

func TestClient_ListDeviceTypes(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	paymentv1.RegisterMasterDeviceServiceServer(srv, &mockMasterDeviceServer{})

	go func() {
		_ = srv.Serve(lis)
	}()
	defer srv.Stop()

	client := paymentclient.New(lis.Addr().String(), "test-key").WithInsecure()
	defer client.Close()

	items, err := client.ListDeviceTypes(context.Background())
	if err != nil {
		t.Fatalf("ListDeviceTypes failed: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].ID != 10 || items[0].Name != "Revo WFV" || items[0].CapacityScan != 50000 {
		t.Errorf("unexpected item: %+v", items[0])
	}
}
