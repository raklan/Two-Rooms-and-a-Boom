export class ComponentBase extends HTMLElement{
    constructor(template){
        super();
        this.attachShadow({mode: 'open'});
        this.shadowRoot.appendChild(template.content.cloneNode(true));
    }
}