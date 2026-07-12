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
            <button type="submit" class="join-form-button secondary">JOIN</button>
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
            <button id="host-btn" class="join-form-button secondary">HOST</button>
        </div>
    </div>
</form>
`
const css = /*css*/`
.join-form-button {
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

.join-form-button:hover {
    background: linear-gradient(
        135deg,
        #343443,
        #252532
    );
    border-color: var(--color-gold-bright);
    color: var(--color-gold-bright);
    box-shadow: 0 0 10px rgba(198, 168, 91, 0.3);
}

.join-form-button:active {
    transform: translateY(1px);
    box-shadow: 0 0 6px rgba(198, 168, 91, 0.2);
}

.join-form-button.secondary {
    background: var(--color-surface);
    border-color: var(--color-muted);
    color: var(--color-cream);
}

.join-form-button.secondary:hover {
    border-color: var(--color-gold-bright);
    color: var(--color-gold-bright);
    box-shadow: 0 0 8px rgba(198, 168, 91, 0.2);
}

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

input {
  padding: 0.8rem 1rem;
  font-size: 0.95rem;
  font-family: var(--font-body);
  color: var(--color-cream);

  background: var(--color-surface);
  border: var(--border-thin);
  border-radius: var(--radius);
  box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.4);

  letter-spacing: 0.08em;
  transition: all 150ms ease;
}

input::placeholder {
  color: var(--color-muted);
  text-transform: uppercase;
  letter-spacing: 0.1em;
}

input:hover {
  border-color: var(--color-gold-bright);
  box-shadow: inset 0 2px 6px rgba(198, 168, 91, 0.15);
}

input:focus {
  outline: none;
  border-color: var(--color-gold-bright);
  box-shadow: 0 0 8px rgba(198, 168, 91, 0.3);
}
`
import { ComponentBase } from "./component-base.js";
export class JoinForm extends ComponentBase{
    constructor(){
        super(html, css);
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