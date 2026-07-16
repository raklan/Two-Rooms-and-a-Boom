const html = /*html*/`
<div class="timer-container">
  <svg class="timer-svg" viewBox="0 0 100 100">
    <!-- Background circle -->
    <circle class="timer-path-elapsed" cx="50" cy="50" r="45"></circle>
    <!-- Animated countdown path -->
    <circle id="timer-path-alert" class="timer-path-alert" cx="50" cy="50" r="45"></circle>
  </svg>
  <span id="timer-label" class="timer-label">0:00</span>
</div>
`

const css = /*css*/`
.timer-container {
  position: relative;
  width: 200px;
  height: 200px;
}

.timer-svg {
  transform: scaleX(-1); /* Flips SVG to count clockwise */
}

.timer-path-elapsed {
  stroke-width: 7px;
  stroke: #e6e6e6;
  fill: none;
}

.timer-path-alert {
  stroke-width: 7px;
  stroke-linecap: round;
  fill: none;
  transition: stroke-dasharray 1s linear, stroke 1s linear;
}

.timer-label {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 2rem;
  font-family: sans-serif;
}
`

import { ComponentBase } from "./component-base.js";

export class TimerBox extends ComponentBase{
    static observedAttributes = ['round', 'round-length'];

    constructor(){
        super(html, css);

        this.round = 0;
        this.roundLength = 0;
        this.time = 0;
        this.timerInterval = null;
    }

    connectedCallback(){

    }

    disconnectedCallback(){

    }

    attributeChangedCallback(name, oldVal, newVal){
        this.round = this.getAttribute('round');
        this.roundLength = this.getAttribute('round-length');

        this.startTimer();
    }

    formatTime(time) {
        const minutes = Math.floor(time / 60);
        let seconds = time % 60;
        if (seconds < 10) seconds = `0${seconds}`;
        return `${minutes}:${seconds}`;
    }

    startTimer = () => {
        this.time = 0;
        const wc = this;

        const label = this.shadowRoot.getElementById("timer-label");
        const path = this.shadowRoot.getElementById("timer-path-alert");

        const startTime = new Date();

        function timerTick(){
            let elapsedSeconds = Math.floor((new Date() - startTime) / 1000);
            wc.time = elapsedSeconds;
            let timeLeft = wc.roundLength - wc.time;
            
            // Update text
            label.innerText = wc.formatTime(timeLeft);
            
            // Update SVG stroke
            let rawTimeFraction = timeLeft / wc.roundLength;
            if(timeLeft !== 0){
                rawTimeFraction = (timeLeft - 1) / wc.roundLength;
            }
            const circleDasharray = `${(rawTimeFraction * 283).toFixed(0)} 283`;
            path.setAttribute("stroke-dasharray", circleDasharray);

            // Color shift to red in last 10 seconds
            if (timeLeft > 10) {
                path.setAttribute("stroke", "gold");
            }else{
                path.setAttribute("stroke", "red")
            }

            if(timeLeft > 0){
                setTimeout(timerTick, 1000)
            }
        }

        this.timerInterval = setTimeout(timerTick, 1000);
    }

    
}

customElements.define('timer-box', TimerBox)