Created from: http://jwtbuilder.jamiekurtz.com/


### Peter (admin)

```console
eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJzb2xvLmlvIiwiaWF0IjoxNjg3NDY5NzU1LCJleHAiOjE3MTkwMDc3ODIsImF1ZCI6Ind3dy5zb2xvLmlvIiwic3ViIjoicGV0ZXJAc29sby5pbyIsIkdpdmVuTmFtZSI6IlBldGVyIiwicm9sZSI6ImFkbWluIn0.8OwUnlJUoW0eBOtA6tK7fBfAGzXOkiCcttwSkmZTVgY
```

```json
{
    "iss": "solo.io",
    "iat": 1687469755,
    "exp": 1719007782,
    "aud": "www.solo.io",
    "sub": "peter@solo.io",
    "GivenName": "Peter",
    "role": "admin"
}
```

### Paul (non-admin)

```console
eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJzb2xvLmlvIiwiaWF0IjoxNjg3NDY5NzU1LCJleHAiOjE3MTkwMDc3ODIsImF1ZCI6Ind3dy5zb2xvLmlvIiwic3ViIjoicGF1bEBzb2xvLmlvIiwiR2l2ZW5OYW1lIjoiUGF1bCIsInJvbGUiOiJndWVzdCJ9.JMbwsbPBS6_9wPQtbZ9jVqr3hHme2VUJzYShhhQudnQ
```

```json
{
    "iss": "solo.io",
    "iat": 1687469755,
    "exp": 1719007782,
    "aud": "www.solo.io",
    "sub": "paul@solo.io",
    "GivenName": "Paul",
    "role": "guest"
}
```

### Expired admin token

```
eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJzb2xvLmlvIiwiaWF0IjoxNjg3NDY5NzU1LCJleHAiOjE2ODc0NzIyNzcsImF1ZCI6Ind3dy5zb2xvLmlvIiwic3ViIjoicGV0ZXJAc29sby5pbyIsIkdpdmVuTmFtZSI6IlBldGVyIiwicm9sZSI6ImFkbWluIn0.An4P2MfQJD40frOSOMZC0ar-N-R7YjseG5RIJ8EBxn0
```

```json
{
    "iss": "solo.io",
    "iat": 1687469755,
    "exp": 1687472277,
    "aud": "www.solo.io",
    "sub": "peter@solo.io",
    "GivenName": "Peter",
    "role": "admin"
}
```

### Non-solo token


```
eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJzb2xvLmlvIiwiaWF0IjoxNjg3NDY5NzU1LCJleHAiOjE3MTkwMDgzMTgsImF1ZCI6Ind3dy5ibGFoLmlvIiwic3ViIjoicGV0ZXJAc29sby5pbyIsIkdpdmVuTmFtZSI6IlBldGVyIiwicm9sZSI6ImFkbWluIn0.MmYP_VhcihkusQXTS6hD1oNET0Pxj4HfmohbOH6v0zo
```

```json
{
    "iss": "solo.io",
    "iat": 1687469755,
    "exp": 1719008318,
    "aud": "www.blah.io",
    "sub": "peter@solo.io",
    "GivenName": "Peter",
    "role": "admin"
}
```
