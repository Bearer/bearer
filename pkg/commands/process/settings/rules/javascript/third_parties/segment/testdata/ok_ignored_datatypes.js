import { Analytics } from '@segment/analytics-node'

const client = new Analytics({ write_key: 'some-write-key' });

client.track({
  event: "some event name",
  userId: user.id
});