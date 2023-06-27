package istio.authz

import future.keywords
import input.attributes.request.http as http_request

default allow := false

# Allow if token is valid, not expired, and if the action is allowed
allow if {
	is_token_valid
	not is_token_expired
	action_allowed
}

# Token is valid if the signature is valid
is_token_valid if {
	print("entering is_token_valid")
	# Verify the signature (note the secret shouldn't be hardcoded like this!)
	print("is_token_valid result", io.jwt.verify_hs256(token, "qwertyuiopasdfghjklzxcvbnm123456"))
	v := io.jwt.verify_hs256(token, "qwertyuiopasdfghjklzxcvbnm123456")
	v == true
}

# Check whether the token is expired
is_token_expired if {
	print("entering is_token_expired")
	# Check the expiration date
	now := time.now_ns() / 1000000000
	exp := now > token_payload.exp
	print("is_token_expired result", exp)
	exp == true
}

# Admin role, if the role in the token payload is admin
is_admin if {
	a := token_payload.role == "admin"
	print("is_admin result", a)
	a == true
}

# Guest role, if the role in the token payload is guest
is_guest if {
	g := token_payload.role == "guest"
	print("is_guest result", g)
	g == true
}

# Solo audience, if the audience in the token payload is www.solo.io
is_solo_audience if {
	# Check the audience
	aud := token_payload.aud == "www.solo.io"
	print("is_solo_audience result", aud)
	aud == true
}

# Action is allowed if:
# - aud is set to www.solo.io and
# - role is set to admin and
# - method is POST and
# - path is /post and
# - schema is valid
action_allowed if {
	is_solo_audience
	is_admin
	m := http_request.method == "POST"
	print("action_allowed result (method == POST)", m)
	p := http_request.path == "/post"
	print("action_allowed result (path == /post)", p)
	is_schema_valid
}

# Action is allowed if:
# - aud is set to www.solo.io and
# - role is set to guest and
# - method is GET and
# - path is /headers
action_allowed if {
	is_solo_audience
	is_guest
	m := http_request.method == "GET"
	print("action_allowed result (method == GET)", m)
	p := http_request.path == "/headers"
	print("action_allowed result (path == /headers)", p)
}

# Schema is valid if the body matches the schema
is_schema_valid if {
	print("entering is_schema_valid")
	print("is_schema_valid result", json.match_schema(http_request.body, schema))
	[match, _] := json.match_schema(http_request.body, schema)
	match == true
}

# Decode the token and return the payload
token_payload := payload if {
	[_, payload, _] := io.jwt.decode(token)
}

# GEt the token from the Authorization header
token := t if {
	# "Authorization": "Bearer <token>"
	t := split(http_request.headers.authorization, " ")[1]
}

# Schema to validate the request body against
schema := {
	"properties": {"id": {"type": "string"}},
	"required": ["id"],
}
