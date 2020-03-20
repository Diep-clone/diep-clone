import * as data from './data';

import io from 'socket.io-client';
import { Obj } from './data/object';
import { drawBackground } from './lib/draw';

export var Socket = io();

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
            "gamemode": "sandbox",
            "gameset": "Gaming",
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

        this.inputSetting = {
            "moveVector": {x:0,y:0},
        }

        this.lastTime = Date.now();
        this.isControlRotate = true;

        this.camera = {
            x: 0,
            y: 0,
            z: 2,
            uiz: 1,
        }

        this.keys = {};

        window.input = {
            blur:function(){},
            execute:function(v){},
            flushInputHooks: function(){},
            get_convar:function(key){},
            keyDown: function(){
                if (this.keys[arguments[0]]) return;
                this.keys[arguments[0]] = true;
                let key = "";
                switch (arguments[0]){
                    case 37:
                        this.input.moveVector.x-=1;
                        key = "moveVector";
                        break;
                    case 38:
                        this.input.moveVector.y-=1;
                        key = "moveVector";
                        break;
                    case 39:
                        this.input.moveVector.x+=1;
                        key = "moveVector";
                        break;
                    case 40:
                        this.input.moveVector.y+=1;
                        key = "moveVector";
                        break;
                    default:
                        break;
                }
                let value = true;
                switch (key){
                    case "moveVector":
                        value = Math.atan2(this.input.moveVector.y,this.input.moveVector.x);
                        break;
                    default:
                        break;
                }
                Socket.emit(key, value);
            },
            keyUp:function(){
                if (!this.keys[arguments[0]]) return;
                this.keys[arguments[0]] = false;
                let key = "";
                switch (arguments[0]){
                    case 37:
                        this.input.moveVector.x+=1;
                        key = "moveVector";
                        break;
                    case 38:
                        this.input.moveVector.y+=1;
                        key = "moveVector";
                        break;
                    case 39:
                        this.input.moveVector.x-=1;
                        key = "moveVector";
                        break;
                    case 40:
                        this.input.moveVector.y-=1;
                        key = "moveVector";
                        break;
                    default:
                        break;
                }
                let value = false;
                switch (key){
                    case "moveVector":
                        value = Math.atan2(this.input.moveVector.y,this.input.moveVector.x);
                        break;
                    default:
                        break;
                }
                Socket.emit(key, value);
            },
            mouse:function(){},
            prevent_right_click: function(){},
            print_convar_help: function(){},
            set_convar: function(key,value){},
            should_prevent_unload: function(){},
            wheel: function(){},
        };

        Socket.emit("login");

        Socket.on("playerSet", function (data, camera) {
            this.playerSetting = data;

            if (this.cv.width <= this.cv.height/9*16){
                this.camera.z = this.cv.height / 900;
            } else {
                this.camera.z = this.cv.width / 1600;
            } 

            this.camera.uiz = this.camera.z;

            this.camera.z *= camera.Z;

            this.camera.x = camera.Pos.X - this.cv.width / 2 / this.camera.uiz / camera.Z;
            this.camera.y = camera.Pos.Y - this.cv.height / 2 / this.camera.uiz / camera.Z;
        }.bind(this));

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
        }.bind(this));
    }

    insertComma = (number) => number.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');

    loop() {
        const tick = Date.now() - this.lastTime;
        this.lastTime = Date.now();

        this.namecv.width =
        this.hpcv.width =
        this.uicv.width = this.cv.width;

        this.namecv.height =
        this.hpcv.height =
        this.uicv.height = this.cv.height;

        switch (this.gameSetting.gameset){
            case "Connecting":
                break;
            case "Gaming":
                drawBackground(this.ctx, this.camera.x, this.camera.y, this.camera.z, this.cv.width, this.cv.height, [{x:-1000,y:-1000,w:2000,h:2000}]);

                this.objectList.forEach((o) => {
                    o.Animate(tick);
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
}