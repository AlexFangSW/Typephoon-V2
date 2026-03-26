```
POST /api/v1/match/queue-in
```
## Request
### Cookie
| Name  | Type   | Default | Example       | Required | Description                                                |
| ---   | ---    | ---     | ---           | ---      | ---                                                        |
| TP_AT | string | N/A     | `xxx.xxx.xxx` |          | Access token, if not provided, will be identified as guest |

## Response
### 200
#### Body
| Name        | Type   | Example       | Nullable | Description                            |
| ---         | ---    | ---           | ---      | ---                                    |
| match_token | string | `xxx.xxx.xxx` |          | Token used for game connection         |

#### Cookie
| Name  | Type   | Example       | Nullable | Description                       |
| ---   | ---    | ---           | ---      | ---                               |
| TP_AT | string | `xxx.xxx.xxx` | v        | Access token generated for guests |


## Example
```bash
curl -XGET http://localhost:8080/api/v1/match/queue-in \
    -H 'Content-Type: application/json' \
    --cookie 'TP_AT=xxx.xxx.xxx'
```
```json
status: 200
{
    "match_token": "xxx.xxx.xxx",
}
```
