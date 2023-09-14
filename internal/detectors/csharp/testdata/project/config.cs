using System;

class Env {
    public static void main() {
      var orderServiceUrl = Environment.GetEnvironmentVariable("ORDER_SERVICE_URL");
      var customerServiceHost = System.Environment.GetEnvironmentVariable(@"CUSTOMER_SERVICE_HOST", EnvironmentVariableTarget.Process);
      var accountId = Environment.GetEnvironmentVariable("ACCOUNT_ID");

      // TEST: ignores other methods on Environment
      var ignore = Environment.Other("IGNORE_ME_HOST");
      // TEST: ignores other classes
      var ignore2 = Other.GetEnvironmentVariable("IGNORE_ME_URL");
    }
}
