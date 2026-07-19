const html = /*html*/`
<div id="select-backdrop">
    <h3 id="select-title">Select a player to share cards with</h3>
    <form id="player-select-form">    
    </form>
</div>
`

const css = /*css*/`
#select-backdrop{
    position: fixed;
    inset: 0;
    display: grid;
    place-items: center;
    z-index: 1000;
    background: rgba(0, 0, 0, 0.75);
}
`


import { ComponentBase } from "./component-base.js";
import { PlayerListUpdateEvent, PlayerSelectedEvent } from "./events.js";
export class PlayerSelector extends ComponentBase{
    static observedAttributes = ['title', 'show'];
    playerList = [];
    show = false;

    constructor(){
        super(html, css)
    }

    connectedCallback(){
        this.shadowRoot.getElementById('player-select-form').addEventListener('submit', this.handlePlayerSelectFormSubmit);

        this.addEventListener(PlayerListUpdateEvent.EVENT_NAME, this.handlePlayerListUpdate);        

        this.dispatchEvent(new PlayerListUpdateEvent([{name: 'ryan', id: 1}, {name: 'christina', id: 2}]))
    }

    disconnectedCallback(){
        this.removeEventListener(PlayerListUpdateEvent.EVENT_NAME, this.handlePlayerListUpdate);
    }

    attributeChangedCallback(name, oldVal, newVal){        
        let title = this.getAttribute('title');
        this.show = this.getAttribute('show') === true.toString();

        this.shadowRoot.getElementById('select-title').innerText = title;

        if(this.show){
            this.shadowRoot.getElementById("select-backdrop").style.visibility = 'visible';
        }else{
            this.shadowRoot.getElementById("select-backdrop").style.visibility = 'hidden';
        }
    }

    handlePlayerListUpdate = (event) => {
        this.playerList = event.newPlayerList;

        let form = this.shadowRoot.getElementById("player-select-form");
        form.replaceChildren();
        for(let player of this.playerList){
            let input = document.createElement('input');
            input.type = 'radio';
            input.id = `player-${player.id}`;
            input.name = 'player';
            input.value = player.id;
            
            let label = document.createElement('label');
            label.setAttribute('for', input.id);
            label.innerText = player.name;

            form.appendChild(input);
            form.appendChild(label);
        }
        let button = document.createElement('button');
        button.innerText = 'Submit';
        form.appendChild(button);
    }

    handlePlayerSelectFormSubmit = (event) => {
        event.preventDefault();
        
        let form = this.shadowRoot.getElementById("player-select-form");

        let selectedId = form['player'].value;
        
        let player = this.playerList.find(p => p.id == selectedId);

        this.parentElement.dispatchEvent(new PlayerSelectedEvent(player));

        this.setAttribute('show', false);
    }
}

customElements.define('player-selector', PlayerSelector)