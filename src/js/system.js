import * as data from './data';

import { Obj } from './data/object';
import { drawBackground, drawText } from './lib/draw';
import { RGB } from './lib/util';
const socket = new WebSocket(`${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/ws`);
console.log(socket);

export default class System {
    constructor() {
        this.cv = document.getElementById("canvas");
        this.ctx = this.cv.getContext("2d");

        this.uicv = document.createElement("canvas");
        this.uictx = this.uicv.getContext("2d");

        this.textinputcontainer = document.getElementById("textInputContainer");
        this.textinput = document.getElementById("textInput");

        this.objectList = [];
        this.uiList = [];

        this.gameSetting = {
            "gamemode": "sandbox",
            "gameset": "Connecting",
        };

        this.playerSetting = {
            "id": "",
            "level": 0,
            "isCanRotate": false,
            "stat": 0,
            "stats": [],
            "maxstats": [],
        };

        this.input = {
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

        this.area = [{X:0,Y:0,W:0,H:0}];

        this.keys = {};
        this.sameKeys = {
            65:37,
            87:38,
            68:39,
            83:40,
        };

        window.input = {
            blur:function(){},
            execute:function(v){},
            flushInputHooks: function(){},
            get_convar:function(key){},
            keyDown: function(){
                if (this.gameSetting.gameset === "SetName") {
                    if (arguments[0] === 13) {
                        this.gameSetting.gameset = "Gaming";
                        this.textinputcontainer.style.top = "-100px";
                        this.socketSend("init",this.textinput.value || "");
                        window['setTyping'](false);
                        return;
                    } else {
                        return;
                    }
                } else if (this.gameSetting.gameset === "Connecting") return;

                if (this.sameKeys[arguments[0]]) arguments[0]=this.sameKeys[arguments[0]];
                if (this.keys[arguments[0]]) return;
                this.keys[arguments[0]] = true;
                let key = "";
                switch (arguments[0]){
                    case 1:
                    case 32:
                        key = "mouseLeft";
                        break;
                    case 3:
                    case 16:
                        key = "mouseRight";
                        break;
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
                    case 79:
                        key = "suicide";
                        break;
                    default:
                        break;
                }
                let value = true;
                switch (key){
                    case "mouseLeft":
                        value = this.keys[1] || this.keys[32];
                        break;
                    case "mouseRight":
                        value = this.keys[3] || this.keys[16];
                        break;
                    case "moveVector":
                        if (this.input.moveVector.x === 0 && this.input.moveVector.y === 0) value = 9;
                        else value = Math.atan2(this.input.moveVector.y,this.input.moveVector.x);
                        break;
                    default:
                        break;
                }
                this.socketSend("input",{"type":key,"value":value})
            }.bind(this),
            keyUp:function(){
                if (this.gameSetting.gameset === "SetName") {
                    return;
                } else if (this.gameSetting.gameset === "Connecting") return;

                if (this.sameKeys[arguments[0]]) arguments[0]=this.sameKeys[arguments[0]];
                if (!this.keys[arguments[0]]) return;
                this.keys[arguments[0]] = false;
                let key = "";
                switch (arguments[0]){
                    case 1:
                    case 32:
                        key = "mouseLeft";
                        break;
                    case 3:
                    case 16:
                        key = "mouseRight";
                        break;
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
                    case "mouseLeft":
                        value = this.keys[1] || this.keys[32];
                        break;
                    case "mouseRight":
                        value = this.keys[3] || this.keys[16];
                        break;
                    case "moveVector":
                        if (this.input.moveVector.x === 0 && this.input.moveVector.y === 0) value = 9;
                        else value = Math.atan2(this.input.moveVector.y,this.input.moveVector.x);
                        break;
                    default:
                        break;
                }
                this.socketSend("input",{"type":key,"value":value})
            }.bind(this),
            mouse:function(){
                this.socketSend("input",{"type":"mouseMove","value":{
                    "x":arguments[0]/this.camera.z+this.camera.x,
                    "y":arguments[1]/this.camera.z+this.camera.y}})
            }.bind(this),
            prevent_right_click: function(){
                return true;
            },
            print_convar_help: function(){},
            set_convar: function(key,value){},
            should_prevent_unload: function(){
                return true;
            },
            wheel: function(){}.bind(this),
        };

        socket.onopen = () => {
            console.log("Successfully Connected");
            this.gameSetting.gameset = "SetName";
            window['setTyping'](true);
        };

        socket.onmessage = msg => {
            const json = JSON.parse(msg.data);

            const event = json.event;
            const data = json.data;

            switch (event) {
                case 'init': {
                    console.log('init you');
                    break;
                }
                case 'objectList': {
                    if (this.gameSetting.gameset === "SetName") {
                        return;
                    }

                    this.camera.z = this.camera.uiz * json.camera.Z;
    
                    this.camera.x = json.camera.Pos.X - this.cv.width / 2 / this.camera.uiz / json.camera.Z;
                    this.camera.y = json.camera.Pos.Y - this.cv.height / 2 / this.camera.uiz / json.camera.Z;

                    data.forEach((obj) => {
                        
                        let isObjEnable = false;
                        
                        this.objectList.forEach((obi) => {
                            if (obi.id === obj.id){
                                obi.ObjSet(obj);
                                isObjEnable = true;
                            }
                        });
                        if (!isObjEnable && !obj.isDead) {
                            let obi = new Obj(obj.id);
                            console.log("new Object " + obj.id);
                            obi.ObjSet(obj);
                            this.objectList.push(obi);
                        }
                    });
                    this.objectList.forEach((obj) => {
                        if (!obj.isEnable){
                            console.log("delete Object " + obj.id);
                            obj.isDelete = true;
                        }
                        obj.isEnable = false;
                    })
                    break;
                }
                case 'playerSet': {
                    this.playerSetting.id = json.id;
                    this.playerSetting.level = json.level;
                    this.playerSetting.isCanRotate = json.isCanRotate;
                    this.playerSetting.stat = json.stat;
                    this.playerSetting.stats = json.stats;
                    this.playerSetting.maxstats = json.maxstats;
                    break;
                }
                case 'area': {
                    this.area = data;
                    break;
                }
            }
        };
        
        socket.onclose = event => {
            this.gameSetting.gameset = "Connecting";
            this.textinputcontainer.style.top = "-100px";
            console.log("Socket Closed Connection: ", event);
        };

        socket.onerror = error => {
            console.error("Socket Error: ", error);
        };

        this.loop();
    }

    socketSend(type,data) {
        if (this.gameSetting.gameset === "Gaming") {
            const body = JSON.stringify({"event":type,"data":data});
            socket.send(body);
        }
    }

    insertComma = (number) => number.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');

    loop() {
        const tick = Date.now() - this.lastTime;
        this.lastTime = Date.now();

        if (this.cv.width <= this.cv.height/9*16){
            this.camera.uiz = this.cv.height / 900;
        } else {
            this.camera.uiz = this.cv.width / 1600;
        }

        this.uicv.width = this.cv.width;
        this.uicv.height = this.cv.height;

        this.ctx.clearRect(0,0,this.cv.width,this.cv.height);

        switch (this.gameSetting.gameset){
            case "Connecting":
                drawText(this.ctx,this.cv.width / 2 / this.camera.uiz, this.cv.height / 2 / this.camera.uiz, this.camera.uiz, 1, new RGB("#FFFFFF"), "Connecting...", 40 * this.camera.uiz);
                break;
            case "SetName":
                let x = this.cv.width / 2 - 160 * this.camera.uiz,
                y = this.cv.height / 2 - 20 * this.camera.uiz,
                w = 320 * this.camera.uiz,
                h = 40 * this.camera.uiz;

                this.textinputcontainer.style.left = x + "px";
                this.textinputcontainer.style.top = y + "px";

                this.textinput.style.width = w + "px";
                this.textinput.style.height = h + "px";
                this.textinput.style.fontSize = this.textinput.style.lineHeight = h - 0.4 + "px";

                this.ctx.save();

                this.ctx.beginPath();

                this.ctx.lineWidth = 5 * this.camera.uiz;
                this.ctx.fillStyle = "#FFFFFF";
                this.ctx.strokeStyle = "#000000";

                this.ctx.moveTo(x, y);
                this.ctx.lineTo(x + w, y);
                this.ctx.lineTo(x + w, y + h);
                this.ctx.lineTo(x, y + h);
                this.ctx.lineTo(x, y);

                this.ctx.fill();
                this.ctx.stroke();

                this.ctx.closePath();

                this.ctx.restore();

                break;
            case "Gaming":
                drawBackground(this.ctx, this.camera.x, this.camera.y, this.camera.z, this.cv.width, this.cv.height, this.area);

                for (let i=0;i<this.objectList.length;){
                    if (this.objectList[i].isDelete){
                        this.objectList.splice(i,1);
                    } else {
                        i++;
                    }
                }

                this.objectList.forEach((o) => {
                    o.Animate(tick);
                });

                this.objectList.forEach((o) => {
                    o.Draw(this.ctx, this.camera);
                });

                this.objectList.forEach((o) => {
                    o.DrawName(this.ctx, this.camera);
                });

                this.objectList.forEach((o) => {
                    o.DrawHPBar(this.ctx, this.camera);
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
}
