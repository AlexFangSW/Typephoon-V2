```
POST /api/v1/match/queue-in
```
## Request
### Header
| Name          | Type   | Default | Example              | Required | Description                                                     |
| ---           | ---    | ---     | ---                  | ---      | ---                                                             |
| Authorization | string | N/A     | `Bearer xxx.xxx.xxx` |          | Auth bearer token, if not provided, will be identified as guest |

## Response
### 200
#### Body
| Name        | Type   | Example       | Nullable | Description                            |
| ---         | ---    | ---           | ---      | ---                                    |
| guest_token | string | `xxx.xxx.xxx` | v        | Auth bearer token generated for guests |
| match_token | string | `xxx.xxx.xxx` |          | Token used for game connection         |


## Example
```bash
curl -XGET http://localhost:8080/api/v1/match/queue-in \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer xxx.xxx.xxx'
```
```json
status: 200
{
    "guest_token": null, 
    "match_token": "xxx.xxx.xxx",
}
```
