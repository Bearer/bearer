render html: x
render inline: "ok"

render html: sanitize(params[:oops])
render inline: "<h1>#{strip_tags(params[:oops])}</h1>"
