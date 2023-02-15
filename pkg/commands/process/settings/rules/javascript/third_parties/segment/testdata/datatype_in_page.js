import { Analytics } from '@segment/analytics-node'

const analytics = new Analytics({ writeKey: 'my-write-key' });

var customer = getCurrentUser();
analytics.page({
  userId: customer.id,
  category: "Shopping Cart",
  properties: {
    path: "/cart/"+customer.bank_account_number
  },
});