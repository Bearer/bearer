risks:
    - detector_id: ruby_third_parties_algolia_client
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/algolia/testdata/ok_no_datatype.rb
          line_number: 1
          parent:
            line_number: 1
            content: Algolia::Search::Client.create('YourApplicationID', 'YourWriteAPIKey')
          content: |
            Algolia::Search::Client.create()
    - detector_id: ruby_third_parties_algolia_index
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/third_parties/algolia/testdata/ok_no_datatype.rb
          line_number: 2
          parent:
            line_number: 2
            content: client.init_index("my_index")
          content: |
            $<CLIENT>.init_index()
components: []


--

