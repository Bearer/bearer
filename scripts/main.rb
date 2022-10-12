require "json"

file = File.read("./recipes.json")

JSON.parse(file).each do |recipe|
  puts "processing #{recipe["name"]}"
  File.open("pkg/classification/db/recipes/#{recipe["name"].downcase.gsub(/(\s|\/|\(|\)|\.|\,|\?)/, "_").gsub(/_{2,}/, "_").gsub(/_$/, "")}.json", "w") do |recipe_file|
    recipe_file << {metadata: {version: "1.0"}}.merge(recipe).to_json
  end
end