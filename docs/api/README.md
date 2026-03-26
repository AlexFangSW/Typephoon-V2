# API Documentation

Frontend Pages:
- In game:
    - Websocket events
    - API for player status
- Result
    - API for result

- AUTH:
    - login
    - logout
    - token-refresh

| Method    | Path                       | Description           | Link                               |
|-----------|----------------------------|-----------------------|------------------------------------|
| POST      | /api/v1/match/queue-in     | Queue in              | [Link](./v1/match/queue-in.md)     |
| Websocket | /api/v1/game/ws            | In game events        | [Link](./v1/game/ws.md)            |
| GET       | /api/v1/game/player/status | In game player status | [Link](./v1/game/player/status.md) |
| GET       | /api/v1/game/result        | Result                | [Link](./v1/game/result.md)        |
| GET       | /api/v1/profile/statistics | Profile statistics    | [Link](./v1/profile/statistics.md) |
| GET       | /api/v1/profile/graphs     | Profile graphs        | [Link](./v1/profile/graphs.md)     |
| GET       | /api/v1/profile/history    | Profile history       | [Link](./v1/profile/history.md)    |
| GET       | /api/v1/auth/login         | Login                 | [Link](./v1/auth/login.md)         |
| POST      | /api/v1/auth/logout        | Logout                | [Link](./v1/auth/logout.md)        |
| POST      | /api/v1/auth/refresh       | Refresh access token  | [Link](./v1/auth/refresh.md)       |
