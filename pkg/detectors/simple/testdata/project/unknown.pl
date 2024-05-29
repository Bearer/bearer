# TEST: ignores lines starting with an asterisk
/*
 * var url = 'https://ignore.test-domain.com'
 */

# TEST: ignores anything after a hash
doSomething(); # https://ignore.test-domain.com

# TEST: ignores anything after a double slash
doSomething(); // https://ignore.test-domain.com

# TEST: ignores hostnames (must be URL with scheme)
makeARequest("ignore.me")

# TEST: extracts URLs
makeARequest("https://url.example.com");

# TEST: extracts multiple URLs in a line, ignoring comment
makeARequest("https://multi-a.example.com/foo?x=1", "https://multi-b.example.com/bar"); // https://ignore.example.com

;; makeARequest("https://multi-a.example.com/foo?x=1", "https://multi-b.example.com/bar"); // https://ignore.example.com

# TEST: works with ports
https://port1.example.com:3000/foo?x=1 https://port2.example.com:3000

# TEST: extracts from square brackets
[http://link.example.com](hey)
