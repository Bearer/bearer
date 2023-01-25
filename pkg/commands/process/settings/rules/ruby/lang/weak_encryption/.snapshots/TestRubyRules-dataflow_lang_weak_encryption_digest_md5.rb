risks:
    - detector_id: ruby_lang_weak_encryption
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/digest_md5.rb
          line_number: 1
          parent:
            line_number: 1
            content: Digest::MD5.hexdigest("hello world")
          content: |
            Digest::MD5.hexdigest()
components: []


--

