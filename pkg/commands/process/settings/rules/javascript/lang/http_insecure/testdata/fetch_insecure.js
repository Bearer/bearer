const insecure_url = "http://example.com/movies.json";

fetch(insecure_url)
	.then((response) => response.json())
	.then((data) => console.log(data));
