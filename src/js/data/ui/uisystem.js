export default class UISystem {
    constructor() {
        this.uicv = document.createElement("canvas");
        this.uictx = this.uicv.getContext("2d");

        this.uiList = [];
    }

    draw(w, h, z) {
        this.uicv.width = w;
        this.uicv.height = h;

        this.uiList.forEach(obj => {
            obj.Draw(this.uictx, z);
        });
    }
}