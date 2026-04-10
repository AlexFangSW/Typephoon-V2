```
POST /api/v1/auth/refresh
```

## Request
### Cookie
| Name  | HttpOnly | SameSite | Secure | Max-Age | Path                 | Description             |
| ---   | ---      | ---      | ---    | ---     | ---                  | ---                     |
| TP_RT | v        | Strict   | v      | 5 min   | /                    | Typephoon access token  |


## Response
Refreshed access token and refresh token
### Cookie
| Name  | HttpOnly | SameSite | Secure | Max-Age | Path                 | Description             |
| ---   | ---      | ---      | ---    | ---     | ---                  | ---                     |
| TP_AT | v        | Strict   | v      | 5 min   | /                    | Typephoon access token  |
| TP_RT | v        | Strict   | v      | 30 days | /api/v1/auth/refresh | Typephoon refresh token |
