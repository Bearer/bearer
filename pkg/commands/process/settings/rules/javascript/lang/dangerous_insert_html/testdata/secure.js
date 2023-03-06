function renderListItem(input) {
	this.ref.insertAdjacentHTML("beforebegin", `<li>fixed list item</li>`);
	this.ref.insertAdjacentHTML("beforebegin", `<li>fixed list item</li>`);
	this.ref.innerHTML = "";
	React.createElement(CssPropTringle, {
		s: 100,
		x: 0,
		y: 0,
	})
}
