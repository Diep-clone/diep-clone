import { uiOpacity } from './console';

export default class UISystem {
    constructor() {
        this.uicv = document.createElement("canvas");
        this.uictx = this.uicv.getContext("2d");

        this.uiList = [];

        this.w = 1000;
        this.h = 1000;
    }

    draw(ctx, w, h, z) {
        this.uicv.width = this.w = w;
        this.uicv.height = this.h = h;

        this.uiList.forEach(obj => {
            obj.draw(this.uictx, 0, 0, w, h, z);
        });

        ctx.save();
        ctx.globalAlpha = uiOpacity;
        ctx.drawImage(this.uicv, 0, 0, w, h);
        ctx.restore();
    }
}