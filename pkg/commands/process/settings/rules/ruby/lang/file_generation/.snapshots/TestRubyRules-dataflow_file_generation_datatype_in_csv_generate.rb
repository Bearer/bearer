data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_generate.rb
              line_number: 5
              field_name: email
              object_name: user_5
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_generate.rb
              line_number: 6
              field_name: first_name
              object_name: user_5
    - name: Lastname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_generate.rb
              line_number: 7
              field_name: last_name
              object_name: user_5
risks:
    - detector_id: ruby_lang_file_generation
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_generate.rb
              line_number: 5
              parent:
                line_number: 1
                content: |-
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
              field_name: email
              object_name: user_5
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_generate.rb
              line_number: 6
              parent:
                line_number: 1
                content: |-
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
              field_name: first_name
              object_name: user_5
        - name: Lastname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_generate.rb
              line_number: 7
              parent:
                line_number: 1
                content: |-
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
              field_name: last_name
              object_name: user_5
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_generate.rb
              line_number: 5
              parent:
                line_number: 1
                content: |-
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
              object_name: user_5
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_generate.rb
              line_number: 6
              parent:
                line_number: 1
                content: |-
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
              object_name: user_5
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_generate.rb
              line_number: 7
              parent:
                line_number: 1
                content: |-
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
              object_name: user_5
components: []


--

