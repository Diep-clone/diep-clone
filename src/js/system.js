import * as data from './data';

import { Obj } from './data/object';
import { drawBackground, drawText } from './lib/draw';
import { RGB, calByte } from './lib/util';
let socket;

export default class System {
    constructor() {
        this.connect();

        this.img = new Image();
        this.img.src = 'background.png';

        this.cv = document.getElementById("canvas");
        this.ctx = this.cv.getContext("2d");

        this.uicv = document.createElement("canvas");
        this.uictx = this.uicv.getContext("2d");

        this.textinputcontainer = document.getElementById("textInputContainer");
        this.textinput = document.getElementById("textInput");

        this.textinputanime = 1;
        this.connectinga = 1;

        this.textinput.style.paddingLeft = "5px";
        this.textinput.style.paddingRight = "5px";

        this.objectList = [];
        this.uiList = [];

        this.gameSetting = {
            "gamemode": "sandbox",
            "isConnecting": true,
            "isNaming": false,
            "isGaming": false,
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
            "mousePos": {x:0,y:0},
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
                if (this.gameSetting.isNaming) {
                    if (arguments[0] === 13) {
                        this.gameSetting.isGaming = true;
                        this.textinputcontainer.style.display = "none";
                        this.textinputcontainer.style.top = "-" + this.textinputcontainer.style.top;
                        
                        var buffer = new ArrayBuffer(31);
                        var view = new DataView(buffer);

                        var string = unescape(encodeURIComponent(this.textinput.value));
                        
                        view.setUint8(0, 0);
                        for (var i = 0; i < string.length; i++) {
                            view.setUint8(i+1, string[i].charCodeAt(0));
                        }
                        socket.send(buffer);
                        
                        window['setTyping'](false);
                        localStorage['name'] = this.textinput.value;
                        this.gameSetting.isNaming = false;
                        return;
                    } else {
                        return;
                    }
                } else if (this.gameSetting.isConnecting) return;

                if (this.sameKeys[arguments[0]]) arguments[0]=this.sameKeys[arguments[0]];
                if (this.keys[arguments[0]]) return;
                this.keys[arguments[0]] = true;
                switch (arguments[0]){
                    case 37:
                        this.input.moveVector.x-=1;
                        break;
                    case 38:
                        this.input.moveVector.y-=1;
                        break;
                    case 39:
                        this.input.moveVector.x+=1;
                        break;
                    case 40:
                        this.input.moveVector.y+=1;
                        break;
                    default:
                        break;
                }
            }.bind(this),
            keyUp:function(){
                if (this.gameSetting.isNaming) {
                    return;
                } else if (this.gameSetting.isConnecting) return;

                if (this.sameKeys[arguments[0]]) arguments[0]=this.sameKeys[arguments[0]];
                if (!this.keys[arguments[0]]) return;
                this.keys[arguments[0]] = false;
                switch (arguments[0]){
                    case 37:
                        this.input.moveVector.x+=1;
                        break;
                    case 38:
                        this.input.moveVector.y+=1;
                        break;
                    case 39:
                        this.input.moveVector.x-=1;
                        break;
                    case 40:
                        this.input.moveVector.y-=1;
                        break;
                    default:
                        break;
                }
            }.bind(this),
            mouse:function(){
                this.input.mousePos = {x:arguments[0], y:arguments[1]};
            }.bind(this),
            prevent_right_click: function(){
                return true;
            },
            print_convar_help: function(){},
            set_convar: function(key,value){},
            should_prevent_unload: function(){
                return true;
            },
            wheel: function(){
                return true;
            },
        };

