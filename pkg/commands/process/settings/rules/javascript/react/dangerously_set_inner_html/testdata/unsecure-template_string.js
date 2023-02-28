function GetListItem(input) {
	<li
		className={"foobar"}
		dangerouslySetInnerHTML={{ __html: `<a href=${input}>home page</a>` }}
	/>;
}
