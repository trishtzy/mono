var fs = require('fs');
var jose = require('node-jose');
var jwt = require('jsonwebtoken');

let kid = 'my-kid';
let jwksFile = fs.readFileSync('jwk.json', { encoding: 'utf8', flag: 'r' });
let input = JSON.parse(jwksFile);

// Replace token value here
let token = 'eyJraWQiOiJteS1raWQiLCJhbGciOiJFUzI1NiJ9.eyJkb2N1bWVudF9pZCI6MX0.MlRZeWIr_lFK1vyfMmmcDIxDm910Lk9xTqV52vPNEyx8j8zYfGowVe2Hj-5H5Iljnkp0MZSjRs6pwU8S86D54g';
let options = {
  algorithm: 'ES256'
}

// {input} is a String or JSON object representing the JWK-set
jose.JWK.asKeyStore(input).
     then(function(result) {
       // {result} is a jose.JWK.KeyStore
       keystore = result;
       // by 'kid'
       let key = keystore.get(kid);
       let pemKey = key.toPEM();
       jwt.verify(token, pemKey, options, function(err, decoded) {
        if (err) {
          console.log(err)
        } else {
          console.log(decoded)
        }
       })
     });
