import requests

# 1. client endpointを


def authorize_req():
    # GET /authorize?response_type=code&client_id=s6BhdRkqt3&state=xyz&redirect_uri=https%3A%2F%2Fclient%2Eexample%2Ecom%2Fcb HTTP/1.1
    # Host: server.example.com
    params = {
        "response_type": "code",
        "client_id": "test_client",
        "state": "test_state",
        "redirect_uri": "http://localhost:9000/callback",
        "scope": "read write",
    }
    resp = requests.get("http://localhost:8080/authorize", params=params)
    print(resp.text)


def main():
    authorize_req()


if __name__ == "__main__":
    main()
