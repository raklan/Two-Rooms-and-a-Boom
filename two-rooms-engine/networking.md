# Networking Flow
Client-Server relationships diagrammed below

## Gameplay Loop 
### Host
```mermaid
sequenceDiagram
    actor Client
    participant Server@{"type": "database"}

    Client ->> Server: /Host
    Server ->> Client: LobbyInfo
    Client ->> Server: StartGame
    Server ->> Client: GameInfo

    loop Repeated 3 times
        Client ->> Server: StartRound
        Server ->> Client: RoundStart
        Note over Client, Server: Round Flow (See Below)
        Server ->> Client: RoundEnd
        Note over Client, Server: Leaders exchange hostages
    end
    Server ->> Client: GameOver
    Server ->> Client: Close
```

### Non-Host
```mermaid
sequenceDiagram
    actor Client
    participant Server@{"type": "database"}
    Client ->> Server: /Join
    Server ->> Client: LobbyInfo
    Server ->> Client: GameInfo

    loop Repeated 3 times
        Server ->> Client: RoundStart
        Note over Client, Server: Round Flow (See Below)
        Server ->> Client: RoundEnd
        Note over Client, Server: Leaders exchange hostages
    end
    Server ->> Client: GameOver
    Server ->> Client: Close
```

## Round Flow
### Leaders Only
#### End-of-Round Hostage Exchange
```mermaid
sequenceDiagram
    actor Leader as Room Leader
    participant Server@{"type":"database"}
    actor Other Leader

    Leader <<->> Other Leader: RoundEnd
    Leader ->> Server: HostageExchange
    Other Leader ->> Server: HostageExchange
    create participant P@{"type":"collections"} as All Players
    Server ->> P: HostageExchangeComplete
    Note over Leader, P: Wait for Host to start next round
```

#### Abdication
```mermaid
sequenceDiagram
    actor Leader as Room Leader
    participant Server@{"type":"database"}
    actor Nominee

    Note over Leader, Nominee: Part of Round Flow
    Leader ->> Server: Abdicate
    Server ->> Nominee: PendingAbdication
    Nominee ->> Server: RespondAbdication
    alt Nominee Rejects
        Server ->> Leader: AbdicationRejected
    else Nominee Accepts
        create participant R@{"type":"collections"} as All Players in Room
        Server ->> R: NewLeader
    end
```
### Everyone
#### Usurption
```mermaid
sequenceDiagram
    actor Usurper
    participant Server@{"type":"database"}

    Note over Usurper, Server: Part of Round Flow
    destroy Usurper
    Usurper ->> Server: Usurp
    create participant R@{"type":"collections"} as All Players in Room
    Server ->> R: PendingUsurption
    R ->> Server: UsurpVote
    Note right of R: Usurper also sends separate vote. Starting Usurption != Voting
    alt Usurption Fails
        Server ->> R: UsurptionFailed
    else Usurption Succeeds
        Server ->> R: NewLeader
    end
```

#### Card Sharing
```mermaid
sequenceDiagram
    actor Player
    participant Server@{"type":"database"}
    actor O as Other Player

    Note over Player, O: Part of Round Flow
    Player ->> Server: CardShare
    Server ->> O: PendingCardShare
    O ->> Server: RespondCardShare
    alt Other Player rejects
        Server ->> Player: CardShareRejected
    else Other Player accepts
        Player <<->> O: SharedCard
    end
```