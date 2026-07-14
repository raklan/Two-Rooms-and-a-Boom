const html = /*html*/`
<div id="waiting-room-container">
    <div>
        <h2>Players</h2>
        <div id="player-list">
            
        </div>
    </div>
    <div>
        <button id="start-game-button" class="lobby-button">Start Game</button>
    </div>
</div>
`

const css = /*css*/`
#player-list{
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    justify-items: center;
}

#waiting-room-container{
    display: flex;
    flex-direction: column;
    justify-content: space-between;
}

.lobby-button {
    width: 100%;
    padding: 0.9rem 1rem;
    margin-top: 0.5rem;

    background: linear-gradient(
        145deg,
        #2a2a36,
        #1f1f28
    );

    border: var(--border-thin);
    border-radius: var(--radius);

    color: var(--color-gold);
    font-family: var(--font-heading);
    font-size: 0.95rem;
    letter-spacing: 0.15em;
    text-transform: uppercase;

    cursor: pointer;
    transition: all 150ms ease;
}

.lobby-button:hover {
    background: linear-gradient(
        135deg,
        #343443,
        #252532
    );
    border-color: var(--color-gold-bright);
    color: var(--color-gold-bright);
    box-shadow: 0 0 10px rgba(198, 168, 91, 0.3);
}
`

import { ComponentBase } from "./component-base.js";
export class WaitingRoom extends ComponentBase {
    constructor(){
        super(html, css);

        this.playerId = null;
        this.isHost = false;
    }

    connectedCallback(){
        ws.addEventListener('message', this.handleWsMessage)
        ws.addEventListener('error', this.handleWsError)

        this.shadowRoot.getElementById("start-game-button", this.startGame)
    }

    disconnectedCallback(){
        ws.removeEventListener('message', this.handleWsMessage);
        ws.removeEventListener('error', this.handleWsError);
    }

    startGame = () => {
        sendWsMessage(ws, 'startGame', {})
    }

    handleWsError = (wsMsg) => {
        console.error('An error occurred with the WebSocket', wsMsg)
        showNotification('Something went wrong trying to establish a connection to the Lobby. Ensure you have the right Room Code and the Name you chose isn\'t taken and try again. If the issue persists, feel free to reach out to me on the Github repository. Sorry about that!', 'Error')
    }

    handleWsMessage = (wsMsg) => {
        var message = JSON.parse(wsMsg.data);

        switch(message.type){
            case 'LobbyInfo':
                this.handleLobbyInfo(message.data);
                break;
            default:
                console.error("Could not handle webSocket message of type", message.type)
                break;
        }
    }

    handleLobbyInfo(data){
        if(!(this.playerId?.length > 0)){
            this.setupNewPlayer(data.playerID, data.lobbyInfo.host.id, data.lobbyInfo.roomCode)
        }

        this.renderPlayerList(data.lobbyInfo.players);
    }

    renderPlayerList(players){
        let playerList = this.shadowRoot.getElementById("player-list");

        playerList.replaceChildren();
        let playerEntry;
        for(let player of players){
            playerEntry = document.createElement("div");
            playerEntry.innerText = player.name;

            if(player.id === this.playerId){
                playerEntry.innerText += ' (You)'
            }

            playerList.appendChild(playerEntry);
        }
    }

    setupNewPlayer(playerId, hostId, roomCode){
        this.playerId = playerId;
        this.isHost = playerId === hostId;

        window.localStorage.setItem('tworooms-connectionInfo', JSON.stringify({
            type: 'rejoin',
            playerId: playerId,
            roomCode: roomCode
        }))
    }
}

customElements.define('waiting-room', WaitingRoom);