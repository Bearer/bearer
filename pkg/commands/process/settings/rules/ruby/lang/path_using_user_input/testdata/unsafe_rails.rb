Rails.root.join(params[:oops])

render(partial: params[:oops])
render_to_string({ file: "/templates/#{params[:oops]}" })

send_file params[:oops], type: "text/html"
