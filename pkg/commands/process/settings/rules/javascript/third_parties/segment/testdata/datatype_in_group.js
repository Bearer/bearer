import { Analytics } from '@segment/analytics-node'

const analytics = new Analytics({ writeKey: 'my-write-key' });

var user = getCurrentUser();
analytics.group({
  userId: user.id,
  groupId: user.job_title,
  traits: {},
});