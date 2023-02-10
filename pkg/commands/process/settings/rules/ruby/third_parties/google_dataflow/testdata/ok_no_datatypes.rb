# Google::Cloud::Dataflow::V1beta3::Snapshots::Client
snapshot_client = Google::Cloud::Dataflow.snapshots
snapshot_client = Google::Cloud::Dataflow::V1beta3::Snapshots::Client.new

# https://cloud.google.com/ruby/docs/reference/google-cloud-dataflow-v1beta3/latest/Google-Cloud-Dataflow-V1beta3-Snapshots-Client#Google__Cloud__Dataflow__V1beta3__Snapshots__Client_get_snapshot_instance_
snapshot = snapshot_client.get_snapshot

snapshot.id = user.id
snapshot.description = "Snapshot taken #{DateTime.now}"
