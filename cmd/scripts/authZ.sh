#/bin/bash

curl "http://localhost:9001/authorize?response_type=code&scope=hoge&client_id=oauth-client-1&redirect_uri=http%3A%2F%2Flocalhost%3A9000%2Fcallback"