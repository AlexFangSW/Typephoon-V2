```
GET /api/v1/auth/login
```
This endpoint triggers OAuth2 workflow.  

## Response
### Cookie
| Name  | HttpOnly | SameSite | Secure | Max-Age | Path                 | Description             |
| ---   | ---      | ---      | ---    | ---     | ---                  | ---                     |
| TP_AT | v        | Strict   | v      | 5 min   | /                    | Typephoon access token  |
| TP_RT | v        | Strict   | v      | 30 days | /api/v1/auth/refresh | Typephoon refresh token |

