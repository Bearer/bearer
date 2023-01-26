data_types:
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/digest_sha1.rb
              line_number: 1
              field_name: first_name
              object_name: user
risks:
    - detector_id: ruby_lang_weak_encryption_with_data
      data_types:
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/digest_sha1.rb
              line_number: 1
              parent:
                line_number: 1
                content: Digest::SHA1.hexidigest(user.first_name)
              field_name: first_name
              object_name: user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/digest_sha1.rb
              line_number: 1
              parent:
                line_number: 1
                content: Digest::SHA1.hexidigest(user.first_name)
              object_name: user
    - detector_id: ruby_lang_weak_encryption
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/digest_sha1.rb
          line_number: 1
          parent:
            line_number: 1
            content: Digest::SHA1.hexidigest(user.first_name)
          content: |
            Digest::SHA1.hexidigest()
components: []


--

