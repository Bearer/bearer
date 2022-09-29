using System.Collections.Generic;

namespace CSharpTest
{
	class CSharpTest
	{
		public string Name { get; set; }

		public Dictionary<string, string> A(string s)
		{
			return new Dictionary<string, string>
			{
				// TEST: ignores strings in initializer indices
				// TEST: finds strings
				{ "ignore.domain.com", "string.example.com" },
				// TEST: finds raw strings
				{ "ignore2.domain.com", @"raw-string.example.com" },
				// TEST: finds interpolated strings
				{ "ignore3.domain.com", $"a.{s}.example.com" },
				// TEST: finds interpolated raw strings
        { "ignore4.domain.com", $@"{s}.interpolated-raw.example.com"}
			};
		}

		public Dictionary<string, string> B()
		{
			return new Dictionary<string, string>
			{
				// TEST: finds strings in initializer values
				// TEST: ignores strings in initializer indices
				["ignore.domain.com"] = "other-dict-literal.example.com",
			};
		}

		public static CSharpTest C() {
			// TEST: finds strings in class property initializers
			return new CSharpTest { Name = "prop.example.com" };
		}

		private void Get(string name) {
		}

		private void Something(string excluded) {
		}

		public List<string> D() {
			// TEST: finds strings assigned to variables
			var s = ".concat.example.com";

      var z = "simple" + ".concat.example.com";

			// TEST: ignores string comparisons
			if ("ignore.domain.com" == "ignore" || "ignore.domain.com" != "ignore") {
				return null;
			}

			return new List<string>
			{
				// TEST: finds strings in list initializers
				"list-string.example.com",
				"subdomain" + s,
				// TEST: ignores indices
        B()["ignore.domain.com"],
			};
		}
	}
}
