import * as data from './data';

import io from 'socket.io-client';
import { Obj } from './data/object';
import { drawBackground } from './lib/draw';

export const Socket = io();

export default class System {
    constructor() {
        this.cv = document.getElementById("canvas");
        this.ctx = this.cv.getContext("2d");

        this.namecv = document.createElement("canvas");
        this.namectx = this.namecv.getContext("2d");

        this.hpcv = document.createElement("canvas");
        this.hpctx = this.hpcv.getContext("2d");

        this.uicv = document.createElement("canvas");
        this.uictx = this.uicv.getContext("2d");

        this.objectList = [];
        this.uiList = [];

        this.gameSetting = {
            "gameMode": "sandbox",
            "gameSet": "Loading",
        };

        this.playerSetting = {
            "id": -1,
            "level": 0,
            "sight": 1,
            "isCanRotate": false,
            "stat": 0,
            "stats": [],
            "maxStats": [],
        };

        this.lastTime = Date.now();
        this.isControlRotate = true;

        this.camera = {
            x: 0,
            y: 0,
            z: 2,
        }

        window.input = data.input;

        Socket.on("playerSet", function (data, camera) {
            this.playerSetting = data;
            this.camera = camera;
        });

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

        console.log(this.playerSetting);

        switch (this.gameSetting.gameSet){
            case "Connecting":
                break;
            case "Gaming":
                drawBackground(this.ctx, this.camera.x, this.camera.y, this.camera.z, this.cv.clientWidth, this.cv.clientWidth, []);

                this.objectList.forEach((o) => {
                    o.Animite(tick);
                    o.Draw(this.ctx, this.camera);
                    o.DrawName(this.namectx, this.camera);
                    o.DrawHPBar(this.hpctx, this.camera);
                });

                this.uiList.forEach((ui) => {
                    ui.Draw(this.uictx);
                });
                break;
            default:
                break;
        }

        requestAnimationFrame(this.loop.bind(this));
    }

    RemoveObject(id) {
        for (let i=0;i<this.objectList.length;i++){
            if (this.objectList[i].id === id){
                this.objectList.splice(i,1);
            }
        }
    }

    Resize() {
        this.cv.width =
        this.namecv.width =
        this.hpcv.width =
        this.uicv.width = window.innerWidth * window.devicePixelRatio;

        this.cv.height =
        this.namecv.height =
        this.hpcv.height =
        this.uicv.height = window.innerHeight * window.devicePixelRatio;

        this.hpctx.lineCap = this.uictx.lineCap = "round";
        this.hpctx.lineJoin = this.uictx.lineJoin = "round";

        this.ctx.imageSmoothingEnabled = 
        this.namectx.imageSmoothingEnabled = 
        this.hpctx.imageSmoothingEnabled = 
        this.uictx.imageSmoothingEnabled = false;
    }
}