package payment

import (
	"app/internal/platform/testutil"
	"bytes"
	"encoding/json"
	"log"
	"os"
	"testing"
)

func TestGetPayments(t *testing.T) {

	db := testutil.SetupTestDB(t)

	logger := log.New(os.Stdout, "app-test", log.LstdFlags)
	repo := NewPaymentRepository(db)
	service := NewPaymentService(repo)
	apiHandler := NewApiPaymentHandler(service, logger)

	httpResponseWriter := testutil.NewMockHTTPResponseWriter()
	httpRequest := testutil.NewMockHTTPRequest("GET", "/api/v1/payments", nil)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("handler panicked: %v", r)
		}
	}()

	apiHandler.getPayments(httpResponseWriter, httpRequest)

	body := httpResponseWriter.Body()
	jsonstring, _ := prettyPrintJSON(string(body))

	t.Logf("Response body: %s", jsonstring)

	if httpResponseWriter.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type to be application/json, got %s", httpResponseWriter.Header().Get("Content-Type"))
	}

}

func prettyPrintJSON(input string) (string, error) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(input), "", "  ")
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
