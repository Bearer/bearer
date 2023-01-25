data_types:
    - name: Email Address
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_open.rb
              line_number: 5
              field_name: email
              object_name: user
    - name: Firstname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_open.rb
              line_number: 6
              field_name: first_name
              object_name: user
    - name: Lastname
      detectors:
        - name: ruby
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_open.rb
              line_number: 7
              field_name: last_name
              object_name: user
risks:
    - detector_id: ruby_lang_file_generation
      data_types:
        - name: Email Address
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_open.rb
              line_number: 5
              parent:
                line_number: 1
                content: |-
                    CSV.open("path/to/user.csv", "wb") do |csv|
                      csv << ["email", "first_name", "last_name"]
                    	users.each do |user|
                    		csv << [
                    			user.email,
                    			user.first_name,
                    			user.last_name
                    		]
                    	end
                    end
              field_name: email
              object_name: user
        - name: Firstname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_open.rb
              line_number: 6
              parent:
                line_number: 1
                content: |-
                    CSV.open("path/to/user.csv", "wb") do |csv|
                      csv << ["email", "first_name", "last_name"]
                    	users.each do |user|
                    		csv << [
                    			user.email,
                    			user.first_name,
                    			user.last_name
                    		]
                    	end
                    end
              field_name: first_name
              object_name: user
        - name: Lastname
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_open.rb
              line_number: 7
              parent:
                line_number: 1
                content: |-
                    CSV.open("path/to/user.csv", "wb") do |csv|
                      csv << ["email", "first_name", "last_name"]
                    	users.each do |user|
                    		csv << [
                    			user.email,
                    			user.first_name,
                    			user.last_name
                    		]
                    	end
                    end
              field_name: last_name
              object_name: user
        - name: Unique Identifier
          stored: false
          locations:
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_open.rb
              line_number: 5
              parent:
                line_number: 1
                content: |-
                    CSV.open("path/to/user.csv", "wb") do |csv|
                      csv << ["email", "first_name", "last_name"]
                    	users.each do |user|
                    		csv << [
                    			user.email,
                    			user.first_name,
                    			user.last_name
                    		]
                    	end
                    end
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_open.rb
              line_number: 6
              parent:
                line_number: 1
                content: |-
                    CSV.open("path/to/user.csv", "wb") do |csv|
                      csv << ["email", "first_name", "last_name"]
                    	users.each do |user|
                    		csv << [
                    			user.email,
                    			user.first_name,
                    			user.last_name
                    		]
                    	end
                    end
              object_name: user
            - filename: pkg/commands/process/settings/rules/ruby/lang/file_generation/testdata/datatype_in_csv_open.rb
              line_number: 7
              parent:
                line_number: 1
                content: |-
                    CSV.open("path/to/user.csv", "wb") do |csv|
                      csv << ["email", "first_name", "last_name"]
                    	users.each do |user|
                    		csv << [
                    			user.email,
                    			user.first_name,
                    			user.last_name
                    		]
                    	end
                    end
              object_name: user
components: []


--

