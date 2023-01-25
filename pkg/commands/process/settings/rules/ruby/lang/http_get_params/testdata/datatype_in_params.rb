uri = URI("https://my.api.com/users/search?ethnic_origin=#{user.ethnic_origin}")

RestClient.get("https://my.api.com/users/search?first_name=#{user.first_name}")