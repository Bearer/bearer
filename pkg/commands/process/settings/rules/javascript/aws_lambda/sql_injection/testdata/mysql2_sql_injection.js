const connection = mysql.createConnection({});
const asyncConn = await mysql.createConnection({});

exports.handler = async function(event, _context) {
  connection.query("SELECT * FROM `user` WHERE name = " + event.customer.name);

  await asyncConn.execute("SELECT * FROM `admin_users` WHERE ID = " + event.admin.id)

  // pool query
  var pool = mysql.createPool()
  pool.query("SELECT * FROM users WHERE name = " + event.user_name, function() {
    // do something
  })
  pool.getConnection(function(_err, conn) {
    conn.query("SELECT * FROM users WHERE name = " + event.user_name, function() {
      // do something
    })
    pool.releaseConnection(conn)
  })
}