event = not_from_handler

YAML.load(event[:ok])

Psych.load(x)

Syck.load("--")

JSON.load(event[:ok])

Oj.load(x)
Oj.object_load("{}") do |json|
end

Marshal.load(x)
Marshal.restore(event[:ok])
