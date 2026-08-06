# Battleship

## Overview

Online multiplayer Battleship for the terminal.

- One server process that is authoritative for all game state. It accepts
  WebSocket connections, hosts any number of concurrent games, and pushes
  state updates to clients whenever anything changes (event-driven: clients
  spend most of their time waiting, and are notified of changes rather than
  polling).
- A client that talks to the server: spins up, connects over WebSocket,
  renders boards as ASCII in the terminal, reads input from the keyboard, and
  reacts to state pushes from the server.
- Game logic (rules, board, ships) lives in the `gamecode` package as
  pure functions, so it is fully unit-testable and
  reusable.
- A shared protocol module defines every message that crosses the wire, so
  client and server can never silently disagree on the format.

Larger vision (not MVP): scale this out into many small services that
can spin up on demand.

## Module Layout

Four separate Go modules:

```
Battleship/
├── protocol/   # shared types + message structs (the wire contract)
├── gamecode/   # pure game rules: placement, fire, turns, win detection (just useful functions)
├── server/     # WebSocket server, game registry, message router, state pushes
└── client/     # WebSocket client + ASCII TUI (rendering + terminal input)
```

`server` and `client` both import `protocol`; `server` imports `gamecode`.
Modules are versioned/released independently so the whole thing can later be
split into deployable services.

## Transport

WebSockets over a single endpoint (`/ws`). JSON messages in both directions.
WebSockets give us bidirectional, low-latency communication, which fits the
game's rhythm: a client submits a move, then waits while the server pushes
updates. This includes updates caused by the other player's moves.

## Service Design

- A game is a stateful object (the "class"): two players, two boards,
  whose turn it is, and a status (waiting / in-progress / finished).
- The server holds a registry: `map[gameID] *Game`. Game IDs (later,
  short human-friendly game codes) let clients create and join games, and let
  the server route moves to the right game.
- One process hosts many games and many clients at once:
  - a goroutine per client connection for reading incoming messages,
  - a single sender goroutine per game so writes never race,
  - per-game locking so concurrent moves from two players can't corrupt state,
  - cleanup when a game ends or a player disconnects (no leaked games).

### Game lifecycle

1. Client A sends `CreateGame` -> server makes a game, returns the game ID,
   status `waiting`.
2. Client B sends `JoinGame` with that ID -> players paired, boards randomly
   generated, status `in-progress`, both clients get a full `StateUpdate`.
3. Players alternate firing. Each `SubmitMove` is validated against the
   current state; every change produces a `StateUpdate` pushed to both.
4. When all of one player's ships are sunk, status becomes `finished` with a
   winner and the game is cleaned up. Clients can start a new game or exit.

## Storage plan

No persistent storage in the MVP. Everything lives in memory: the game
registry and every game's state exists only while the server runs. Games
finish and are garbage-collected. This is fine because a game only lives for
the duration of a match, nothing is needed in server post restart. If we
later want persistence (history, stats, resume-on-reconnect), the protocol
and gamecode packages are designed so a storage layer can be added behind the
server without touching client or rules.

## API plan

All traffic is JSON over the WebSocket. Every message has a `type` discriminator
so both sides can dispatch on message kind.

### From Client (messages)

- `CreateGame` — ask the server to start a new game; server replies with the
  game ID.
- `JoinGame` — join an existing game by ID (later: by short game code).
- `SubmitMove` — fire at a coordinate on the opponent's board.

### From Server (messages)

- `StateUpdate` — full authoritative game state after any change: your board,
  the opponent's board in fog-of-war (only hits and misses visible, never
  their ships), whose turn it is, game status, and the result of the last
  move (hit / miss / sank). Sent to both players whenever anything changes,
  so both stay in sync.
- `ErrorResponse` — a rejected message (bad move, game full, game not found,
  not your turn).

## Game Design

- Battleship rules: standard 5-ship fleet on a 10×10 grid.
- Random auto-placement for the MVP: the server generates a valid,
  seeded-random board per player so
  games are reproducible when debugging. Manual placement can come later.
- Firing: pure function — target coordinate → hit / miss / sank + updated
  board. Ships are sunk when all their cells are hit.
- Turns: strictly alternating; the server rejects a move made out of turn.
- Win: all ships sunk → game over, winner reported.
- Client UI: ASCII boards drawn in the terminal (own board + opponent
  board with fog-of-war), read from stdin via bufio. Commands: `fire A5`,
  `new`, `quit`. The TUI re-renders on every `StateUpdate` and prints
  hit/miss/sink/win messages.
