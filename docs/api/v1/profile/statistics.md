```
GET /api/v1/profile/statistics
```
## Request
### Cookie
| Name  | Type   | Default | Example       | Required | Description  |
| ---   | ---    | ---     | ---           | ---      | ---          |
| TP_AT | string | N/A     | `xxx.xxx.xxx` | v        | Access token |

## Response
### 200
#### Body
| Name    | Type       | Example | Nullable | Description        |
| ---     | ---        | ---     | ---      | ---                |
| best    | ResultItem | N/A     |          | All time best      |
| avg_10  | ResultItem | N/A     |          | Average of last 10 |
| avg_all | ResultItem | N/A     |          | Average overall    |

#### ResultItem
| Name | Type   | Example                | Nullable | Description                            |
| ---  | ---    | ---                    | ---      | ---                                    |
| ts   | string | `2026-03-25T08:22:08Z` |          | Result create timestamp RFC3339 format |
| wpm  | number | `123`                  |          | Word per minute                        |
| acc  | number | `98`                   |          | Accuracy                               |

### 401
#### Body
| Name  | Type   | Example        | Nullable | Description       |
| ---   | ---    | ---            | ---      | ---               |
| error | string | `Unauthorized` |          | Error description |

## Example
```bash
curl -XGET http://localhost:8080/api/v1/profile/statistics \
    -H 'Content-Type: application/json' \
    --cookie 'TP_AT=xxx.xxx.xxx'
```
```json
status: 200
{
    "best": {
        "ts": "2026-03-25T08:22:08Z",
        "wpm": 123,
        "acc": 98
    },
    "avg_10": {
        "ts": "2026-03-25T08:22:08Z",
        "wpm": 123,
        "acc": 98
    },
    "avg_all": {
        "ts": "2026-03-25T08:22:08Z",
        "wpm": 123,
        "acc": 98
    }
}
```
