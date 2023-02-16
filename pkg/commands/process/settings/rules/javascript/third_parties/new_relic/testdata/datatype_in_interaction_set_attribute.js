const newrelic = require("newrelic")

router.addRoute("/order", () => {
  const user = getCurrentUser();
  newrelic.interaction()
    .setAttribute("username", user.first_name)
    .setAttribute("postal-code", user.post_code);
  renderCart(user);
});

// alternative syntax
var interaction = newrelic.interaction()
interaction.setAttribute("email", user.email_address)