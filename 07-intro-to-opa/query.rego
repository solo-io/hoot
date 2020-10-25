package app.api

default allow = false

# make sure that we only allow requests with matching claims:
allow {
    claim := input.claims[_]
    data.claims[claim].methods[_] == input.method

    prefix := data.claims[claim].path_prefix
    prefix == array.slice(input.path, 0, count(prefix))
}

# a reader is someone that uses GET or HEAD request
is_reader {
    input.methods[_] == "GET"
} {
    input.methods[_] == "HEAD"
}

# everyone, including anonymous users can read posts
allow {
    is_reader
    input.path[0] == "posts"
}
