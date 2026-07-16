const html = /*html*/`
<div id="card-container">
    <button id="close">X</button>
    <img id="role-card" src="" alt="Your Role Card" />
</div>
`

const css = /*css*/`
#card-container {
    position: fixed;
    inset: 0;
    display: grid;
    place-items: center;
    z-index: 1000;
    background: rgba(0, 0, 0, 0.75);
}

#close{
    position: absolute;
    top: 0;
    right: 0;
    font-size: 72px;
    color: white;
    background-color: transparent;
    border: none;
    cursor: pointer;
}

#role-card {
    width: 100vw;
    height: 100dvh;
    object-fit: contain;
}

@media (max-width:1000px){
    #close {
        font-size: 156px;
    }
}
`

import { ComponentBase } from "./component-base.js";
export class CardView extends ComponentBase{
    static observedAttributes = ['team', 'role', 'show']

    team = '';
    role = '';

    constructor(){
        super(html, css);

    }

    connectedCallback(){
        this.team = this.getAttribute("team");
        this.role = this.getAttribute("role");
        this.show = this.getAttribute("show") === true.toString();

        this.shadowRoot.getElementById("close").addEventListener("click", this.close)
    }

    diconnectedCallback(){

    }

    attributeChangedCallback(name, oldVal, newVal){
        this.team = this.getAttribute("team");
        this.role = this.getAttribute("role");
        this.show = this.getAttribute("show") === true.toString();

        if(this.show){
            let newSrc = `/assets/images/cards/${this.team}/${this.role}.PNG`
            this.shadowRoot.getElementById("role-card").setAttribute("src", newSrc);
            this.shadowRoot.getElementById("card-container").style.visibility = '';
        }else{
            this.shadowRoot.getElementById("card-container").style.visibility = 'hidden';
        }

    }

    close = () =>{
        this.setAttribute("show", false);
    }
}

customElements.define('card-view', CardView)