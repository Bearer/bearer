YAML.load(params[:oops])

Psych.load(params[:oops])

Syck.load(params[:oops])

JSON.load(params[:oops])

Oj.load(params[:oops]) do |json|
end
Oj.object_load(params[:oops])

Marshal.load(params[:oops])
Marshal.restore(params[:oops])
