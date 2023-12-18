import "./css/style.css";

import htmx from "htmx.org";
import { Chart, Colors, BarController, BarElement, Legend } from "chart.js";

Chart.register(Colors, BarController, BarElement, Legend);

class ResultChart extends HTMLElement {
  score: number = 0;
  max: number = 10;
  topTitle: string = "title";

  constructor() {
    super();

    this.score = +this.getAttribute("score")!;
    this.max = +this.getAttribute("max")!;
    this.topTitle = this.getAttribute("title")!;

    this.attachShadow({ mode: "open" });
  }

  connectedCallback() {
    this.render();
  }

  render() {
    // Create a canvas element in the shadow DOM
    const canvas = document.createElement("canvas");
    canvas.width = 400; // Set your desired width
    canvas.height = 200; // Set your desired height
    this.shadowRoot!.appendChild(canvas);

    const ctx = canvas.getContext("2d")!;

    new Chart(ctx, {
      data: {
        labels: ["Score"],
        datasets: [
          {
            data: { count: this.score, min: 0, max: this.max },
            backgroundColor: "rgba(75, 192, 192, 0.2)",
            borderColor: "rgba(75, 192, 192, 1)",
            borderWidth: 1,
          },
        ],
      },
      type: "bar",
      options: {
        indexAxis: "y",
        elements: {
          bar: {
            borderWidth: 2,
          },
        },
        responsive: true,
        plugins: {
          legend: {
            position: "right",
          },
          title: {
            display: true,
            text: this.topTitle,
          },
        },
      },
    });
  }
}

declare global {
  interface Window {
    htmx: typeof htmx;
  }
}

window.htmx = htmx;
window.customElements.define("results-chart", ResultChart);
