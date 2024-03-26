curl \
    -H "Authorization: Bearer dummy_token:1" \
    "http://localhost:8000/authorize"

curl \
    -H "Authorization: Bearer dummy_token:1" \
    -v \
    "http://localhost:8000/authorize"