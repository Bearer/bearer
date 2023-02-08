data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/algolia/testdata/datatype_in_save_object.rb
              line_number: 4
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/algolia/testdata/datatype_in_save_object.rb
              line_number: 6
              field_name: email
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_third_parties_algolia
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/algolia/testdata/datatype_in_save_object.rb
              line_number: 4
              parent:
                line_number: 4
                content: 'index.save_object({ email: user.email }, { auto_generate_object_id_if_not_exist: true })'
              field_name: email
              object_name: user
              subject_name: User
            - filename: pkg/commands/process/settings/rules/ruby/third_parties/algolia/testdata/datatype_in_save_object.rb
              line_number: 6
              parent:
                line_number: 6
                content: 'index.save_objects([{ email: user.email }], { auto_generate_object_id_if_not_exist: true })'
              field_name: email
              object_name: user
              subject_name: User
    - detector_id: ruby_third_parties_algolia_client
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/algolia/testdata/datatype_in_save_object.rb
          line_number: 1
          parent:
            line_number: 1
            content: Algolia::Search::Client.create('YourApplicationID', 'YourWriteAPIKey')
          content: |
            Algolia::Search::Client.create()
    - detector_id: ruby_third_parties_algolia_index
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/algolia/testdata/datatype_in_save_object.rb
          line_number: 2
          parent:
            line_number: 2
            content: client.init_index("my_index")
          content: |
            $<CLIENT>.init_index()
components: []


--

