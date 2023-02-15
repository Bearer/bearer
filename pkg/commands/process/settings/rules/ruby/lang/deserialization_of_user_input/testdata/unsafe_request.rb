YAML.load(request.env[:oops])

Psych.load(request.env[:oops])

Syck.load(request.env[:oops])

JSON.load(request.env[:oops])

Oj.load(request.env[:oops])
Oj.object_load(request.env[:oops]) do |json|
end

Marshal.load(request.env[:oops])
Marshal.restore(request.env[:oops])
