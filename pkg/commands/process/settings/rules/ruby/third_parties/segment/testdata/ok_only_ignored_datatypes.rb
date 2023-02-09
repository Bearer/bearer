analytics = Segment::Analytics.new(write_key: "ABC123F")
analytics.track(user_id: user.id, event: "order created")