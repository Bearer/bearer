const secure_url = "https://example.com/movies.json";

fetch(secure_url)
	.then((response) => response.json())
	.then((data) => console.log(data));
