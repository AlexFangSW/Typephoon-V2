```
GET /api/v1/profile/history
```
## Request
### Header
| Name          | Type   | Default | Example              | Required | Description       |
| ---           | ---    | ---     | ---                  | ---      | ---               |
| Authorization | string | N/A     | `Bearer xxx.xxx.xxx` | v        | Auth bearer token |

### Query
| Name | Type   | Default | Example | Required | Description                         |
| ---  | ---    | ---     | ---     | ---      | ---                                 |
| size | number | `20`    | `20`    |          | Number of results to return, max 50 |
| page | number | `1`     | `1`     |          | Page number                         |

## Response
### 200
#### Body
| Name      | Type              | Example | Nullable | Description              |
| ---       | ---               | ---     | ---      | ---                      |
| items     | array[ResultItem] | N/A     |          | A list of result items   |
| count     | number            | `10`    |          | Number of retuned items  |
| total     | number            | `9999`  |          | Total number of results  |
| next_page | boolean           | `true`  |          | Is there a next page     |
| prev_page | boolean           | `false` |          | Is there a previous page |

#### ResultItem
| Name | Type   | Example                | Nullable | Description                            |
| ---  | ---    | ---                    | ---      | ---                                    |
| ts   | string | `2026-03-25T08:22:08Z` |          | Result create timestamp RFC3339 format |
| wpm  | number | `123`                  |          | Word per minute                        |
| acc  | number | `98`                   |          | Accuracy                               |
| rank | number | `1`                    |          | Rank for the game                      |

### 401
#### Body
| Name  | Type   | Example        | Nullable | Description       |
| ---   | ---    | ---            | ---      | ---               |
| error | string | `Unauthorized` |          | Error description |

## Example
```bash
curl -XGET http://localhost:8080/api/v1/profile/history \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer xxx.xxx.xxx'
```
```json
status: 200
{
    "total": 999,
    "count": 20,
    "prev_page": false,
    "next_page": true,
    "items": [
        {
            "ts": "2026-03-25T08:22:08Z",
            "wpm": 123,
            "acc": 98,
            "rank": 1,
        },
        {
            "ts": "2026-03-25T08:22:08Z",
            "wpm": 123,
            "acc": 98,
            "rank": 1,
        },
        ...
        {
            "ts": "2026-03-25T08:22:08Z",
            "wpm": 123,
            "acc": 98,
            "rank": 1,
        } 
    ]
}
```
