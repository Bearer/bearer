require(path);

app.get("/bad", (req, _res) => {
	try {
		const { path } = req.query;

		require(path);
	} catch (err) {
		// handle error
	}
});
