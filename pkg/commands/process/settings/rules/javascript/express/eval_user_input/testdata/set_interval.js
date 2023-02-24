app.post("/:id", (req, res) => {
	userInput = req.params.id;
	var command = "new Function('" + userInput + "')";
	setInterval(command);
});
