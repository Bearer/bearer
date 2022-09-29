// classes are supported
public partial class Customer
{
    // properties are supported
    private int id;
    private User name;
    private string address;
    public List<Order> orders;
    // intialized properties are supported
    private const string ConsumerKeyKey = "oauth_consumer_key";
    // shorthands for getters and setters are supported as properties
    public int Age { get; set; }
    
    // constructors are supported
    public Customer(int constructor1)
    {
    // not suported
		var v = new { Amount = 108, Message = "Hello" };
    }
    
    // type of a is undefined
    public void test( a, b int) {
    }

    public static string dosomething(int a, string b){
        // delegate parameters are supported
        string resource = request.Parameters.Where(delegate(Parameter p)
        {
            return (p.Type == ParameterType.UrlSegment);

        })
    }
}