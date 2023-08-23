let user = { a: 123 }
let nested = { ...user, b: 456 }

call({
  x: { n: nested },
  y: { c: 4 },
})
