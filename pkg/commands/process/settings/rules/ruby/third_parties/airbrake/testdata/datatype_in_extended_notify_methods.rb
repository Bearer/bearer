Airbrake.notify_request_sync(
  method: 'GET',
  route: "/users/#{user.first_name}",
  status_code: 200,
  timing: 123.45 # ms
)
Airbrake.notify_request(
  method: 'GET',
  route: "/users/#{user.first_name}",
  status_code: 200,
  timing: 123.45 # ms
)


Airbrake.notify_request_sync(
  {
    current_user: current_user.email
  },
  request_id: 123
)
Airbrake.notify_request(
  {
    current_user: current_user.email
  },
  request_id: 123
)


Airbrake.notify_query_sync(
  method: 'GET',
  route: "/users/#{user.first_name}",
  query: 'SELECT * FROM foos'
)
Airbrake.notify_query(
  method: 'GET',
  route: "/users/#{user.first_name}",
  query: 'SELECT * FROM foos'
)


Airbrake.notify_query_sync(
  {
    user: user.email
  },
  request_id: 123
)
Airbrake.notify_query(
  {
    user: user.email
  },
  request_id: 123
)


Airbrake.notify_performance_breakdown_sync(
  method: 'GET',
  route: "/users/#{user.first_name}",
  response_type: 'json',
  groups: { db: 24.0, view: 0.4 }, # ms
  timing: 123.45 # ms
)
Airbrake.notify_performance_breakdown(
  method: 'GET',
  route: "/users/#{user.first_name}",
  response_type: 'json',
  groups: { db: 24.0, view: 0.4 }, # ms
  timing: 123.45 # ms
)


Airbrake.notify_performance_breakdown_sync(
  {
    user: user.email
  },
  request_id: 123
)
Airbrake.notify_performance_breakdown(
  {
    user: user.email
  },
  request_id: 123
)


Airbrake.notify_queue_sync(
  queue: "emails",
  error_count: 1,
  groups: { redis: 24.0, sql: 0.4 }, # ms
  timing: 0.05221 # ms
)
Airbrake.notify_queue(
  queue: "emails",
  error_count: 1,
  groups: { redis: 24.0, sql: 0.4 }, # ms
  timing: 0.05221 # ms
)


Airbrake.notify_queue_sync(
  {
    user: user.email
  },
  job_id: 123
)
Airbrake.notify_queue(
  {
    user: user.email
  },
  job_id: 123
)