package payload

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

const (
	// receiverTypeEmail sets receiver type to email address.
	ReceiverTypeEmail = "EmailAddress"

	// receiverTypePhone sets receiver type to phone number.
	ReceiverTypePhone = "PhoneNumber"

	// receiverTypeUserID sets receiver type to user id.
	ReceiverTypeUserID = "UserID"
)

type (
	//MassPayment payload for mass payment request.
	MassPayment struct {
		User         string `nvp_field:"USER"`
		Password     string `nvp_field:"PWD"`
		Signature    string `nvp_field:"SIGNATURE"`
		Version      string `nvp_field:"VERSION"`
		Method       string `nvp_field:"METHOD"`
		EmailSubject string `nvp_field:"EMAILSUBJECT"`
		CurrencyCode string `nvp_field:"CURRENCYCODE"`
		ReceiverType string `nvp_field:"RECEIVERTYPE"`
		Items        []MassPaymentItem
	}

	// MassPaymentItem contains data about an individual mass payment
	// item.
	MassPaymentItem struct {
		Email  string  `nvp_field:"L_EMAIL"`
		Phone  string  `nvp_field:"L_RECEIVERPHONE"`
		UserID string  `nvp_field:"L_RECEIVERID"`
		Amount float64 `nvp_field:"L_AMT"`
		ID     string  `nvp_field:"L_UNIQUEID"`
		Note   string  `nvp_field:"L_NOTE"`
	}
)

// NewMassPayment creates a new MassPayment struct with defaults
func NewMassPayment(user string, password string, signature string, receiverType string) *MassPayment {
	return &MassPayment{
		User:         user,
		Password:     password,
		Signature:    signature,
		Method:       "MassPay",
		CurrencyCode: "GBP",
		ReceiverType: receiverType,
	}
}

// AddItem Adds an item to the mass payment items array.
func (mp *MassPayment) AddItem(item MassPaymentItem) {
	mp.Items = append(mp.Items, item)
}

// Serialize Convert struct into NVP key=value format for the masspayment.
func (mp MassPayment) Serialize() (string, error) {
	data := url.Values{}
	mpType := reflect.TypeOf(mp)
	mpValue := reflect.ValueOf(mp)

	for i := 0; i < mpType.NumField(); i++ {
		field := mpType.Field(i)
		if fieldTag, ok := field.Tag.Lookup("nvp_field"); ok {
			valueContents := mpValue.Field(i).String()
			if valueContents != "" {
				data.Set(fieldTag, mpValue.Field(i).String())
			}
		}
	}

	for i, item := range mp.Items {
		item.Serialize(&data, i)
	}

	if len(mp.Items) == 0 {
		return "", errors.New("Expected at least one mass payment item")
	}

	return data.Encode(), nil
}

// Serialize Convert mass payment item into key=value pair and add to existing
// Values struct.
func (mpi MassPaymentItem) Serialize(data *url.Values, index int) {
	itemType := reflect.TypeOf(mpi)
	itemValue := reflect.ValueOf(mpi)

	for i := 0; i < itemType.NumField(); i++ {
		field := itemType.Field(i)
		if fieldTag, ok := field.Tag.Lookup("nvp_field"); ok {
			value := itemValue.Field(i)

			switch field.Type.Kind() {
			case reflect.String:
				valueContents := value.String()
				if valueContents != "" {
					data.Set(fmt.Sprintf("%s%d", fieldTag, index), value.String())
				}
			case reflect.Float64:
				data.Set(fmt.Sprintf("%s%d", fieldTag, index), strconv.FormatFloat(value.Float(), 'f', 2, 64))
			}
		}
	}
}
