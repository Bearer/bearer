const knex = require('knex')({
  client: 'mysql',
})

exports.handler =  async function(event, _context) {
  var cartDetails = knex.select('user.cart_details')
    .from('users')
    .whereRaw('name = '+ event.user.name)

  return cartDetails
}