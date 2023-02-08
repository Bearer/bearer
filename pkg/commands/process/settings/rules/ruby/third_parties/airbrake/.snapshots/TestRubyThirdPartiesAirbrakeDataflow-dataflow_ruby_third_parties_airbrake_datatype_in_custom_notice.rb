data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_custom_notice.rb
              line_number: 5
              field_name: email
              object_name: current_user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_airbrake
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/airbrake/testdata/datatype_in_custom_notice.rb
              line_number: 5
              parent:
                line_number: 4
                content: |-
                    def to_airbrake
                        { params: { user: current_user.email } }
                      end
              field_name: email
              object_name: current_user
              subject_name: User
components: []


--

