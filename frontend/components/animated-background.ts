export class AnimationProcessor {
  private ctx: CanvasRenderingContext2D;
  private pointsArray: any[];
  private ticker: NodeJS.Timeout | undefined;

  private canvas: HTMLCanvasElement;

  private pointsCount: number;
  private colors: string[];

  private radius: number;
  private center: { x: number; y: number };

  constructor(
    canvas: HTMLCanvasElement,
    pointsCount = 150,
    transparency = 1,
    radius = 400,
    center = { x: 0, y: 0 },
  ) {
    this.canvas = canvas;
    transparency;
    this.pointsCount = pointsCount;
    this.colors = [
      `rgba(246, 211, 101, ${transparency})`,
      `rgba(253, 160, 133, ${transparency})`,
      `rgba(246, 211, 101, ${transparency})`,
      `rgba(253, 160, 133, ${transparency})`,
      `rgba(246, 211, 101, ${transparency})`,
      `rgba(253, 160, 133, ${transparency})`,
    ];
    this.radius = radius;
    this.center = center;

    this.ctx = canvas.getContext("2d") as CanvasRenderingContext2D;

    this.resizeCanvas();

    this.pointsArray = Array.from({ length: this.pointsCount }, () => {
      return {
        x: Math.random() * canvas.width,
        y: Math.random() * canvas.height,
        radius: Math.random() * 2,
        dx: Math.random() * 2 - 1,
        dy: Math.random() * 2 - 1,
        color: this.colors[Math.floor(Math.random() * this.colors.length)],
      };
    });
    window.addEventListener("resize", this.resizeCanvas.bind(this));
  }

  private resizeCanvas() {
    this.canvas.width =
      this.canvas.parentElement?.clientWidth || window.innerWidth;
    this.canvas.height =
      this.canvas.parentElement?.clientHeight || window.innerHeight;
  }

  public updateCenter(x: number, y: number) {
    this.center = { x, y };
  }

  public updateRadius(radius: number) {
    this.radius = radius;
  }

  private draw() {
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);

    // const gradient = this.ctx.createRadialGradient(
    //   this.center.x,
    //   this.center.y,
    //   this.radius * 0.1,
    //   this.center.x,
    //   this.center.y,
    //   this.radius,
    // );

    // gradient.addColorStop(0, "rgba(255, 195, 40, 1)");
    // gradient.addColorStop(1, "rgba(253, 60, 180, 0)");

    // this.ctx.fillStyle = gradient;
    // this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);

    function createCustomRadialGradient(
      ctx: CanvasRenderingContext2D,
      centerX: number,
      centerY: number,
      innerRadius: number,
      outerRadius: number,
      innerColor: { r: number; g: number; b: number; a: number },
      outerColor: { r: number; g: number; b: number; a: number },
    ) {
      const imageData = ctx.createImageData(
        ctx.canvas.width,
        ctx.canvas.height,
      );
      const data = imageData.data;

      for (let y = 0; y < ctx.canvas.height; y++) {
        for (let x = 0; x < ctx.canvas.width; x++) {
          const dx = x - centerX;
          const dy = y - centerY;
          const distance = Math.sqrt(dx * dx + dy * dy);

          // Normalize distance to [0, 1]
          const t = Math.min(
            Math.max((distance - innerRadius) / (outerRadius - innerRadius), 0),
            1,
          );

          // Interpolate colors
          const r = innerColor.r + t * (outerColor.r - innerColor.r);
          const g = innerColor.g + t * (outerColor.g - innerColor.g);
          const b = innerColor.b + t * (outerColor.b - innerColor.b);
          const a = innerColor.a + t * (outerColor.a - innerColor.a);

          const index = (y * ctx.canvas.width + x) * 4;

          data[index] = Math.round(r);
          data[index + 1] = Math.round(g);
          data[index + 2] = Math.round(b);
          data[index + 3] = Math.round(a * 255);
        }
      }

      ctx.putImageData(imageData, 0, 0);
    }

    createCustomRadialGradient(
      this.ctx,
      this.center.x,
      this.center.y,
      this.radius * 0.1,
      this.radius,
      { r: 255, g: 195, b: 40, a: 1 },
      { r: 253, g: 60, b: 180, a: 0 },
    );

    this.pointsArray.forEach((point) => {
      this.ctx.beginPath();
      this.ctx.arc(point.x, point.y, point.radius, 0, Math.PI * 2);
      this.ctx.fillStyle = point.color;
      this.ctx.fill();
    });
  }

  private update() {
    this.pointsArray.forEach((point) => {
      point.x += point.dx;
      point.y += point.dy;

      if (point.x < -10 || point.x > this.canvas.width + 10) {
        point.dx *= -1;
      }

      if (point.y < -10 || point.y > this.canvas.height + 10) {
        point.dy *= -1;
      }
    });
  }

  public start() {
    this.ticker = setInterval(() => {
      this.update();
      this.draw();
    }, 1000 / 60); // 60 FPS
  }

  public stop() {
    if (this.ticker) {
      clearInterval(this.ticker);
      window.removeEventListener("resize", this.resizeCanvas.bind(this));
      this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
    }
  }
}
