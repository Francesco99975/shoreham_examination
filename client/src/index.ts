import "./css/style.css";

import htmx from "htmx.org";
import {
  Chart,
  Colors,
  BarController,
  BarElement,
  Legend,
  CategoryScale,
  LinearScale,
} from "chart.js";

import Annotation from "chartjs-plugin-annotation";

Chart.register(
  Colors,
  BarController,
  BarElement,
  Annotation,
  Legend,
  CategoryScale,
  LinearScale
);

class ResultChart extends HTMLElement {
  test: string = "";
  score: number = 0;
  max: number = 0;
  low: number = 0;
  high: number = 0;
  chart!: Chart<"bar", number[], string>;

  constructor() {
    super();

    const id = this.getAttribute("id")!;
    this.setAttribute("id", this.id + "pr");
    this.test = this.getAttribute("test")!;
    this.score = +this.getAttribute("score")!;
    this.max = +this.getAttribute("max")!;
    this.low = +this.getAttribute("low")!;
    this.high = +this.getAttribute("high")!;
    this.attachShadow({ mode: "open" });
  }

  render(id: string) {
    this.shadowRoot!.innerHTML = `
            <style>
              div {
                width: 100%;
                display: flex;
                justify-content: center;
              }
              canvas {
                width: 100%;
                height: 300px;
              }
            </style>
           <div>
               <canvas id="${id}"></canvas>
           </div>

          `;

    const canvas = this.shadowRoot!.getElementById(id) as HTMLCanvasElement;

    const ctx = canvas.getContext("2d");

    if (ctx) {
      this.chart = new Chart(ctx, {
        data: {
          labels: ["Patient Score"],
          datasets: [
            {
              label: "Score",
              data: [this.score],
              backgroundColor: "#22c55e",
              borderColor: "#14532d",
              borderWidth: 1,
              barThickness: 25,
              type: "bar",
            },
          ],
        },
        type: "bar",
        options: {
          scales: {
            y: {
              beginAtZero: true,
              grid: {
                color: "#eeeee4",
              },
              ticks: {
                color: "#eeeee4",
              },
            },
            x: {
              beginAtZero: true,
              max: this.max,
              grid: {
                color: "#eeeee4",
              },
              ticks: {
                color: "#eeeee4",
              },
            },
          },

          indexAxis: "y",
          elements: {
            bar: {
              borderWidth: 2,
              backgroundColor: "#22c55e",
              borderColor: "#14532d",
            },
          },
          responsive: true,
          plugins: {
            legend: {
              position: "right",
              labels: {
                color: "#eeeee4",
              },
            },
            title: {
              display: false,
            },
            annotation: {
              annotations: [
                {
                  type: "box",
                  drawTime: "beforeDraw",
                  backgroundColor: "rgba(9, 188, 138, 0.6)",
                  borderWidth: 0,
                  xMin: 0,
                  xMax: this.low,
                  label: {
                    drawTime: "afterDraw",
                    color: "#eeeee4",
                    display: false,
                    content: "Normal Range",
                    position: { x: "center", y: "end" },
                  },
                },
                {
                  type: "box",
                  drawTime: "beforeDraw",
                  backgroundColor: "rgba(255, 225, 168, 0.6)",
                  borderWidth: 0,
                  xMin: this.low,
                  xMax: this.high,
                  label: {
                    drawTime: "afterDraw",
                    color: "#eeeee4",
                    display: false,
                    content: "Moderate Range",
                    position: { x: "center", y: "end" },
                  },
                },
                {
                  type: "box",
                  drawTime: "beforeDraw",
                  backgroundColor: "rgba(167, 38, 8, 0.6)",
                  borderWidth: 0,
                  xMin: this.high,
                  xMax: this.max,
                  label: {
                    drawTime: "afterDraw",
                    color: "#eeeee4",
                    display: true,
                    content: "Severe Range",
                    position: { x: "center", y: "end" },
                  },
                },
              ],
            },
          },
        },
      });
    }
  }
}

declare global {
  interface Window {
    htmx: typeof htmx;
  }
}

window.htmx = htmx;
window.customElements.define("results-chart", ResultChart);
