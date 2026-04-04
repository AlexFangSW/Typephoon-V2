```
GET /api/v1/profile/me
```
## Request
### Cookie
| Name  | Type   | Default | Example       | Required | Description  |
| ---   | ---    | ---     | ---           | ---      | ---          |
| TP_AT | string | N/A     | `xxx.xxx.xxx` | v        | Access token |


## Response
### 200
#### Body
| Name | Type   | Example | Nullable | Description |
| ---  | ---    | ---     | ---      | ---         |
| id   | string | abc123  |          | User ID     |
| name | string | Alex    |          | Username    |

### 401
#### Body
| Name  | Type   | Example        | Nullable | Description       |
| ---   | ---    | ---            | ---      | ---               |
| error | string | `Unauthorized` |          | Error description |

## Example
```bash
curl -XGET http://localhost:8080/api/v1/profile/me \
    -H 'Content-Type: application/json' \
    --cookie 'TP_AT=xxx.xxx.xxx'
```
```json
status: 200
{
    "name": "Alex",
    "id": "abc123"
}
```
