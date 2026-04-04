# WebSocket Events

## Event Types
| Name | Paylaod | Usage | 
| ---|---|---|
| connect | Connect| On client connect | 
| words | Words| Words for this game | 
| countdown | Countdown| Countdown before game starts | 
| start | Start| Game start | 
| ping | Ping| Client ping to let the server now it is still connected | 
| keystroke | Keystroke| On every client keystroke | 
| finish | Finish| Server notifiy clients that the game has finished | 
| correction | Correction| Correct client state | 
| tooLate | TooLate| When client connects to a started game | 


## Event Structure
#### Header
| Name | Type | Required | Description |
|---|---|---|---|
| userID | string | v |User ID |

#### Body
| Name | Type | Required | Description |
|---|---|---|---|
| type | string | v | Name of event |
| payload | any |  | Event detail |


## Paylaods
### Connect
No payload

Example:
```json
{
  "type": "connect"
}
```

### Words
| Name | Type | Description |
|---|---|---|
| text | string | Text for this game |

Example:
```json
{
  "type": "words",
  "payload": {
    "text": "aaa bbb ccc ddd"
  }
}
```

### Countdown
| Name | Type | Description |
|---|---|---|
| number | number | Current countdown number|

Example:
```json
{
  "type": "countdown",
  "payload": {
    "number": 5
  }
}
```

### Start
No payload

Example:
```json
{
  "type": "start"
}
```

### Ping
No payload

Example:
```json
{
  "type": "ping"
}
```

### Keystroke
|Name | Type | Description|
|---|---|---|
|wordIndex | number | Current word index |
|charIndex | number | Current character index|
|key | string| What key is typed|

Example:
```json
{
  "type": "keystroke",
  "payload": {
    "wordIndex": 5,
    "charIndex": 6,
    "key": "a"
  }
}
```

### Finish
No payload

Example:
```json
{
  "type": "finish"
}
```

### Correction
|Name | Type | Description|
|---|---|---|
|text | string | Correct recored input|

Example:
```json
{
  "type": "correction",
  "payload": {
    "text": "aaa bbbb ccccc"
  }
}
```

### Too Late
|Name | Type | Description|
|---|---|---|
|text | string | Correct recored input|

Example:
```json
{
  "type": "tooLate"
}
```
