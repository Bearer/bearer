public class Greet {
  const Greeting = "Hello World";

  public static void main(String[] args)
  {
    var s = Greeting + "!";
    s += "!!";

    String s2 = "hey ";
    s2 += args[0];
    s2 += " there";
  }
}
