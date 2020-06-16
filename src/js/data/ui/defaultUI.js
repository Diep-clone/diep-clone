export default class DefaultUI {
    constructor(parent, w, h, mx, my) {
        this.w = w || 100;
        this.h = h || 100;
        this.mx = mx || "mid"; // "left" "mid" "right"
        this.my = my || "mid"; // "up" "mid" "down"
        this.parent = parent;
    }

    setPosition(x, y, w, h) {
        let pos = {
            x: 0,
            y: 0,
        }

        switch (this.mx) { 
            case "mid":
                pos.x = x + w / 2 - this.w / 2;
                break;
            case "right":
                pos.x = x + w - this.w;
                break;
            default:
                pos.x = x;
                break;
        }

        switch (this.my) {   
            case "mid":
                pos.y = y + h / 2 - this.h / 2;
                break;
            case "down":
                pos.y = y + h - this.h;
                break;
            default:
                pos.y = y;
                break;
        }

        return pos;
    }
}