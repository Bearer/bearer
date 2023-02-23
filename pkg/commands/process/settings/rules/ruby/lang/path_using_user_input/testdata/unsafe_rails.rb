Rails.root.join(params[:oops])

render(partial: params[:oops])
render_to_string({ file: "/templates/#{params[:oops]}" })
