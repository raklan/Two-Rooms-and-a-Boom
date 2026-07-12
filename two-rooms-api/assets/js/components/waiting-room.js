const html = /*html*/`
<div id="waiting-room-container">
    <div>
        <h2>Players</h2>
        <div id="player-list">
            
        </div>
    </div>
</div>
`

const css = /*css*/`
#player-list{
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    justify-items: center;
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
    }

    disconnectedCallback(){
        ws.removeEventListener('message', this.handleWsMessage);
        ws.removeEventListener('error', this.handleWsError);
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
            this.setupNewPlayer(data.playerId, data.lobbyInfo.host.id)
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

    setupNewPlayer(playerId, hostId){
        this.playerId = playerId;
        this.isHost = playerId === hostId;
    }
}

customElements.define('waiting-room', WaitingRoom);