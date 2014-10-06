// To calculate a salted password
c = require('./models/cryptoUtil');
console.log( c.createHash('pwd', 'salt')); 

// To update the password in db
// via the MongoDB shell
mongo mongoUri -u mongoUser -p mongoPwd
use dbname
db.users.update({_id:ObjectId("user_id")}, {$set: {password : "salted_password"}})
