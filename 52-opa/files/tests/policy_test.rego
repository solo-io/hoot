package istio.authz

import future.keywords

test_no_token_request_denied if {
	not allow with input as {"attributes": {"request": {"http": {
		"method": "GET",
		"path": "/headers",
	}}}}
}

test_invalid_token_request_denied if {
	not allow with input as {"attributes": {"request": {"http": {
		"headers": {
			"authorization": "Bearer something"
		},
		"method": "GET",
		"path": "/headers",
	}}}}
}

test_expired_token_request_denied if {
	not allow with input as {"attributes": {"request": {"http": {
		"headers": {
			"authorization": "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJzb2xvLmlvIiwiaWF0IjoxNjg3NDY5NzU1LCJleHAiOjE2ODc0NzIyNzcsImF1ZCI6Ind3dy5zb2xvLmlvIiwic3ViIjoicGV0ZXJAc29sby5pbyIsIkdpdmVuTmFtZSI6IlBldGVyIiwicm9sZSI6ImFkbWluIn0.An4P2MfQJD40frOSOMZC0ar-N-R7YjseG5RIJ8EBxn0"
		},
		"method": "GET",
		"path": "/headers",
	}}}}
}

test_guest_get_headers_allowed if {
	allow with input as {"attributes": {"request": {"http": {
		"headers": {
			"authorization": "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJzb2xvLmlvIiwiaWF0IjoxNjg3NDY5NzU1LCJleHAiOjE3MTkwMDc3ODIsImF1ZCI6Ind3dy5zb2xvLmlvIiwic3ViIjoicGF1bEBzb2xvLmlvIiwiR2l2ZW5OYW1lIjoiUGF1bCIsInJvbGUiOiJndWVzdCJ9.JMbwsbPBS6_9wPQtbZ9jVqr3hHme2VUJzYShhhQudnQ"
		},
		"method": "GET",
		"path": "/headers",
	}}}}
}

test_admin_get_headers_denied if {
	not allow with input as {"attributes": {"request": {"http": {
		"headers": {
			"authorization": "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJzb2xvLmlvIiwiaWF0IjoxNjg3NDY5NzU1LCJleHAiOjE3MTkwMDc3ODIsImF1ZCI6Ind3dy5zb2xvLmlvIiwic3ViIjoicGV0ZXJAc29sby5pbyIsIkdpdmVuTmFtZSI6IlBldGVyIiwicm9sZSI6ImFkbWluIn0.8OwUnlJUoW0eBOtA6tK7fBfAGzXOkiCcttwSkmZTVgY"
		},
		"method": "GET",
		"path": "/headers",
	}}}}
}

test_admin_post_without_json_denied if {
	not allow with input as {"attributes": {"request": {"http": {
		"headers": {
			"authorization": "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJzb2xvLmlvIiwiaWF0IjoxNjg3NDY5NzU1LCJleHAiOjE3MTkwMDc3ODIsImF1ZCI6Ind3dy5zb2xvLmlvIiwic3ViIjoicGV0ZXJAc29sby5pbyIsIkdpdmVuTmFtZSI6IlBldGVyIiwicm9sZSI6ImFkbWluIn0.8OwUnlJUoW0eBOtA6tK7fBfAGzXOkiCcttwSkmZTVgY"
		},
		"method": "POST",
		"path": "/post",
	}}}}
}

test_admin_post_invalid_json_denied if {
	not allow with input as {"attributes": {"request": {"http": {
		"headers": {
			"authorization": "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJzb2xvLmlvIiwiaWF0IjoxNjg3NDY5NzU1LCJleHAiOjE3MTkwMDc3ODIsImF1ZCI6Ind3dy5zb2xvLmlvIiwic3ViIjoicGV0ZXJAc29sby5pbyIsIkdpdmVuTmFtZSI6IlBldGVyIiwicm9sZSI6ImFkbWluIn0.8OwUnlJUoW0eBOtA6tK7fBfAGzXOkiCcttwSkmZTVgY"
		},
		"method": "POST",
		"path": "/post",
		"body": "{\"foo\": \"bar\"}"
	}}}}
}

test_admin_post_json_use_number_denied if {
	not allow with input as {"attributes": {"request": {"http": {
		"headers": {
			"authorization": "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJzb2xvLmlvIiwiaWF0IjoxNjg3NDY5NzU1LCJleHAiOjE3MTkwMDc3ODIsImF1ZCI6Ind3dy5zb2xvLmlvIiwic3ViIjoicGV0ZXJAc29sby5pbyIsIkdpdmVuTmFtZSI6IlBldGVyIiwicm9sZSI6ImFkbWluIn0.8OwUnlJUoW0eBOtA6tK7fBfAGzXOkiCcttwSkmZTVgY"
		},
		"method": "POST",
		"path": "/post",
		"body": "{\"id\": 123 }"
	}}}}
}

test_admin_post_valid_json_allowed if {
	allow with input as {"attributes": {"request": {"http": {
		"headers": {
			"authorization": "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJzb2xvLmlvIiwiaWF0IjoxNjg3NDY5NzU1LCJleHAiOjE3MTkwMDc3ODIsImF1ZCI6Ind3dy5zb2xvLmlvIiwic3ViIjoicGV0ZXJAc29sby5pbyIsIkdpdmVuTmFtZSI6IlBldGVyIiwicm9sZSI6ImFkbWluIn0.8OwUnlJUoW0eBOtA6tK7fBfAGzXOkiCcttwSkmZTVgY"
		},
		"method": "POST",
		"path": "/post",
		"body": "{\"id\": \"hello\"}"
	}}}}
}