# Typephoon V2
A redesign of the previous [Typephoon](https://github.com/AlexFangSW/Typephoon_api) project.  
Simplified architecture and implimentation while also preserving scallability and performance.

## Architecture
![Architecture](./docs/images/architecture.svg)

### General Server
Serves general APIs for things like "user profile", "authentication" ... etc.  

The web page will also be served here, we will still be using a frontend framework, but it will be built and 
served on one of the endpoints.

### Match Making Service
Currently the match making procedure is very simple:
We have a line of players waiting, the match making will check the wait list periodicly,
if there are players in line, we group them until no players are left. 
Games can start even if players are not full.

All data required for logic will be in memory.

It is designed with "Active-Passive" architecture. 
A matchmaker becomes active once it aquires the lock in the external cache.

### Game Server
This is implimented with the concept of [dedicated game servers](https://en.wikipedia.org/wiki/Game_server).  
We can have multiple game servers and each game server can host multiple games,
all player in the same game will be connected to the same game server via proxy.  

The game server will be the source of truth for all events, it will correct the client if any missmatch happens.  

> Originally thought of using [Agones](https://github.com/googleforgames/agones) to manage a pool of game servers,
but that seems a bit overkill for my project.

### Proxy
Direct the player connection to the correct game server according the the data stored in cache.   
It will also do some validation to make sure the connection is from a valid player.

## Screenshots
...

## Development
...
