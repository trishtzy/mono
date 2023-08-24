var jwt = require('jsonwebtoken');
var jwksClient = require('jwks-rsa');
var fs = require('fs');

let jwksFilepath = 'jwk.json';

var client = jwksClient({
  jwksUri: '',
  getKeysInterceptor: () => {
    const file = fs.readFileSync(jwksFilepath, { encoding: 'utf8', flag: 'r' });
    const fileObject = JSON.parse(file)
    return fileObject.keys;
  }
});

function getKey(header, _){
  client.getSigningKey(header.kid, function(err, key) {
    if (err) {
      console.log(err)
      return
    }

    var signingKey = key.publicKey || key.rsaPublicKey;

    jwt.verify(token, signingKey, options, function(err, decoded) {
      if (err) {
        console.log(err)
      } else {
        console.log(decoded)
      }
    });
  });
}

let token = 'eyJhbGciOiJQUzI1NiJ9.eyJoaSI6InRyaWNpYSJ9.DUJh7_xbh8N3OHhxrmjg8NwgNfaOF8jENjfghWS-l2uBQht6FH35DT0pVt-eEGM3TlU50swGkezaSTcOX3LV5VRebmNrvpw02buxde30Tg-69ldjBg85ptZSbzvlat56ICtGInYoTDJ3wPzWlk5NWTqI-P8MT3dBL2_DP49d_m4QDLY40A1a8-Ld9GZ2bIBkJqzPw0nX2TZBhY8hp2ujZL25p3zs5DLg8V562zNngZYzECrVGqK06JEDb3tp45VztrOUeppoVJbiZ7yiij_gwBPa7P8Y_anTl9K0-iivIGL_zaJDpyAXBSPSkhddWObOlw1yXWX7LHH5P6z5qI4b4A';
let options = {
  algorithm: 'RS256'
}

let header = {
  kid: 'my-kid'
}
getKey(header)

