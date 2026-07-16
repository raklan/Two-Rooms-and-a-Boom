const html = /*html*/`
<div id="view-container">
    <div id="timer-container">
        <p id="text"></p>
    </div>
    <div id="controls">
        <button id="start-round">Start Next Round</button>
        <button id="show-card">See my Card</button>
    </div> 
    <card-view id="card-view" team="blue" role="player" show="false"></card-view>   
</div>
`;

const css = /*css*/``;

import { ComponentBase } from "./component-base.js";
import { CardView } from "./card-view.js";
export class GameplayView extends ComponentBase{
    constructor(){
        super(html, css);

        this.playerId = null;
        this.isHost = false;
        this.team = null;
        this.role = null;
        this.room = null;
    }

    connectedCallback(){
        ws.addEventListener('message', this.handleWsMessage);
        ws.addEventListener('error', this.handleWsError);

        this.shadowRoot.getElementById("start-round").onclick = this.startNextRound;
        this.shadowRoot.getElementById("show-card").onclick = this.showCard;
    }

    disconnectedCallback(){
        ws.removeEventListener('message', this.handleWsMessage);
        ws.removeEventListener('error', this.handleWsError);
    }
    handleWsMessage = (wsMsg) => {
        var message = JSON.parse(wsMsg.data);

        switch(message.type){
            case 'LobbyInfo':
                this.handleLobbyInfo(message.data);
                break;
            case 'GameState':
                this.handleGameState(message.data);
                break;
            case 'RoundStart':
                this.handleRoundStart(message.data);
                break;
            default:
                console.error("Could not handle webSocket message of type", message.type)
                break;
        }
    }

    handleWsError = (wsMsg) => {
        console.error('An error occurred with the WebSocket', wsMsg)
        showNotification('Something went wrong trying to establish a connection to the Lobby. Ensure you have the right Room Code and the Name you chose isn\'t taken and try again. If the issue persists, feel free to reach out to me on the Github repository. Sorry about that!', 'Error')
    }

    // #region Message Handlers
    handleLobbyInfo(data){
        if(!(this.playerId?.length > 0)){
            this.playerId = data.playerID;
            this.isHost = data.playerID === data.lobbyInfo.host.Id;
        }
    }

    handleGameState(gameState){
        let thisPlayer = gameState.players.find(p => p.id == this.playerId);

        this.room = thisPlayer.room;
        this.role = thisPlayer.role;
        this.team = thisPlayer.team;

        let cardView = this.shadowRoot.getElementById("card-view");
        cardView.setAttribute('team', this.team);
        cardView.setAttribute('role', this.role);
    }

    handleRoundStart(roundData){
        console.log('round starting', roundData)
    }

    startNextRound(){
        sendWsMessage(ws, 'StartRound', {})
    }

    showCard = () => {
        this.shadowRoot.getElementById("card-view").setAttribute("show", true);
    }
}

customElements.define('gameplay-view', GameplayView);