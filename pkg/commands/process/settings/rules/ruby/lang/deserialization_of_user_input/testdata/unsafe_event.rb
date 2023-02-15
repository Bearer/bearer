def handler(event:, context:)
  YAML.load(event["oops"])

  Psych.load(event["oops"])

  Syck.load(event["oops"])

  JSON.load(event["oops"])

  Oj.load(event["oops"])
  Oj.object_load(event["oops"]) do |json|
  end

  Marshal.load(event["oops"])
  Marshal.restore(event["oops"])
end
