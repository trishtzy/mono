require('jwt')

optional_parameters = { kid: 'my-kid', use: 'sig', alg: 'RS512' }
jwk = JWT::JWK.new(OpenSSL::PKey::RSA.new(2048), optional_parameters)
payload = { data: 'data' }
token = JWT.encode(payload, jwk.signing_key, jwk[:alg], kid: jwk[:kid])

# JSON Web Key Set for advertising your signing keys
jwks_hash = JWT::JWK::Set.new(jwk).export
puts token
puts "\n"
puts jwks_hash.to_json
