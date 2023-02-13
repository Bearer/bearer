# Google::Cloud::Dataflow::V1beta3::Messages::Client and others
client = Google::Cloud::Dataflow.messages
client = Google::Cloud::Dataflow::V1beta3::Messages::Client.new

# client config
client_config = client.configure
client.configure do |config|
  config.metadata = { current_user_id: current_user.email }
end

# config metadata
# https://cloud.google.com/ruby/docs/reference/google-cloud-dataflow-v1beta3/latest/Google-Cloud-Dataflow-V1beta3-Messages-Client-Configuration#Google__Cloud__Dataflow__V1beta3__Messages__Client__Configuration_metadata_instance_
client_config.metadata("current_user_id": current_user.email)
client_config.metadata = { current_user_id: current_user.email }
