import { drawBar, drawText } from "../../lib/draw";
import { RGB } from '../../lib/util';

export default class DefaultUI {
    constructor(x, y, w, h, mx, my, c) {
        this.x = x;
        this.y = y;
        this.w = w || 0;
        this.h = h || 0;
        this.mx = mx || "mid"; // "left" "mid" "right"
        this.my = my || "mid"; // "up" "mid" "down"
        this.color = c;
        this.childs = [];
        this.isEnable = true;
    }

    addChild(obj) {
        this.childs.push(obj);
        return this;
    }

    setEnable(isEnable) {
        this.isEnable = isEnable;
    }

    setPosition(x, y, w, h, z) {
        let pos = {
            sx: 0,
            sy: 0,
        }

        switch (this.mx) { 
            case "mid":
                pos.sx = x + w / 2 - this.w / 2 * z;
                break;
            case "right":
                pos.sx = x + w - this.w * z;
                break;
            default:
                pos.sx = x;
                break;
        }
        pos.sx += this.x * z;

        switch (this.my) {   
            case "mid":
                pos.sy = y + h / 2 - this.h / 2 * z;
                break;
            case "down":
                pos.sy = y + h - this.h * z;
                break;
            default:
                pos.sy = y;
                break;
        }
        pos.sy += this.y * z;

        return pos;
    }

    draw(ctx, x, y, w, h, z) {
        if (!this.isEnable) return;

        let {sx, sy} = this.setPosition(x, y, w, h, z);

        this.drawThis(ctx, sx, sy, z);

        this.childs.forEach((o) => {
            o.draw(ctx, sx, sy, this.w * z, this.h * z, z);
        });
    }
}

export class Bar extends DefaultUI {
    setPer(p) {
        this.per = p;
        return this;
    }

    drawThis(ctx, sx, sy, z) {
        if (!this.isEnable) return;
        drawBar(ctx, sx / (z * 3), sy / (z * 3) + this.h / 6, this.h / 6, this.w / 3, z * 3, 1, this.per, this.color);
    }
}

export class Text extends DefaultUI {
    setText(t, s, d) {
        this.text = t;
        this.size = s;
        this.dir = d;
        return this;
    }

    drawThis(ctx, sx, sy, z) {
        if (!this.isEnable) return;
        drawText(ctx, sx / z, sy / z + this.h / 2, z, 1, new RGB("#ffffff"), this.text, this.size, this.dir);
    }
}

export class List extends DefaultUI {
    setList(list) {
        list.forEach((o) => {
            this.addChild(o);
        });
        return this;
    }

    setPadding(p) {
        this.padding = p;
        return this;
    }

    draw(ctx, x, y, w, h, z) {
        if (!this.isEnable) return;
        
        let {sx, sy} = this.setPosition(x, y, w, h, z);

        this.childs.forEach((o) => {
            o.draw(ctx, sx, sy, this.w, this.h, z);
            sy += (o.h + (this.padding || 0)) * z;
        });
    }
}
