custom_metadata = Google::Cloud::Dataflow::V1beta3::ParameterMetadata::CustomMetadataEntry.new
custom_metadata.value = "ip: #{customer.ip_address}"

template_metadata = Google::Cloud::Dataflow::V1beta3::TemplateMetadata.new
template_metadata.description ="ip: #{customer.ip_address}"
template_metadata.name ="ip: #{customer.ip_address}"
