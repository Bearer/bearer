package test.example;

public class Env {
    public static void main (String[] args) {
      String accountId = System.getenv("ACCOUNT_ID");
      String orderServiceUrl = System.getenv("ORDER_SERVICE_URL") + "/path?x=" + accountId;

      // TEST: ignores other methods on System
      String ignore = System.other("IGNORE_ME_HOST");
      // TEST: ignores other packages
      String ignore2 = Other.getenv("IGNORE_ME_URL");

      System.out.println(someVar["ignored.domain.com"]);
    }
}
