const html = /*html*/`
<form id="join-box">
    <div style="display: flex">
        <input
        class="text-input"
        name="playerName"
        type="text"
        placeholder="Your Name"
        autocomplete="off"
        autocapitalize="off"
        style="flex-grow: 1"
        />
    </div>
    
    <div id="join-box-buttons">
        <div id="roomCodeJoin" class="join-form-submit">
            <input
                class="text-input"
                name="roomCode"
                type="text"
                maxlength="4"
                placeholder="Room Code"
                autocomplete="off"
                style="text-transform: uppercase;"
            />
            <button type="submit" class="button">JOIN</button>
        </div>
        <div style="flex-grow: 1">
            <h3 style="text-align: center">- OR -</h3>
        </div>
        <div id="maxPlayersHost" class="join-form-submit">
            <input
                class="text-input"
                name="maxPlayers"
                type="number"
                placeholder="Max Players"
                autocomplete="off"
            />
            <button id="host-btn" class="button">HOST</button>
        </div>
    </div>
</form>
`
const css = /*css*/`
.join-form-submit{
    display: flex;
    flex-direction: column;
}

#join-box input{
    margin-bottom: 0.5rem;
}

#join-box-buttons{
    display: flex;
    justify-content: space-between;
    flex-wrap: wrap;
}
`
import { ComponentBase } from "./component-base.js";
import { CardView } from "./card-view.js";
import { buttonCSS, inputCSS } from "./css-snippets.js";
export class JoinForm extends ComponentBase{
    constructor(){
        super(html, css);

        this.shadowRoot.adoptedStyleSheets.push(buttonCSS(), inputCSS())
    }

    connectedCallback(){
        this.joinForm = this.shadowRoot.getElementById("join-box");

        this.joinForm.addEventListener("submit", this.handleJoinFormSubmit);
    }

    disconnectedCallback() {
        this.joinForm.removeEventListener("submit", this.handleJoinFormSubmit);
    }

    handleJoinFormSubmit = (event) => {
        event.preventDefault();

        const formData = new FormData(event.target);

        const playerName = (formData.get("playerName") || "").toString();

        const roomCode = (formData.get("roomCode") || "").toString().toUpperCase();
        const maxPlayers = (formData.get("maxPlayers") || "").toString();

        const connectionInfo = {
        type: roomCode?.length > 0 ? "join" : "host",
        playerName,
        maxPlayers,
        roomCode,
        };

        localStorage.setItem(
        "tworooms-connectionInfo",
        JSON.stringify(connectionInfo)
        );

        window.open("/play", "_self");
    }
}

customElements.define('join-form', JoinForm)