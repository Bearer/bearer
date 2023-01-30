data_types:
    - name: Physical Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/digest_md5.rb
              line_number: 1
              field_name: address
              object_name: user
              subject_name: User
risks:
    - detector_id: ruby_lang_weak_encryption_with_data
      data_types:
        - name: Physical Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/digest_md5.rb
              line_number: 1
              parent:
                line_number: 1
                content: Digest::MD5.hexdigest(user.address)
              field_name: address
              object_name: user
              subject_name: User
    - detector_id: ruby_lang_weak_encryption
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption_with_data/testdata/digest_md5.rb
          line_number: 1
          parent:
            line_number: 1
            content: Digest::MD5.hexdigest(user.address)
          content: |
            Digest::MD5.hexdigest()
components: []


--

