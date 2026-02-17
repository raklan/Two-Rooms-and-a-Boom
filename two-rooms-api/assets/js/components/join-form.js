const template = document.createElement('template');
const innerHTML = `
<style>
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
</style>
<form id="join-box">
    <input
      class="text-input"
      name="playerName"
      type="text"
      placeholder="ENTER YOUR NAME"
      autocomplete="off"
      autocapitalize="off"
    />
    <input
      class="text-input"
      name="roomCode"
      type="text"
      maxlength="4"
      placeholder="ENTER ROOM CODE"
      autocomplete="off"
      style="text-transform: uppercase;"
    />
    <button type="submit" class="join-form-button">JOIN</button>
</form>
<div style="align-self: center;">
    <h3 style="text-align: center;">-OR-</h3>
    <button id="host-btn" class="join-form-button secondary">HOST</button>
</div>
`
import { ComponentBase } from "./component-base.js";
export class JoinForm extends HTMLElement{
    constructor(){
        super();
        this.innerHTML = innerHTML;
    }

    connectedCallback(){
        this.joinForm = this.shadowRoot.getElementById("join-box");
        this.hostBtn = this.shadowRoot.getElementById("host-btn");

        this.joinForm.addEventListener("submit", this.handleSubmit);
        this.hostBtn.addEventListener("click", this.handleHost); 
    }

    disconnectedCallback() {
        this.joinForm.removeEventListener("submit", this.handleSubmit);
        this.hostBtn.removeEventListener("click", this.handleHost);
    }

    handleJoinFormSubmit = (event) => {
        event.preventDefault();

        const formData = new FormData(this.form);

        const roomCode = (formData.get("roomCode") || "").toString().toUpperCase();
        const playerName = (formData.get("playerName") || "").toString();

        const connectionInfo = {
        type: "join",
        roomCode,
        playerName,
        };

        localStorage.setItem(
        "tworooms-connectionInfo",
        JSON.stringify(connectionInfo)
        );

        window.open("/play", "_self");
    }

    handleHost = () => {
        console.log('hosting')
    }
}

customElements.define('join-form', JoinForm)