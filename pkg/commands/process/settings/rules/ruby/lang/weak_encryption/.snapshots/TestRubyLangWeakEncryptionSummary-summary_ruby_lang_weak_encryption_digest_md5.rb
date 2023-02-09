low:
    - rule_dsrid: DSR-7
      rule_display_id: ruby_lang_weak_encryption
      rule_description: Avoid weak encryption libraries.
      rule_documentation_url: https://curio.sh/reference/rules/ruby_lang_weak_encryption
      line_number: 1
      filename: pkg/commands/process/settings/rules/ruby/lang/weak_encryption/testdata/digest_md5.rb
      parent_line_number: 1
      parent_content: Digest::MD5.hexdigest("hello world")


--

