package paths

import (
	"fmt"
	"log"
	"os"
)

var (
	apiUrl = "something.com"
)

type Customer struct {
	Client Client
}

type Client struct {
	Id string
}

func getOnlyPath(platform string) (error, []string) {
	messages, err := delivery_client.Call(
		"GET",
		"/api/delivery-messages",
	)

	if err != nil {
		log.Fatalln(err)
	}

	return nil, messages
}

func getWithUrlConcatenated(platformId string, customer Customer) (error, []string) {
	messages, err := delivery_client.Call(
		"GET",
		apiUrl+"/api/shop/customers/"+customer.Client.Id+"/transactions/"+platformId,
	)

	if err != nil {
		log.Fatalln(err)
	}

	return nil, messages
}

func getWithUrlPassedAsArg(platformId string, customerId string) (error, []string) {
	messages, err := delivery_client.Call(
		"GET",
		apiUrl,
		"/api/shop/customers/"+customerId+"/transactions/"+platformId,
	)

	if err != nil {
		log.Fatalln(err)
	}

	return nil, messages
}

func getWithUrlFullyInterpolated(platformId string, customerId string, page string, filters string) (error, []string) {
	url := fmt.Sprintf(
		"%s:%s/api/shop/delivery-messages"+"?num_page=%s&filters[]=%s",
		os.Getenv("CUSTOMERS_FOO"),
		os.Getenv("CUSTOMER_FOO_PORT"),
		page,
		filters,
	)

	messages, err := delivery_client.Call(
		"GET",
		url,
	)

	if err != nil {
		log.Fatalln(err)
	}

	return nil, messages
}

func getTransactionsWithUrlInterpolated(platformId string, customerId string) (error, []string) {
	url := fmt.Sprintf(
		"%s/api/shop/customers/",
		apiUrl,
	)

	messages, err := delivery_client.Call(
		"GET",
		url+customerId+"/transactions/"+platformId,
	)

	if err != nil {
		log.Fatalln(err)
	}

	return nil, messages
}
