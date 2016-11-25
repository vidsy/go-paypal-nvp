package payload

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

const (
	// ReceiverTypeEmail sets receiver type to email address.
	ReceiverTypeEmail = "EmailAddress"

	// ReceiverTypePhone sets receiver type to phone number.
	ReceiverTypePhone = "PhoneNumber"

	// ReceiverTypeUserID sets receiver type to user id.
	ReceiverTypeUserID = "UserID"
)

type (
	// MassPayment payload for mass payment request.
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
func NewMassPayment(currency string, receiverType string) *MassPayment {
	return &MassPayment{
		Method:       "MassPay",
		CurrencyCode: currency,
		ReceiverType: receiverType,
	}
}

// Credentials sets credentials and API version.
func (mp *MassPayment) SetCredentials(user string, password string, signature string, apiVersion string) {
	mp.User = user
	mp.Password = password
	mp.Signature = signature
	mp.Version = apiVersion
}

// AddItem adds an item to the mass payment items array.
func (mp *MassPayment) AddItem(item MassPaymentItem) {
	mp.Items = append(mp.Items, item)
}

// Total gives the total of all payments in the mass payment.
func (mp MassPayment) Total() float64 {
	total := 0.0
	for _, item := range mp.Items {
		total += item.Amount
	}

	return total
}

// Serialize convert struct into NVP key=value format for the masspayment.
func (mp MassPayment) Serialize() (string, error) {
	data := url.Values{}
	massPaymentType := reflect.TypeOf(mp)
	massPaymentValue := reflect.ValueOf(mp)

	for i := 0; i < massPaymentType.NumField(); i++ {
		field := massPaymentType.Field(i)
		if fieldTag, ok := field.Tag.Lookup("nvp_field"); ok {
			valueContents := massPaymentValue.Field(i).String()
			if valueContents != "" {
				data.Set(fieldTag, massPaymentValue.Field(i).String())
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

// Serialize convert mass payment item into key=value pair and add to existing
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
