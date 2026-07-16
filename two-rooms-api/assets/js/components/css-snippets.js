function buttonCSS(){
    let css = /*css*/`
        .button {
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

        .button:hover {
            background: linear-gradient(
                135deg,
                #343443,
                #252532
            );
            border-color: var(--color-gold-bright);
            color: var(--color-gold-bright);
            box-shadow: 0 0 10px rgba(198, 168, 91, 0.3);
        }
    `
    let stylesheet = new CSSStyleSheet();
    stylesheet.replaceSync(css);
    return stylesheet;
}

function inputCSS(){
    let css = /*css*/`
        .text-input {
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
        
        .text-input::placeholder {
          color: var(--color-muted);
          text-transform: uppercase;
          letter-spacing: 0.1em;
        }
        
        .text-input:hover {
          border-color: var(--color-gold-bright);
          box-shadow: inset 0 2px 6px rgba(198, 168, 91, 0.15);
        }
        
        .text-input:focus {
          outline: none;
          border-color: var(--color-gold-bright);
          box-shadow: 0 0 8px rgba(198, 168, 91, 0.3);
        }
    `

    let stylesheet = new CSSStyleSheet();
    stylesheet.replaceSync(css);
    return stylesheet;
}

export { buttonCSS, inputCSS }