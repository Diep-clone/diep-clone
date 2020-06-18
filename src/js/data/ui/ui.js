import { drawBar, drawText } from "../../lib/draw";

export default class DefaultUI {
    constructor(x, y, w, h, mx, my, c) {
        this.x = x; // 여기에서 의미하는 xy 값은 절대 좌표값이 아니라 상대 좌표값을 의미합니다.
        this.y = y;
        this.w = w || 100;
        this.h = h || 100;
        this.mx = mx || "mid"; // "left" "mid" "right" 상위 오브젝트에 대해 중심점을 잡아주는 역할을 해줍니다.
        this.my = my || "mid"; // "up" "mid" "down"
        this.color = c;
        this.childs = [];
    }

    addChild(obj) {
        this.childs.push(obj);
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
        pos.x += this.x;

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
        pos.y += this.y;

        return pos;
    }

    draw(ctx, x, y, w, h, z) {
        let {sx, sy} = this.setPosition(x, y, w, h);

        drawThis(ctx, sx, sy, z);

        this.childs.forEach((o) => {
            o.draw(ctx, sx, sy, this.w, this.h, z);
        });
    }
}

export class Bar extends DefaultUI {
    setPer(p) {
        this.per = p;
    }

    drawThis(ctx, sx, sy, z) {
        drawBar(ctx, sx, sy - this.h / 2, this.h / 2, this.w, z, 1, this.per, this.color);
    }
}

export class Text extends DefaultUI {
    setText(t, s, d) {
        this.text = t;
        this.size = s;
        this.dir = d;
    }

    drawThis(ctx, sx, sy, z) {
        drawText(ctx, sx + w / 2, sy + h / 2, z, 1, new RGB("#000000"), this.text, this.size, this.dir);
    }
}

export class List extends DefaultUI {
    setList(list) {
        list.forEach((o) => {
            this.addChild(o);
        });
    }

    setPadding(p) {
        this.padding = p;
    }

    draw(ctx, x, y, w, h, z) {
        let {sx, sy} = this.setPosition(x, y, w, h);

        this.childs.forEach((o) => {
            o.draw(ctx, sx, sy, this.w, this.h, z);
            sx += o.h + this.padding;
        });
    }
}
