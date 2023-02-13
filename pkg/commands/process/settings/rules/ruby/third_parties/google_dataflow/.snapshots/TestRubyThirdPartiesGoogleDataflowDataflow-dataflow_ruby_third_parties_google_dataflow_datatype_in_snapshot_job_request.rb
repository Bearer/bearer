data_types:
    - name: IP address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_snapshot_job_request.rb
              line_number: 2
              field_name: ip_address
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_google_dataflow
      data_types:
        - name: IP address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_snapshot_job_request.rb
              line_number: 2
              parent:
                line_number: 2
                content: 'req.description = "Current user: #{user.ip_address}"'
              field_name: ip_address
              object_name: user
              subject_name: User
    - detector_id: google_dataflow_description_classes
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_snapshot_job_request.rb
          line_number: 1
          parent:
            line_number: 1
            content: Google::Cloud::Dataflow::V1beta3::SnapshotJobRequest.new
          content: |
            Google::Cloud::Dataflow::$<_>::SnapshotJobRequest.new
components: []


--

