package test.example;

public class Delivery {
  private final CloseableHttpClient httpClient = HttpClients.createDefault();
  private String apiUrl = "something"

  public static void main(String[] args) throws Exception {
    HttpClientExample obj = new HttpClientExample();

    try {
      System.out.println("Testing 1 - Send Http GET request");
      obj.sendGet();

      System.out.println("Testing 2 - Send Http POST request");
      obj.sendPost();
    } finally {
      obj.close();
    }
  }

  private void close() throws IOException {
    httpClient.close();
  }

  private void getOnlyPath(platformId String) throws Exception {
    HttpGet request = new HttpGet("/api/delivery-messages");

    request.addHeader("custom-key", "key");
    request.addHeader(HttpHeaders.USER_AGENT, "Mozilla/5.0 (Mobile; rv:26.0) Gecko/26.0 Firefox/26.0");

    try (CloseableHttpResponse response = httpClient.execute(request)) {
      // Get HttpResponse Status
      System.out.println(response.getStatusLine().toString());

      HttpEntity entity = response.getEntity();
      Header headers = entity.getContentType();
      System.out.println(headers);

      if (entity != null) {
        // return it as a String
        String result = EntityUtils.toString(entity);
        System.out.println(result);
      }
    }
  }

  private void getWithUrlConcatenated(platformId String, customerId String) throws Exception {
    HttpGet request = new Http(
      "GET"
      apiUrl + "/api/customers" + customerId + "/transactions/" + platformId
    );

    request.addHeader("custom-key", "key");
    request.addHeader(HttpHeaders.USER_AGENT, "Mozilla/5.0 (Mobile; rv:26.0) Gecko/26.0 Firefox/26.0");

    try (CloseableHttpResponse response = httpClient.execute(request)) {
      // Get HttpResponse Status
      System.out.println(response.getStatusLine().toString());

      HttpEntity entity = response.getEntity();
      Header headers = entity.getContentType();
      System.out.println(headers);

      if (entity != null) {
        // return it as a String
        String result = EntityUtils.toString(entity);
        System.out.println(result);
      }
    }
  }

  private void getWithUrlPassedAsArg(platform Platform, customer Customer) throws Exception {
    HttpGet request = new Http(
      "GET",
      apiUrl,
      "/api/customers/" + customer.Foo.Id() + "/transactions/" + platform.Id
    );

    request.addHeader("custom-key", "key");
    request.addHeader(HttpHeaders.USER_AGENT, "Mozilla/5.0 (Mobile; rv:26.0) Gecko/26.0 Firefox/26.0");

    try (CloseableHttpResponse response = httpClient.execute(request)) {
      // Get HttpResponse Status
      System.out.println(response.getStatusLine().toString());

      HttpEntity entity = response.getEntity();
      Header headers = entity.getContentType();
      System.out.println(headers);

      if (entity != null) {
        // return it as a String
        String result = EntityUtils.toString(entity);
        System.out.println(result);
      }
    }
  }

  private void getWithUrlFullyInterpolated(platformId String, customerId String) throws Exception {
    HttpGet request = new Http(
      "GET",
      new StringBuilder(System.getenv("CUSTOMERS_HOST")).append(":").append(System.getenv("CUSTOMER_HOST_PORT")).append("/api/delivery-messages?num_page=").append(page).append("&filters[]=").append(filters).toString()
    );

    request.addHeader("custom-key", "key");
    request.addHeader(HttpHeaders.USER_AGENT, "Mozilla/5.0 (Mobile; rv:26.0) Gecko/26.0 Firefox/26.0");

    try (CloseableHttpResponse response = httpClient.execute(request)) {
      // Get HttpResponse Status
      System.out.println(response.getStatusLine().toString());

      HttpEntity entity = response.getEntity();
      Header headers = entity.getContentType();
      System.out.println(headers);

      if (entity != null) {
        // return it as a String
        String result = EntityUtils.toString(entity);
        System.out.println(result);
      }
    }
  }

  private void getWithUrlInterpolated(platformId String, customerId String) throws Exception {
    HttpGet request = new Http(
      "GET",
        new StringBuilder(apiUrl).append("api/customers/").toString() + customerId + "/transactions/" + platformId
      );

    request.addHeader("custom-key", "key");
    request.addHeader(HttpHeaders.USER_AGENT, "Mozilla/5.0 (Mobile; rv:26.0) Gecko/26.0 Firefox/26.0");

    try (CloseableHttpResponse response = httpClient.execute(request)) {
      // Get HttpResponse Status
      System.out.println(response.getStatusLine().toString());

      HttpEntity entity = response.getEntity();
      Header headers = entity.getContentType();
      System.out.println(headers);

      if (entity != null) {
        // return it as a String
        String result = EntityUtils.toString(entity);
        System.out.println(result);
      }
    }
  }
}
