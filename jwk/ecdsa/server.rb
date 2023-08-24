require('jwt')

payload = {document_id: 1}

# JWT supports 3 ECDSA Curves - ES256, ES384, ES512
# See: https://datatracker.ietf.org/doc/html/rfc7518#page-9
optional_parameters = { kid: 'my-kid', use: 'sig', alg: 'ES256' }
jwk = JWT::JWK.new(OpenSSL::PKey::EC.generate('prime256v1'), optional_parameters)

# Encoding
token2 = JWT.encode(payload, jwk.signing_key, jwk[:alg], kid: jwk[:kid])

# JSON Web Key Set for advertising your signing keys
jwks_hash = JWT::JWK::Set.new(jwk).export

puts "token2: #{token2}"
puts "jwks_hash: #{jwks_hash.to_json}"
