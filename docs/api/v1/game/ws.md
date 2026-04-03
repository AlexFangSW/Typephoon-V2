# WebSocket Events

| Name | Body | Usage | 
| ---|---|---|
| connect | ConnectBody | On client connect | 
| words | WordsBody | Words for this game | 
| countdown | CountdownBody | Countdown before game starts | 
| start | StartBody | Game start | 
| ping | PingBody | Client ping to let the server now it is still connected | 
| keystroke | KeystrokeBody | On every client keystroke | 
| finish | FinishBody | Server notifiy clients that the game has finished | 
| correction | CorrectionBody | Correct client state | 
| tooLate | TooLateBody | When client connects to a started game | 
