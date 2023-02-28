function GetListItem(input) {
	const sanitizedInput = sanitize(input);
	return (
		<li
			className={"foobar"}
			dangerouslySetInnerHTML={{ __html: sanitizedInput }}
		/>
	);
}
