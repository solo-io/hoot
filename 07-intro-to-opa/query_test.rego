package app.api

test_post_allowed_by_post_writer {
    allow with input as {
        "method": "POST",
        "path": ["posts","alice","hello_post"],
        "sub": "alice",
        "claims": [
            "write:posts",
        ]
    }
}

test_post_not_allowed_by_review_writer{
    not allow with input as {
        "method": "POST",
        "path": ["posts","alice","hello_post"],
        "sub": "alice",
        "claims": [
            "write:reviews"
        ]
    }
}
test_can_read_with_no_claims{
    not allow with input as {
        "method": "GET",
        "path": ["posts","alice","hello_post"],
    }
}