export default class DefaultUI {
    constructor(x, y, w, h, mx, my) {
        this.x = x; // 여기에서 의미하는 xy 값은 절대 좌표값이 아니라 상대 좌표값을 의미합니다.
        this.y = y;
        this.w = w || 100;
        this.h = h || 100;
        this.mx = mx || "mid"; // "left" "mid" "right" 상위 오브젝트에 대해 중심점을 잡아주는 역할을 해줍니다.
        this.my = my || "mid"; // "up" "mid" "down"
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

        drawThis(ctx, sx, sy);

        this.childs.forEach((o) => {
            o.draw(ctx, sx, sy, this.w, this.h, z);
        });
    }
}

export class Bar extends DefaultUI {
    drawThis(ctx, sx, sy) {
        ctx.save();

        

        ctx.restore();
    }
}