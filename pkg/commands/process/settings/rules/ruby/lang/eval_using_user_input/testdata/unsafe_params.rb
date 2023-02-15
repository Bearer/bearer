RubyVM::InstructionSequence.compile(params["oops"])

a.eval(params["oops"], "test")

a.instance_eval(params["oops"])

a.class_eval(params["oops"])

a.module_eval(params["oops"])

eval(params["oops"])

instance_eval(params["oops"], "test")

class_eval(params["oops"])

module_eval(params["oops"])
