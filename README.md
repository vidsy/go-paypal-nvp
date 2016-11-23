<h1 align="center">go-paypal-nvp</h1>

<p align="center">
  Library for interacting with the <a href="https://developer.paypal.com/docs/classic/api/NVPAPIOverview/">PayPal NVP endpoint</a>.
</p>

## Usage

Currently the library supports the following methods:

* MassPayment

With others coming soon.

### Authentication

Only API credentials are supported at present, however the client takes a client that implements the `TransportClient` interface so 
a http client can be created with client cert authentication setup and passed in.

### Example

Below is an example setting up and preforming a `Mass Payment`:

```go
package main

import (
	"fmt"

	"github.com/vidsy/go-paypalnvp"
	"github.com/vidsy/go-paypalnvp/payload"
)

func main() {

	client := paypalnvp.NewClient(
		nil,
		paypalnvp.Sandbox,
		"user",
		"password",
		"signature",
	)

	massPayment := payload.NewMassPayment("GBP", payload.ReceiverTypeEmail)
	massPaymentItem := payload.MassPaymentItem{
		Email:  "tech+paypal-buyer@vidsy.co",
		Amount: 100.50,
		ID:     "123456789",
		Note:   "Vidsy payment going out",
	}
	massPayment.AddItem(massPaymentItem)

	response, err := client.Execute(massPayment)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Successful: %t\n", response.Successful())
	fmt.Printf("Errors count: %d\n", response.ErrorCount())

	if response.ErrorCount() > 0 {
		fmt.Printf("Error code: %s\n", response.Errors[0].Code)
		fmt.Printf("Long message: %s\n", response.Errors[0].LongMessage)
	}
}
```
