## Overview
- Service with game state
- Game code which talks to server


## Service Design
- Game as a class
- server holds list of classes mapped to game ids



#### Storage plan
Single source of truth game states stored on server, each player has a local copy of ONLY their maps for UI rendering.

#### API plan

##### From Client:
- Submitted move
- Create game / join game (use game code - later)
##### From Server:
- State update

## Game Design
- UI

Carrier (5 holes), Battleship (4 holes), Destroyer (2 holes), and two 3-hole ships (a Cruiser and a Submarine)

shipIds = {Carrier: 0, Battleship: 1, Destroyer: 2, Cruiser: 3, Submarine: 4}
#### Server game states
- Match ID [String]
- whos_turn [bool]
- Player 1 game state [2d list] (full grid, x for hits p1 has made on p2, o for misses p1 has made)
- Player 2 game state [2d list] (full grid, x for hits, o for misses)
- p1 ships [2d list] (full grid, shipId + x if hit e.g 2 or 2x)
- p2 ships [2d list] (full grid, shipId + x if hit e.g 2 or 2x)
- Player 1 ship counter [Map{ship name: num coords remaining}]
- Player 2 ship counter [Map{ship name: num coords remaining}]

#### client game states
- Match ID [String]
- whos_turn [bool]
- Player game state [2d list] (full grid, x for hits, o for misses)
- player ships [2d list] (full grid, shipId + x if hit e.g 2 or 2x)
- Player ship counter [Map{ship name: num coords remaining}]