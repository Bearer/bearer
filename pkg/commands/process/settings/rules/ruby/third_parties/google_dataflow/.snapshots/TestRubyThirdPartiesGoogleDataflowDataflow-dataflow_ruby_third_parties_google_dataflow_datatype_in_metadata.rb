data_types:
    - name: IP address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_metadata.rb
              line_number: 2
              field_name: ip_address
              object_name: customer
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_metadata.rb
              line_number: 5
              field_name: ip_address
              object_name: customer
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_metadata.rb
              line_number: 6
              field_name: ip_address
              object_name: customer
              subject_name: User
risks:
    - detector_id: ruby_third_parties_google_dataflow
      data_types:
        - name: IP address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_metadata.rb
              line_number: 2
              parent:
                line_number: 2
                content: 'custom_metadata.value = "ip: #{customer.ip_address}"'
              field_name: ip_address
              object_name: customer
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_metadata.rb
              line_number: 5
              parent:
                line_number: 5
                content: 'template_metadata.description ="ip: #{customer.ip_address}"'
              field_name: ip_address
              object_name: customer
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_metadata.rb
              line_number: 6
              parent:
                line_number: 6
                content: 'template_metadata.name ="ip: #{customer.ip_address}"'
              field_name: ip_address
              object_name: customer
              subject_name: User
    - detector_id: google_dataflow_template_metadata_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_metadata.rb
          line_number: 4
          parent:
            line_number: 4
            content: Google::Cloud::Dataflow::V1beta3::TemplateMetadata.new
          content: |
            Google::Cloud::Dataflow::$<_>::TemplateMetadata.new
    - detector_id: google_dataflow_value_classes
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_metadata.rb
          line_number: 1
          parent:
            line_number: 1
            content: Google::Cloud::Dataflow::V1beta3::ParameterMetadata::CustomMetadataEntry.new
          content: |
            Google::Cloud::Dataflow::$<_>::ParameterMetadata::CustomMetadataEntry.new
components: []


--

