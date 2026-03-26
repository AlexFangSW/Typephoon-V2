```
GET /api/v1/profile/graphs
```
## Request
### Header
| Name          | Type   | Default | Example              | Required | Description       |
| ---           | ---    | ---     | ---                  | ---      | ---               |
| Authorization | string | N/A     | `Bearer xxx.xxx.xxx` | v        | Auth bearer token |

### Query
| Name | Type   | Default | Example | Required | Description                           |
| ---  | ---    | ---     | ---     | ---      | ---                                   |
| size | number | `10`    | `10`    |          | Number of results to return, max 1000 |

## Response
### 200
#### Body
| Name  | Type              | Example | Nullable | Description              |
| ---   | ---               | ---     | ---      | ---                      |
| items | array[ResultItem] | N/A     |          | A list of result items   |
| count | number            | `10`    |          | Number of returned items |

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
curl -XGET http://localhost:8080/api/v1/profile/graphs \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer xxx.xxx.xxx'
```
```json
status: 200
{
    "count": 10,
    "items": [
        {
            "ts": "2026-03-25T08:22:08Z",
            "wpm": 123,
            "acc": 98,
        },
        {
            "ts": "2026-03-25T08:22:08Z",
            "wpm": 123,
            "acc": 98,
        },
        ...
        {
            "ts": "2026-03-25T08:22:08Z",
            "wpm": 123,
            "acc": 98,
        } 
    ]
}
```
