const html = /*html*/`
<div id="view-container">
    <div id="timer-container">
        <timer-box id="round-timer"></timer-box>
        <button id="start-round">Start Next Round</button>
    </div>
    <div id="controls">
        <button id="show-card">See my Card</button>
    </div> 
    <div id="gameover-alert">
        <div style="border: 1px solid var(--color-gold); text-align: center; background-color: black; padding: 15px">
            <div id="victory">Your team won!</div>
            <a href="/" style="color: var(--color-gold)">Next</a>
        </div>
    </div>
    <card-view id="card-view" team="blue" role="player" show="false"></card-view>   
</div>
`;

const css = /*css*/`
#gameover-alert{
    color: var(--color-gold);
    background-color: rgba(0,0,0,0.75);
    position: fixed;
    z-index: 1000;
    inset: 0;
    display: grid;
    place-items: center;
    visibility: hidden;
}
`;

import { ComponentBase } from "./component-base.js";
import { CardView } from "./card-view.js";
import { TimerBox } from "./timer-box.js";
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
            case 'RoundEnd':
                this.handleRoundEnd(message.data);                
                break;
            case 'GameOver':
                this.handleGameOver(message.data);
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
            this.isHost = data.playerID === data.lobbyInfo.host.id;
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

        if(this.isHost){
            this.shadowRoot.getElementById("start-round").style.visibility = '';
        }else{
            this.shadowRoot.getElementById("start-round").style.visibility = 'hidden';
        }
    }

    handleRoundStart(roundData){
        console.log('round starting', roundData)

        this.room = roundData.room;

        const timer = this.shadowRoot.getElementById("round-timer");
        timer.setAttribute("round", roundData.roundNumber);
        timer.setAttribute("round-length", roundData.roundLength);

        this.shadowRoot.getElementById("start-round").style.visibility = 'hidden';
    }

    handleRoundEnd(roundData){
        console.log('round ending', roundData);

        if(this.isHost){
            this.shadowRoot.getElementById("start-round").style.visibility = 'visible';
        }else{
            this.shadowRoot.getElementById("start-round").style.visibility = 'hidden';
        }
    }

    handleGameOver(gameOverData){
        let gameOverMessage = `Your team ${gameOverData.victory ? 'won' : 'lost'}!`

        this.shadowRoot.getElementById("victory").innerText = gameOverMessage;
        this.shadowRoot.getElementById('gameover-alert').style.visibility = 'visible';
    }

    startNextRound(){
        sendWsMessage(ws, 'StartRound', {})
    }

    showCard = () => {
        this.shadowRoot.getElementById("card-view").setAttribute("show", true);
    }
}

customElements.define('gameplay-view', GameplayView);