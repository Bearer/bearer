response = Net::HTTP.post_form("https://my.api.com/users/search", { email: user.email })

HTTParty.post("https://my.api.com/users/search", body: { email: user.email })

