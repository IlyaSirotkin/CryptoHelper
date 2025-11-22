package exchange_datasource

import (
	setup "cryptoHelper/setup"
	"testing"
)

func TestExtractData(t *testing.T) {
	err := setup.SetENVreading("../../../config/env_file.env")
	if err != nil {
		t.Error(err)
	}
	currenciesToTest := []string{"BTC", "ETH", "ONDO", "XRP", "ADA", "SOL"}
	exc := NewExchange()
	for _, el := range currenciesToTest {
		price, err := exc.ExtractData(el)
		if err != nil {
			t.Errorf("ExtractData with %s return error %e", el, err)
		} else if price == 0.0 {
			t.Errorf("ExtractData return 0.0 with %s", el)
		}
	}
}
