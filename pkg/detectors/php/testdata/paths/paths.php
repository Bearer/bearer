<?php
public function getOnlyPath(Platform $platform): string
{

    // we want to ignore these
    require '/require/delivery-messages'.$test;
    require_once '/require/delivery-messages'.$test;

    $test = "/somepath";

    $response = $this->client->request(
        'GET',
        '/api/delivery-messages1'.$test // supported
    );

    if (Response::HTTP_OK === $response->getStatusCode()) {
        /** @var DeliveryMessage[] $deliveryMessages */
        $deliveryMessages = $this->serializer->deserialize($response->getContent(), DeliveryMessage::class.'[]', JsonEncoder::FORMAT);

        $deliveryMessages = array_filter(
            $deliveryMessages,
            static function (DeliveryMessage $deliveryMessage) use ($platform) {
                return $deliveryMessage->platform === (string) $platform;
            }
        );

        return array_pop($deliveryMessages)->message ?? '';
    }

    return '';
}

public function getWithURLConcatenated(Platform $platform, Customer $customer): string
{
    $response = $this->client->request(
        'GET',
        $apiURL . "/api/customers/" . $customer->Foo->id . "/transactions/" . $platform->Bar->id, // supported
    );

    if (Response::HTTP_OK === $response->getStatusCode()) {
        /** @var DeliveryMessage[] $deliveryMessages */
        $deliveryMessages = $this->serializer->deserialize($response->getContent(), DeliveryMessage::class.'[]', JsonEncoder::FORMAT);

        $deliveryMessages = array_filter(
            $deliveryMessages,
            static function (DeliveryMessage $deliveryMessage) use ($platform) {
                return $deliveryMessage->platform === (string) $platform;
            }
        );

        return array_pop($deliveryMessages)->message ?? '';
    }

    return '';
}

public function getWithURLPassedAsArg(PlatformID $platformID, CustomerID $customerID): string
{
    $response = $this->client->request(
        'GET',
        $apiURL,
        "/api/customers/" . $customerID . "/transactions/" . $platformID, // supported
    );

    if (Response::HTTP_OK === $response->getStatusCode()) {
        $deliveryMessages = $this->serializer->deserialize($response->getContent(), DeliveryMessage::class.'[]', JsonEncoder::FORMAT);

        $deliveryMessages = array_filter(
            $deliveryMessages,
            static function (DeliveryMessage $deliveryMessage) use ($platform) {
                return $deliveryMessage->platform === (string) $platform;
            }
        );

        return array_pop($deliveryMessages)->message ?? '';
    }

    return '';
}

public function getWithURLFullyInterpolated(PlatformID $platformID, CustomerID $customerID): string
{
    $response = $this->client->request(
        'GET',
        "{$_ENV["CUSTOMERS_HOST"]}:{$port}/api/delivery-messages?num_page={$page}&filters[]={$filters}", // supported
    );

    if (Response::HTTP_OK === $response->getStatusCode()) {
        $deliveryMessages = $this->serializer->deserialize($response->getContent(), DeliveryMessage::class.'[]', JsonEncoder::FORMAT);

        $deliveryMessages = array_filter(
            $deliveryMessages,
            static function (DeliveryMessage $deliveryMessage) use ($platform) {
                return $deliveryMessage->platform === (string) $platform;
            }
        );

        return array_pop($deliveryMessages)->message ?? '';
    }

    return '';
}

public function getTransactionsWithURLInterpolated(PlatformID $platformID, CustomerID $customerID): string
{
    $response = $this->client->request(
        'GET',
        "{$apiURL}api/customers/" . $customerID . "/transactions/" . $platformID, // supported
    );

    if (Response::HTTP_OK === $response->getStatusCode()) {
        /** @var DeliveryMessage[] $deliveryMessages */
        $deliveryMessages = $this->serializer->deserialize($response->getContent(), DeliveryMessage::class.'[]', JsonEncoder::FORMAT);

        $deliveryMessages = array_filter(
            $deliveryMessages,
            static function (DeliveryMessage $deliveryMessage) use ($platform) {
                return $deliveryMessage->platform === (string) $platform;
            }
        );

        return array_pop($deliveryMessages)->message ?? '';
    }

    return '';
}
?>