risks:
    - detector_id: ruby_lang_weak_encryption
      locations:
        - filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/digest_sha1.rb
          line_number: 1
          parent:
            line_number: 1
            content: Digest::SHA1.hexidigest("hello world")
          content: |
            Digest::SHA1.hexidigest()
components: []


--

