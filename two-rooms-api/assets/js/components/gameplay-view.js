const html = /*html*/``;

const css = /*css*/``;

import { ComponentBase } from "./component-base.js";
export class GameplayView extends ComponentBase{
    constructor(){
        super(html, css);
    }

    connectedCallback(){
        ws.addEventListener('message', this.handleWsMessage);
        ws.addEventListener('error', this.handleWsError);
    }

    disconnectedCallback(){
        ws.removeEventListener('message', this.handleWsMessage);
        ws.removeEventListener('error', this.handleWsError);
    }

    handleWsMessage = (wsMsg) => {
        var message = JSON.parse(wsMsg.data);

        switch(message.type){
            case '':
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
}

customElements.define('gameplay-view', GameplayView);