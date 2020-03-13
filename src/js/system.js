import * as data from './data';

import io from 'socket.io-client';
import { Obj } from './data/object';
import { drawBackground } from './lib/draw';

export const Socket = io();

export default class System {
    constructor() {
        this.cv = document.getElementById("canvas");
        this.ctx = this.cv.getContext("2d");

        this.statcv = document.createElement("canvas");
        this.statctx = this.statcv.getContext("2d");

        this.uicv = document.createElement("canvas");
        this.uictx = this.uicv.getContext("2d");

        this.objectList = [];
        this.uiObjectList = [];

        this.colorList = data.colorList;
        this.expList = data.expList;

        this.lastTime = Date.now();
        this.isControlRotate = true;

        this.camera = {
            x: 0,
            y: 0,
            z: 2,
        }

        window.input = data.input;

        Socket.on("objectList", function (list) {
            list.forEach((obj) => {
                let isObjEnable = false;
                this.objectList.forEach((obi) => {
                    if (obi.id === obj.id){
                        obi.ObjSet(obj);
                        isObjEnable = true;
                    }
                });
                if (!isObjEnable) {
                    let obi = new Obj(obj.id);
                    obi.ObjSet(obj);
                    this.objectList.push(obi);
                }
            });
        });
    }

    insertComma = (number) => number.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');

    loop() {
        const tick = Date.now() - this.lastTime;
        this.lastTime = Date.now();

        drawBackground(this.ctx, this.camera.x, this.camera.y, this.camera.z, this.cv.clientWidth, this.cv.clientWidth, []);

        this.objectList.forEach((o) => {
            o.Animite(tick);
            o.Draw(this.ctx, this.camera);
            o.DrawStatus(this.statctx, this.camera);
        });

        

        requestAnimationFrame(this.loop.bind(this));
    }
}