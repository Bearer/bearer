analytics = Segment::Analytics.new(write_key: "ABC123F")
analytics.track(user_id: user.id, event: "account sign in", context: { ip: user.ip_address })