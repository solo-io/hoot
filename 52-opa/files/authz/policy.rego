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
	v := io.jwt.verify_hs256(token, "qwertyuiopasdfghjklzxcvbnm123456")
	print("is_token_valid result", v)
	v == true
}

# Check whether the token is expired
is_token_expired if {
	# Check the expiration date
	now := time.now_ns() / 1000000000
	now > token_payload.exp
}

# Admin role, if the role in the token payload is admin
is_admin if {
	token_payload.role == "admin"
}

# Guest role, if the role in the token payload is guest
is_guest if {
	token_payload.role == "guest"
}

# Solo audience, if the audience in the token payload is www.solo.io
is_solo_audience if {
	# Check the audience
	token_payload.aud == "www.solo.io"
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
	http_request.method == "POST"
	http_request.path == "/post"
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
	http_request.method == "GET"
	http_request.path == "/headers"
}

# Schema is valid if the body matches the schema
is_schema_valid if {
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
