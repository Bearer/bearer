RubyVM::InstructionSequence.compile(request.env["oops"])

a.eval(request.env["oops"], "test")

a.instance_eval(request.env["oops"])

a.class_eval(request.env["oops"])

a.module_eval(request.env["oops"])

eval(request.env["oops"])

instance_eval(request.env["oops"], "test")

class_eval(request.env["oops"])

module_eval(request.env["oops"])
