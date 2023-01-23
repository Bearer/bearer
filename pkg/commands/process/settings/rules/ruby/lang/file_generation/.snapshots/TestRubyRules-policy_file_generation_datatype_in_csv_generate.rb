critical:
    - policy_name: ""
      policy_dsrid: DSR-4
      policy_display_id: ruby_lang_file_generation
      policy_description: Do not write sensitive data to static files.
      line_number: 5
      filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_generate.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: |-
        CSV.generate do |csv|
        	csv << ["email", "first_name", "last_name"]
        	users.each do |user_5|
        		csv << [
        			user_5.email,
        			user_5.first_name,
        			user_5.last_name
        		]
        	end
        end
    - policy_name: ""
      policy_dsrid: DSR-4
      policy_display_id: ruby_lang_file_generation
      policy_description: Do not write sensitive data to static files.
      line_number: 6
      filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_generate.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: |-
        CSV.generate do |csv|
        	csv << ["email", "first_name", "last_name"]
        	users.each do |user_5|
        		csv << [
        			user_5.email,
        			user_5.first_name,
        			user_5.last_name
        		]
        	end
        end
    - policy_name: ""
      policy_dsrid: DSR-4
      policy_display_id: ruby_lang_file_generation
      policy_description: Do not write sensitive data to static files.
      line_number: 7
      filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_generate.rb
      category_groups:
        - PII
      parent_line_number: 1
      parent_content: |-
        CSV.generate do |csv|
        	csv << ["email", "first_name", "last_name"]
        	users.each do |user_5|
        		csv << [
        			user_5.email,
        			user_5.first_name,
        			user_5.last_name
        		]
        	end
        end


--