        this.loop();
    }

    connect() {
        socket = new WebSocket(`${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/ws`);

        socket.binaryType = 'arraybuffer';

        socket.onopen = () => {
            console.log("Successfully Connected");
            this.gameSetting.isNaming = true;
            window['setTyping'](true);
            this.textinput.value = localStorage['name'] || '';
            this.textinputcontainer.style.display = "block";
            this.gameSetting.isConnecting = false;
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
                    if (this.gameSetting.isNaming) {
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
            this.gameSetting.isConnecting = true;
            this.gameSetting.isNaming = false;
            this.gameSetting.isGaming = false;
            this.textinputcontainer.style.display = "none";
            this.textinputcontainer.style.top = "-" + this.textinputcontainer.style.top;
            console.log("Socket Closed Connection: ", event);
            this.connect();
        };

        socket.onerror = error => {
            console.error("Socket Error: ", error);
            socket.close();
        };
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

        if (this.gameSetting.isConnecting) {
            this.textinputanime = 1;
            this.connectinga = Math.min(this.connectinga + 0.2,1);
        }
        
        if (this.gameSetting.isGaming) {
            drawBackground(this.ctx, this.camera.x, this.camera.y, this.camera.z, this.cv.width, this.cv.height, this.area);

            var buffer = new ArrayBuffer(14);
            var view = new DataView(buffer);
            view.setUint8(0, 1);
            if (this.input.moveVector.x === 0 && this.input.moveVector.y === 0) {
                view.setFloat32(1, 9.);
            } else {
                view.setFloat32(1, Math.atan2(this.input.moveVector.y,this.input.moveVector.x));
            }
            view.setFloat32(5,this.input.mousePos.x/this.camera.z+this.camera.x);
            view.setFloat32(9,this.input.mousePos.y/this.camera.z+this.camera.y);
            view.setUint8(13,
                ((this.keys[1] || this.keys[32]) || 0)
                + ((this.keys[3] || this.keys[16]) || 0) * 2
                + (this.keys[79] || 0) * 4
            );

            socket.send(buffer);

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
        } else {
            this.ctx.drawImage(this.img,
                this.cv.width / 2 - this.img.width * this.camera.uiz / 2 / 2.4,
                this.cv.height / 2 - this.img.height * this.camera.uiz / 2 / 2.4,
                this.img.width * this.camera.uiz / 2.4,
                this.img.height * this.camera.uiz / 2.4);
            drawText(this.ctx,this.cv.width / 2 / this.camera.uiz, this.cv.height / 2 / this.camera.uiz, this.camera.uiz, this.connectinga, new RGB("#FFFFFF"), "Connecting...", 60);
        }

        if (this.gameSetting.isNaming) {
            if (this.textinput.value) this.textinput.value = calByte.cutByteLength(this.textinput.value,15);

            let x = this.cv.width / 2 - 166 * this.camera.uiz,
            y = (this.cv.height / 2 - 21 * this.camera.uiz) * (1-this.textinputanime),
            w = 332 * this.camera.uiz,
            h = 42 * this.camera.uiz;

            this.textinputanime *= 0.95;
            this.connectinga = Math.max(this.connectinga - 0.2,0);

            this.textinputcontainer.style.left = window['unscale'](x) + "px";
            this.textinputcontainer.style.top = window['unscale'](y) + "px";

            this.textinput.style.width = window['unscale'](w) + "px";
            this.textinput.style.height = window['unscale'](h) + "px";
            this.textinput.style.fontSize = this.textinput.style.lineHeight = window['unscale'](h - 0.4) + "px";

            drawText(this.ctx, (x + w / 2) / this.camera.uiz, y / this.camera.uiz - 11, this.camera.uiz, 1, new RGB("#FFFFFF"), "This is the tale of...", 20.2);
            drawText(this.ctx, (x + w / 2) / this.camera.uiz, y / this.camera.uiz + 57, this.camera.uiz, 1, new RGB("#FFFFFF"), "(press enter to spawn)", 11.8);

            this.ctx.save();

            this.ctx.beginPath();

            this.ctx.lineCap = "round";
            this.ctx.lineJoin = "round";
            this.ctx.lineWidth = 4.5 * this.camera.uiz;
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
        }

        requestAnimationFrame(this.loop.bind(this));
    }
}
