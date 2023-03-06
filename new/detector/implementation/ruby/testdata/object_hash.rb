nested = {
  "one" => 42,
  "two" => "hi"
}

call({
  x: { n: nested },
  y: { b: 4 }
})
