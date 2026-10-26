package external

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientDecodesMockResponses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/fns/v1/companies/7700000000":
			_, _ = w.Write([]byte(`{"inn":"7700000000","status":"ACTIVE","legalForm":"OOO","okvedSection":"G","registeredAt":"2023-09-01","taxDebtMinor":500}`))
		case "/bki/v1/reports/7700000000":
			_, _ = w.Write([]byte(`{"inn":"7700000000","company":null,"owner":{"found":true,"overdueCount12m":2,"maxOverdueDays24m":45,"activeLoans":1}}`))
		default:
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}))
	defer server.Close()
	c := &Client{BaseURL: server.URL, HTTP: server.Client()}

	registry, e := c.Registry(context.Background(), "7700000000")
	if e != nil || registry.Status != "ACTIVE" || registry.OkvedSection != "G" || registry.TaxDebtMinor != 500 {
		t.Fatalf("registry decoded incorrectly: %+v %v", registry, e)
	}

	report, e := c.CreditReport(context.Background(), "7700000000")
	if e != nil || report.Company != nil || report.Owner == nil || report.Owner.OverdueCount12m != 2 {
		t.Fatalf("credit report decoded incorrectly: %+v %v", report, e)
	}

	if _, e = c.Account(context.Background(), "7700000000"); e == nil {
		t.Fatal("unavailable source returned no error")
	}
}

func TestClientErrorsDoNotContainINN(t *testing.T) {
	c := &Client{BaseURL: "http://127.0.0.1:1", HTTP: http.DefaultClient}
	_, e := c.Registry(context.Background(), "7712345678")
	if e == nil || strings.Contains(e.Error(), "7712345678") {
		t.Fatalf("error leaks INN or is missing: %v", e)
	}
}
