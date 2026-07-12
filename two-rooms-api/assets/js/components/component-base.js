export class ComponentBase extends HTMLElement{
    constructor(html, css){
        super();
        this.attachShadow({mode: 'open'});

        const template = document.createElement("template");
        template.innerHTML = html;
        
        const styleSheet = new CSSStyleSheet();
        styleSheet.replaceSync(css);

        this.shadowRoot.appendChild(template.content.cloneNode(true));
        this.shadowRoot.adoptedStyleSheets = [styleSheet];
    }
}