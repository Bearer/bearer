def handler(event:, context:)
  RubyVM::InstructionSequence.compile(event["oops"])

  a.eval(event["oops"], "test")

  a.instance_eval(event["oops"])

  a.class_eval(event["oops"])

  a.module_eval(event["oops"])

  eval(event["oops"])

  instance_eval(event["oops"], "test")

  class_eval(event["oops"])

  module_eval(event["oops"])
end
