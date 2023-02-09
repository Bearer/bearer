analytics = Segment::Analytics.new(write_key: "ABC123F")
analytics.alias(user_id: user.email, previous_id: "some id")