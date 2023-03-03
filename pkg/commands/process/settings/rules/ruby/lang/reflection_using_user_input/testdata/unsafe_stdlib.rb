Object.const_get(params[:class])
Object.const_set(params[:class], 42)
Object.remove_const(params[:class])

method(params[:method])

x.define_method(params[:method]) {}

params[:method].to_sym.to_proc

bad_things(&params[:method].to_sym)
x.bad_things(&params[:method].to_sym)
