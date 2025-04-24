package inridereports

import (
    "testing"

    "bytes"
    "crowdsourcing/tools/osrmtool"
    "encoding/json"
    "net/http"
    "net/http/httptest"

    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
)

func TestRemoveExcessEvents(t *testing.T) {
	locations := []osrmtool.Location{
		{1.0, 1.0},
		{2.0, 2.0},
		{3.0, 3.0},
		{4.0, 4.0},
	}
	events := []Event{
		{Uid: "1"},
		{Uid: "2"},
		{Uid: "3"},
		{Uid: "4"},
	}

	maxAllowed := 2
	removeExcessEvents(&locations, &events, maxAllowed)

	if len(locations) != maxAllowed || len(events) != maxAllowed {
		t.Errorf("expected %d events and locations after trimming, got %d and %d", maxAllowed, len(events), len(locations))
	}
}

func TestIsReportExpired(t *testing.T) {
	report := Event{
		UpdateTimestamp: 1000,
		TtlMinutes:      1,
	}

	atTimestamp := int64(160000)
	if !isReportExpired(report, atTimestamp) {
		t.Error("expected report to be expired")
	}

	atTimestamp = int64(1000 + 30*1000)
	if isReportExpired(report, atTimestamp) {
		t.Error("expected report to NOT be expired")
	}
}

func TestExtractEventsFromClusters(t *testing.T) {
	events := []Event{
		{Uid: "a", PositiveConfirmation: 1, NegativeConfirmation: 0, IsAggregated: false},
		{Uid: "b", PositiveConfirmation: 2, NegativeConfirmation: 1, IsAggregated: true},
		{Uid: "c", PositiveConfirmation: 3, NegativeConfirmation: 2, IsAggregated: false},
	}
	clusters := [][]int64{{0, 1}, {2}}
	centers := []int64{1, 2}

	clustered := map[string]Event{}
	startTime := int64(9999)

	extractEventsFromClusters(clustered, events, clusters, centers, startTime)

	if len(clustered) != 2 {
		t.Fatalf("expected 2 clustered events, got %d", len(clustered))
	}

	for _, evt := range clustered {
		if !evt.IsAggregated {
			t.Error("expected event to be marked as aggregated")
		}
		if evt.UpdateTimestamp != startTime {
			t.Errorf("expected timestamp to be %d, got %d", startTime, evt.UpdateTimestamp)
		}
	}
}



func TestDeserializeRequest(t *testing.T) {
	e := echo.New()

	body := map[string]any{
		"type": POLICE_ONLINE,
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	report, err := DeserializeRequest(ctx)
	assert.NoError(t, err)
	assert.Equal(t, InRideReportType(POLICE_ONLINE), report.Type)
	assert.True(t, report.Online)
}

func TestDeserializeRequest_Invalid(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{invalid json}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	_, err := DeserializeRequest(ctx)
	assert.Error(t, err)
}

func TestSerializeResponse(t *testing.T) {
	reportId := "abc123"
	resp := SerializeResponse(reportId)

	assert.Equal(t, reportId, resp.ReportId)
}

func TestDeserializeConfirmRequest(t *testing.T) {
	e := echo.New()

	body := map[string]any{}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/confirm/123", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/confirm/:report_id")
	ctx.SetParamNames("report_id")
	ctx.SetParamValues("123")

	confirm, err := DeserializeConfirmRequest(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "123", confirm.ReportId)
}

func TestDeserializeConfirmRequest_Invalid(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/confirm/123", bytes.NewBufferString("{bad json}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetPath("/confirm/:report_id")
	ctx.SetParamNames("report_id")
	ctx.SetParamValues("123")

	_, err := DeserializeConfirmRequest(ctx)
	assert.Error(t, err)
}

func TestSerializeConfirmResponse(t *testing.T) {
	resp := SerializeConfirmResponse("dummy")
	assert.IsType(t, ReportConfirmationResponse{}, resp)
}
