risks:
    - detector_id: google_dataflow_value_classes
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/google_dataflow/testdata/ok_ignored_datatypes_only.rb
          line_number: 1
          parent:
            line_number: 1
            content: Google::Cloud::Dataflow::V1beta3::StructuredMessage::Parameter.new
          content: |
            Google::Cloud::Dataflow::$<_>::StructuredMessage::Parameter.new
components: []


--

