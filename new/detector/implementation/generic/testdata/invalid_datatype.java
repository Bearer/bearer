public class HashingAssignment extends AssignmentEndpoint {
  public String getMd5(HttpServletRequest request) {
    String secret = SECRETS[new Random().nextInt(SECRETS.length)];
  }
}