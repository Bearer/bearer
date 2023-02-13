Google::Apis::AnalyticsreportingV4::User.new(user_id: user.email)

Google::Apis::AnalyticsreportingV4::UserActivitySession.update!(
  session_id: DateTime.now + user.ip_address
)