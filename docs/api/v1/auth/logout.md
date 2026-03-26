```
POST /api/v1/auth/logout
```

## Request
### Cookie
| Name  | HttpOnly | SameSite | Secure | Max-Age | Path                 | Description             |
| ---   | ---      | ---      | ---    | ---     | ---                  | ---                     |
| TP_AT | v        | Strict   | v      | 5 min   | /                    | Typephoon access token  |


## Response
Removes access (TP_AT) and refresh token (TP_RT) from cookie.  
The refresh token is invalidated.  
