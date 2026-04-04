```
GET /api/v1/game/status
```
## Request
### Header
| Name | Type | Required | Description |
| ---| ---| ---|---|
| match_token | string (JWT) | v| Identify user and game id|

## Response
### 200
#### Body
| Name | Type | Example | Nullable | Description |
| ---| ---|--- | ---|---|
| N/A  | array[PlayerItem] | N/A | | Array of player status |

#### PlayerItem
| Name | Type | Example | Nullable | Description |
| ---| ---|--- | ---|---|
| connected  | bool | true |  | connection status |
| finish_ts   | string | `2026-03-25T08:22:08Z` | v | Finish timestamp RFC3339 format, could be null if the player didn't finish |
| wpm   | number | 99  | v | Word per minute |
| acc   | number | 90 | v | Accuracy |
| name   | string | Alex |  | Player name|
| id   | string | abc123 |  | Player ID|
| me   | bool | true |  | Is this the current player |
