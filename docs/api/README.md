# API Documentation

Frontend Pages:
- Home Page
    - Queue In
- Match making page:
    - Quit
- In game:
    - Websocket events
    - API for player status
- Result
    - API for result
- Profile
    - statistics
        - best 
        - total avg
        - avg of 5
    - graphs
        - query: date range
        - x: time
        - y: wpm / acc
    - history 
        - time / wpm / acc / rank 
        - sort xxx

| Method    | Path                       | Description           | Link     |
|-----------|----------------------------|-----------------------|----------|
| POST      | /api/v1/match/queue-in     | Queue in              | [Link]() |
| POST      | /api/v1/match/quit         | Quit                  | [Link]() |
| Websocket | /api/v1/game/ws            | In game events        | [Link]() |
| GET       | /api/v1/game/player/status | In game player status | [Link]() |
| GET       | /api/v1/game/result        | Result                | [Link]() |
| GET       | /api/v1/profile/statistics | Profile statistics    | [Link]() |
| GET       | /api/v1/profile/graphs     | Profile graphs        | [Link]() |
| GET       | /api/v1/profile/history    | Profile history       | [Link]() |
