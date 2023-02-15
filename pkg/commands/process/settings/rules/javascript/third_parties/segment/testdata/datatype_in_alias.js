import { Analytics } from "@segment/analytics-node"

const prodAnalytics = new Analytics({ writeKey: "product-write-key" });
const appAnalytics = new Analytics({ writeKey: "application-write-key" });

var user = getCurrentUser();
appAnalytics.alias({
  previousId: user.email,
  userId: user.id,
});
