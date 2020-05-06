export default class UISystem {
    constructor() {
        this.uicv = document.createElement("canvas");
        this.uictx = this.uicv.getContext("2d");

        this.uiList = [];
    }

    draw(ctx, w, h, z) {
        this.uicv.width = w;
        this.uicv.height = h;

        this.uiList.forEach(obj => {
            
        });

        ctx.save();
        ctx.globalAlpha = 0.7;
        ctx.drawImage(this.uicv, 0, 0, w, h);
        ctx.restore();
    }
}