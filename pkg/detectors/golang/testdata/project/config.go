package config

import (
	"os"

	sneeky "os"
)

var url = "http://" + os.Getenv("ORDER_SERVICE_DOMAIN") + "/whatever"
var orderServiceUrl = os.Getenv("ORDER_SERVICE_URL")

var orderServiceUrl = os.Getenv("ORDER_SERVICE_ADDR")
var orderServiceUrl = os.Getenv("ORDER_SERVICE_SVC")

MustMapEnv(&svc.productCatalogSvcAddr, "PRODUCT_CATALOG_SERVICE_ADDR")
MustGetEnv("CURRENCY_SERVICE_ADDR")


var userServiceHost, _ = sneeky.LookupEnv(`USER_SERVICE_HOST`)
var accountId = os.Getenv("ACCOUNT_ID")
var other = os.Other("IGNORE_ME_HOST")
