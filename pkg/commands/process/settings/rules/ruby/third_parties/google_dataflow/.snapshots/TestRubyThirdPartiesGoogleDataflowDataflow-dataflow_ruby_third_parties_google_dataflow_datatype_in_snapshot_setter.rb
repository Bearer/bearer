data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_snapshot_setter.rb
              line_number: 9
              field_name: email_address
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_google_dataflow
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_snapshot_setter.rb
              line_number: 9
              parent:
                line_number: 9
                content: 'snapshot.description = "Current user: #{user.email_address}"'
              field_name: email_address
              object_name: user
              subject_name: User
    - detector_id: google_dataflow_client_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_snapshot_setter.rb
          line_number: 2
          parent:
            line_number: 2
            content: Google::Cloud::Dataflow.snapshots
          content: |
            Google::Cloud::Dataflow.$<METHOD>
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_snapshot_setter.rb
          line_number: 3
          parent:
            line_number: 3
            content: Google::Cloud::Dataflow::V1beta3::Snapshots::Client.new
          content: |
            Google::Cloud::Dataflow::$<_>::$<_>::Client.new
    - detector_id: google_dataflow_description_classes
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_snapshot_setter.rb
          line_number: 6
          parent:
            line_number: 6
            content: snapshot_client.get_snapshot
          content: |
            $<VAR>.get_snapshot
    - detector_id: google_dataflow_snapshots_client_init
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_snapshot_setter.rb
          line_number: 2
          parent:
            line_number: 2
            content: Google::Cloud::Dataflow.snapshots
          content: |
            Google::Cloud::Dataflow.snapshots
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/datatype_in_snapshot_setter.rb
          line_number: 3
          parent:
            line_number: 3
            content: Google::Cloud::Dataflow::V1beta3::Snapshots::Client.new
          content: |
            Google::Cloud::Dataflow::$<_>::Snapshots::Client.new
components: []


--

