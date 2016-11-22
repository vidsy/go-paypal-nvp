package nvppayload_test

import (
	"testing"

	"github.com/vidsy/go-paypalnvp/nvppayload"
)

func TestMassPayment(t *testing.T) {
	t.Run(".AddItem", func(t *testing.T) {
		t.Run("AddsToItemArray", func(t *testing.T) {
			item := nvppayload.MassPaymentItem{}
			massPayment := nvppayload.NewMassPayment("user", "password", "signature", nvppayload.ReceiverTypeEmail)
			massPayment.AddItem(item)

			if len(massPayment.Items) != 1 {
				t.Fatalf("Expected 1 item, got: %d", len(massPayment.Items))
			}
		})
	})

	t.Run(".Serialize", func(t *testing.T) {
		t.Run("ReturnsErrorWhenNoDataSet", func(t *testing.T) {
			massPayment := nvppayload.NewMassPayment("user", "password", "signature", nvppayload.ReceiverTypeEmail)
			_, err := massPayment.Serialize()

			if err == nil {
				t.Fatalf("Expected error, got: %v", err)
			}
		})

		t.Run("ReturnsCorrectlySerializedPayload", func(t *testing.T) {
			massPayment := nvppayload.NewMassPayment("user", "password", "signature", nvppayload.ReceiverTypeEmail)
			massPayment.Version = "1.0"
			massPayment.EmailSubject = "Test email"
			itemOne := nvppayload.MassPaymentItem{
				Email:  "test@test.com",
				Amount: 1.50,
				ID:     "123456789",
				Note:   "A test transaction",
			}
			itemTwo := nvppayload.MassPaymentItem{
				Email:  "test@testtwo.com",
				Amount: 1.60,
				ID:     "1234567810",
				Note:   "Another test transaction",
			}

			massPayment.AddItem(itemOne)
			massPayment.AddItem(itemTwo)

			expectedPayload := `CURRENCYCODE=GBP&EMAILSUBJECT=Test+email&L_AMT0=1.50&L_AMT1=1.60&L_EMAIL0=test%40test.com&L_EMAIL1=test%40testtwo.com&L_NOTE0=A+test+transaction&L_NOTE1=Another+test+transaction&L_UNIQUEID0=123456789&L_UNIQUEID1=1234567810&METHOD=MassPay&PWD=password&RECEIVERTYPE=EmailAddress&SIGNATURE=signature&USER=user&VERSION=1.0`
			payload, _ := massPayment.Serialize()

			if expectedPayload != payload {
				t.Fatalf("Expected payload to be: '%s', got '%s'", expectedPayload, payload)
			}
		})
	})
}
