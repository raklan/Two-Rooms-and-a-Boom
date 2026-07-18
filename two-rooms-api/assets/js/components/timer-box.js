const html = /*html*/`
<div class="timer-container">
  <svg class="timer-svg" viewBox="0 0 100 100">
    <!-- Background circle -->
    <circle class="timer-path-elapsed" cx="50" cy="50" r="45"></circle>
    <!-- Animated countdown path -->
    <circle id="timer-path-alert" class="timer-path-alert" cx="50" cy="50" r="45"></circle>
  </svg>
  <span id="timer-label" class="timer-label">
    <span id="round">Round 1</span>
    <span id="time-left">0:00</span>
  </span>
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
  transition: stroke-dasharray 1s linear, stroke var(--transition-speed) cubic-bezier(0.83, 0.08, 1, 1);
}

.timer-label #round{
  font-size: small;
  text-decoration: underline;
}

.timer-label {
  position: absolute;
  top: 45%;
  left: 50%;
  text-align: center;
  transform: translate(-50%, -50%);
  font-size: 2rem;
  font-family: sans-serif;
  color: var(--color-gold);
}
`

import { ComponentBase } from "./component-base.js";

export class TimerBox extends ComponentBase {
  static observedAttributes = ['round', 'round-length'];

  constructor() {
    super(html, css);

    this.round = 0;
    this.roundLength = 0;
    this.time = 0;
    this.timerInterval = null;
  }

  connectedCallback() {

  }

  disconnectedCallback() {

  }

  attributeChangedCallback(name, oldVal, newVal) {
    let newRound = this.getAttribute('round');
    this.roundLength = this.getAttribute('round-length');

    this.shadowRoot.getElementById("round").innerText = `Round ${newRound}`

    if (this.round !== newRound) {
      this.round = newRound;
      this.resetTimer();
      this.startTimer();
    }
  }

  formatTime(time) {
    if(time <= 0){
      return "0:00"
    }

    const minutes = Math.floor(time / 60);
    let seconds = time % 60;
    if (seconds < 10) seconds = `0${seconds}`;
    return `${minutes}:${seconds}`;
  }

  resetTimer = () => {
    const path = this.shadowRoot.getElementById("timer-path-alert");
    //The transition time will slowly fade the color to red over the course of the timer
    path.style.setProperty("--transition-speed", `0s`)
    path.setAttribute("stroke", "var(--color-gold)")
    path.setAttribute("stroke-dasharray", "283")
    path.style.setProperty("--transition-speed", `${Math.max(this.roundLength - 10, 10)}s`)
  }

  startTimer = () => {
    this.time = 0;
    const wc = this;

    const label = this.shadowRoot.getElementById("time-left");
    const path = this.shadowRoot.getElementById("timer-path-alert");

    const startTime = new Date();

    function timerTick() {
      let elapsedSeconds = Math.floor((new Date() - startTime) / 1000);
      wc.time = elapsedSeconds;
      let timeLeft = wc.roundLength - wc.time;

      // Update text
      label.innerText = wc.formatTime(timeLeft);

      // Update SVG stroke
      let rawTimeFraction = timeLeft / wc.roundLength;
      if (timeLeft !== 0) {
        rawTimeFraction = (timeLeft - 1) / wc.roundLength;
      }
      const circleDasharray = `${(rawTimeFraction * 283).toFixed(0)} 283`;
      path.setAttribute("stroke-dasharray", circleDasharray);
      path.setAttribute("stroke", "red")

      if (timeLeft > 0) {
        setTimeout(timerTick, 1000)
      }
    }

    this.timerInterval = setTimeout(timerTick, 0);
  }


}

customElements.define('timer-box', TimerBox)