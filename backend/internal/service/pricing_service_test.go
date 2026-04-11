package service

import (
	"encoding/json"
	"testing"
)

type MockCubepathClient struct {
	Resp map[string]json.RawMessage
	Err  error
}

func (m *MockCubepathClient) Get(path string) (json.RawMessage, error) {
	return m.Resp[path], m.Err
}
func TestGetPlans(t *testing.T) {
	mock := &MockCubepathClient{
		Resp: map[string]json.RawMessage{
			"/vps/plans": json.RawMessage(`{"locations": [{"location_name": "test"}]}`),
		},
		Err: nil,
	}

	service := NewPricingService()
	plansResp, err := service.GetPlans(mock)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if plansResp.Locations[0].LocationName != "test" {
		t.Fatalf("expected location name test, got: %v", plansResp.Locations[0].LocationName)
	}
}
