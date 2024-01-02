# Authroization Code Grant
Authorization Code Grant is one of grant types defined in [RFC 6749](https://datatracker.ietf.org/doc/html/rfc6749#section-4.1). This grant is used to obtain both access tokens and refresh tokens.

## Glossaries
- [Access token](https://datatracker.ietf.org/doc/html/rfc6749#section-1.4)
  - credentials used to access protected resources.
- [Refresh token](https://datatracker.ietf.org/doc/html/rfc6749#section-1.5)
  - credentials used to obtain access tokens.
  - it is not used for authentication/authorization directly.
- [User Agent](https://datatracker.ietf.org/doc/html/rfc2616#section-1.3)
  - The client which initiates a request. These are often browsers, editors, spiders(web-traversing robots), or other end user tools.

## Flow
For more details, please refer [the seciont-4.1](https://datatracker.ietf.org/doc/html/rfc6749#section-4.1).

```mermaid
sequenceDiagram
    actor ResourceOwner
    Client->>UserAgent: 1. do an authroization request
    UserAgent->>AuthorizationServer: 2. send the authorization request
    AuthorizationServer->>UserAgent: 3. request a user authentication
    UserAgent->>ResourceOwner: 4. request the user authentication
    ResourceOwner->>UserAgent: 5. provide the information for the user authentication
    UserAgent->>AuthorizationServer: 6. send the information for the user authentication
    AuthorizationServer->>Client: 7. return the authorization response
    Client->>AuthorizationServer: 8. send a token issue request
    AuthorizationServer->>Client: 9. return an access token issue response
```

### Authorization Request/Response
#### [Request](https://datatracker.ietf.org/doc/html/rfc6749#section-4.1.1)

- required
  - `response_type`: must be set to "code".
  - `client_id`
- recommended paramters
  - `state`
- optional
  - `redirect_uri`
  - `scope`

example

```plaintext
GET /authorize?response_type=code&client_id=s6BhdRkqt3&state=xyz&redirect_uri=https%3A%2F%2Fclient%2Eexample%2Ecom%2Fcb HTTP/1.1
Host: server.example.com
```

#### [Response](https://datatracker.ietf.org/doc/html/rfc6749#section-4.1.2)
- required
  - `code`
  - `state`

example

```plaintext
HTTP/1.1 302 Found
Location: https://client.example.com/cb?code=SplxlOBeZQQYbYS6WxSbIA&state=xyz
```

#### [Error Response](https://datatracker.ietf.org/doc/html/rfc6749#section-4.1.2.1)
- required
  - `error`
  - `state`
- optional
  - `error_description`
  - `error_uri`

example

```plaintext
HTTP/1.1 302 Found
Location: https://client.example.com/cb?error=access_denied&state=xyz
```

### Token Issue Request/Response
#### [Request](https://datatracker.ietf.org/doc/html/rfc6749#section-4.1.3)

- required
  - `grant_type`: must be set to "authorization_code".
  - `code`
  - `redirect_uri`
  - `client_id`

example

```plaintext
POST /token HTTP/1.1
Host: server.example.com
Authorization: Basic czZCaGRSa3F0MzpnWDFmQmF0M2JW
Content-Type: application/x-www-form-urlencoded
grant_type=authorization_code&code=SplxlOBeZQQYbYS6WxSbIA&redirect_uri=https%3A%2F%2Fclient%2Eexample%2Ecom%2Fcb
```

#### [Response](https://datatracker.ietf.org/doc/html/rfc6749#section-4.1.4), [Successful Response](https://datatracker.ietf.org/doc/html/rfc6749#section-5.1)

- required
  - `access_token`
  - `token_type`
- recommended
  - `expires_in`
- optional
  - `refresh_token`
  - `scope`

example

```plaintext
HTTP/1.1 200 OK
Content-Type: application/json;charset=UTF-8
Cache-Control: no-store
Pragma: no-cache

{
    "access_token":"2YotnFZFEjr1zCsicMWpAA",
    "token_type":"example",
    "expires_in":3600,
    "refresh_token":"tGzv3JOkF0XG5Qx2TlKWIA",
    "example_parameter":"example_value"
}
```

#### [Error Response](https://datatracker.ietf.org/doc/html/rfc6749#section-5.2)

- required
  - `error`
- optional
  - `error_description`
  - `error_uri`

example

```plaintext
HTTP/1.1 400 Bad Request
Content-Type: application/json;charset=UTF-8
Cache-Control: no-store
Pragma: no-cache

{
    "error":"invalid_request"
}
```
