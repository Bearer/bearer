const algoliasearch = require("algoliasearch")
const myAlgolia = algoliasearch("123", "123")

const index = myAlgolia.initIndex("test_index")

// saveObject
const userObj = { user_id: user.ip_address }
index
  .saveObject(userObj, { autoGenerateObjectIDIfNotExist: true })
  .then(console.log("obj saved"))

index.saveObjects([{ email: user.email }])




