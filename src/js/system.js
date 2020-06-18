import * as data from './data/console';

import { Obj } from './data/object';
import UISystem from './data/ui/uisystem';
import { drawBackground, drawText } from './lib/draw';
import { RGB, calByte } from './lib/util';
let socket;

/*
    전체 게임 시스템을 담당하고 있습니다.
    코드가 전체적으로 더럽습니다. 정돈이 필요한 파일입니다.
*/

export default class System {
    constructor() {
        this.connect(); // connect

        this.img = new Image(); // main background
        this.img.src = 'background.png';

        this.cv = document.getElementById("canvas"); // canvas
        this.ctx = this.cv.getContext("2d"); // ctx

        this.gameui = new UISystem();

        this.textinputcontainer = document.getElementById("textInputContainer"); // name input
        this.textinput = document.getElementById("textInput");

        this.textinputanime = 1; // name input animation (down)
        this.connectinga = 1;
        this.panela = 0.4;

        this.textinput.style.paddingLeft = "5px";
        this.textinput.style.paddingRight = "5px";

        this.objectList = []; // object List

        this.scoreboard = []; // scoreboard

        this.gameSetting = { // Game's data
            "gamemode": "sandbox",
            "isConnecting": true,
            "isNaming": false,
            "isGaming": false,
        };

        this.playerSetting = { // my Tank's data
            "obj": null, // object
            "id": 0, // object's id
            "team": "",
            "level": 0,
            "stat": 0,
            "stats": [0,0,0,0,0,0,0,0],
            "maxstats": [0,0,0,0,0,0,0,0],
        };

        this.input = { // input data
            "moveVector": {x:0,y:0},
            "mousePos": {x:0,y:0},
            "e":0,
            "c":0,
        }

        this.lastTime = Date.now();
        this.isControlRotate = true; // it is possible to control own tank's rotate?
        this.pingList = [];

        this.camera = { // camera data
            x: 0,
            y: 0,
            z: 2,
            uiz: 1, // ui scale (affected by screen size)
        }

        this.area = []; // area list

        this.keys = {}; // list of the keys currently pressed
        this.sameKeys = { // considered as the same key
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
            keyDown: function() {
                if (this.gameSetting.isNaming) {
                    if (arguments[0] === 13) { //
                        this.gameSetting.isGaming = true; // game start
                        this.textinputcontainer.style.display = "none"; // hide name input
                        this.textinputcontainer.style.top = "-" + this.textinputcontainer.style.top;
                        
                        let buffer = new ArrayBuffer(31);
                        let view = new DataView(buffer);

                        let string = unescape(encodeURIComponent(this.textinput.value));
                        
                        view.setUint8(0, 2);
                        for (let i = 0; i < string.length; i++) {
                            view.setUint8(i+1, string[i].charCodeAt(0));
                        }

                        socket.send(buffer);
                        
                        window['setTyping'](false); // 
                        localStorage['name'] = this.textinput.value; // save the name to localStorage
                        this.gameSetting.isNaming = false;
                        return;
                    } else {
                        return;
                    }
                } else if (this.gameSetting.isConnecting) return; // When naming or connecting, ignore all key pressing

                if (this.sameKeys[arguments[0]]) arguments[0]=this.sameKeys[arguments[0]];
                if (this.keys[arguments[0]]) return;
                this.keys[arguments[0]] = true;
                switch (arguments[0]) {
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
                    case 69:
                        this.input.e = 1 - this.input.e;
                        break;
                    case 67:
                        this.input.c = 1 - this.input.c;
                        break;
                    default:
                        break;
                }
            }.bind(this),
            keyUp:function() {
                if (this.gameSetting.isNaming) { // When naming or connecting, ignore all key pressing
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
            mouse:function() {
                this.input.mousePos = {x:arguments[0], y:arguments[1]}; // mouse pos
            }.bind(this),
            prevent_right_click: function() {
                return true;
            },
            print_convar_help: function(){},
            set_convar: function(key,value){},
            should_prevent_unload: function() {
                return true;
            },
            wheel: function() {
                return true;
            },
        };

        this.loop(); // start loop

        setInterval(() => {
            this.pingList.push(Date.now()); // piong
            let buffer = new ArrayBuffer(1);
            let view = new DataView(buffer);

            view.setUint8(0, 3);

            socket.send(buffer);
        }, 100);
    }

    connect() { // connect to websocket
        socket = new WebSocket(`${window.location.protocol === 'https:' ? 'wss' : 'ws'}://${location.host}/ws`);

        socket.binaryType = 'arraybuffer';

        socket.onopen = () => {
            console.log("Successfully Connected");
            let buffer = new ArrayBuffer(1);
            let view = new DataView(buffer);
            view.setUint8(0, 0);
            socket.send(buffer);

            this.gameSetting.isNaming = true;
            window['setTyping'](true);
            this.textinput.value = localStorage['name'] || '';
            this.textinputcontainer.style.display = "block";
            this.gameSetting.isConnecting = false;
        };

        socket.onmessage = msg => {
            let view = new DataView(msg.data);

            switch (view.getUint8(0)){
                case 0: { // get Camera value, objectList value
                    if (!this.gameSetting.isGaming) {
                        return;
                    }
                    this.camera.z = this.camera.uiz * view.getFloat64(17);
    
                    this.camera.x = view.getFloat64(1) - this.cv.width / 2 / this.camera.z;
                    this.camera.y = view.getFloat64(9) - this.cv.height / 2 / this.camera.z;

                    let i = 32;

                    if (view.getInt8(25)) {
                        this.playerSetting.id = view.getUint32(26);
                        this.playerSetting.level = view.getUint16(30);
                        this.playerSetting.stat = view.getUint8(32);
                        for (let j = 0; j < 8; j++) {
                            this.playerSetting.stats[j] = view.getUint8(33+j);
                        }
                        for (let j = 0; j < 8; j++) {
                            this.playerSetting.maxstats[j] = view.getUint8(41+j);
                        }
                        i = 49;
                        let len = view.getUint8(i);
                        i++;
                        let u = new Uint8Array(msg.data.slice(i,i+len));
                        this.playerSetting.team = String.fromCharCode.apply(null, u);
                        i+=len;
                    }

                    let imObjList = [];

                    while (i < msg.data.byteLength) {
                        let isObjEnable = false;

                        let obj = {};

                        obj.id = view.getUint32(i);
                        i+=4;

                        obj.x = view.getFloat64(i);
                        i+=8;

                        obj.y = view.getFloat64(i);
                        i+=8;

                        obj.r = view.getFloat64(i);
                        i+=8;

                        obj.dir = view.getFloat64(i);
                        i+=8;

                        obj.mh = view.getFloat64(i);
                        i+=8;

                        obj.h = view.getFloat64(i);
                        i+=8;

                        obj.opacity = view.getFloat64(i);
                        i+=8;

                        obj.score = view.getUint32(i);
                        i+=4;

                        let gunLen = view.getUint8(i);
                        i++;
                        obj.guns = [];
                        for (; gunLen > 0; i++,gunLen--) {
                            obj.guns.push(view.getUint8(i));
                        }

                        obj.hitTime = view.getUint8(i);
                        i++;

                        obj.isDead = view.getUint8(i);
                        i++;

                        let len = view.getUint8(i);
                        i++;

                        let u = new Uint8Array(msg.data.slice(i,i+len))
                        
                        obj.team = String.fromCharCode.apply(null, u);
                        i+=len;

                        len = view.getUint8(i);
                        i++;

                        u = new Uint8Array(msg.data.slice(i,i+len));
                        obj.type = new TextDecoder().decode(u);
                        i += len;
                        len = view.getUint8(i);
                        i++;

                        u = new Uint8Array(msg.data.slice(i,i+len));
                        obj.name = new TextDecoder().decode(u);
                        i += len;

                        this.objectList.forEach((obi) => {
                            if (obi.id === obj.id){
                                obi.ObjSet(obj);
                                isObjEnable = true;
                                imObjList.push(obi);
                            }
                        });
                        if (!isObjEnable && !obj.isDead) {
                            let obi = new Obj(obj.id);
                            if (obj.id === this.playerSetting.id) {
                                this.playerSetting.obj = obi;
                            }
                            obi.ObjSet(obj);
                            imObjList.push(obi);
                        }
                    }

                    this.objectList = imObjList;

                    return;
                }
                case 1: { // get Area list
                    for (let i = 1, j = 0; i < msg.data.byteLength; i+=4) {
                        this.area[j++] = view.getInt32(i);
                    }
                    return;
                }
                case 2: { // correct Connecting and start Naming
                    this.gameSetting.isNaming = true;
                    window['setTyping'](true);
                    this.textinput.value = localStorage['name'] || '';
                    this.textinputcontainer.style.display = "block";
                    this.gameSetting.isConnecting = false;
                    return;
                }
                case 3: { // scoreboard
                    
                    return;
                }
                case 4: { // ping
                    console.log(Date.now() - this.pingList[0]);
                    this.pingList.shift();
                    return;
                }
                default: {
                    return;
                }
            }
        };
        
        socket.onclose = event => { // reconnecting
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
        //console.log(Math.floor(1000 / tick)); // print fps value
        this.lastTime = Date.now();

        if (this.cv.width <= this.cv.height/9*16) { //
            this.camera.uiz = this.cv.height / 900;
        } else {
            this.camera.uiz = this.cv.width / 1600;
        }

        this.ctx.clearRect(0,0,this.cv.width,this.cv.height); // clear canvas

        if (this.gameSetting.isConnecting) { // connecting text
            this.textinputanime = 1;
            this.connectinga = Math.min(this.connectinga + 0.01 * tick,1);
        }
        
        if (this.gameSetting.isGaming) {
            console.log(data.colorList);

            drawBackground(this.ctx, this.camera.x, this.camera.y, this.camera.z, this.cv.width, this.cv.height, this.area);

            let buffer = new ArrayBuffer(14); // send input data
            let view = new DataView(buffer);
            view.setUint8(0, 1);
            if (this.input.moveVector.x === 0 && this.input.moveVector.y === 0) {
                view.setFloat32(1, 9.);
            } else {
                view.setFloat32(1, Math.atan2(this.input.moveVector.y,this.input.moveVector.x)); // move rotate
            }

            let x, y;
            if (this.input.c && this.playerSetting.obj) {
                x = Math.cos(this.playerSetting.obj.dir + 0.001 * tick) * 200 + this.playerSetting.obj.x;
                y = Math.sin(this.playerSetting.obj.dir + 0.001 * tick) * 200 + this.playerSetting.obj.y;
            } else {
                x = this.input.mousePos.x/this.camera.z + this.camera.x;
                y = this.input.mousePos.y/this.camera.z + this.camera.y;
            }
            if (this.playerSetting.obj) {
                this.playerSetting.obj.dir = Math.atan2(y - this.playerSetting.obj.y, x - this.playerSetting.obj.x);
            }

            view.setFloat32(5, x); // mouse pos
            view.setFloat32(9, y);
            view.setUint8(13,
                ((this.keys[1] || this.keys[32] || this.input.e) || 0) // is Shot?
                + ((this.keys[3] || this.keys[16]) || 0) * 2 // is Right Click?
                + (this.keys[79] || 0) * 4 // is Suicide? (o key)
            );

            socket.send(buffer);

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

            this.gameui.draw(this.ctx, this.cv.width, this.cv.height, this.camera.uiz);

        } else {
            this.ctx.drawImage(this.img, // draw Background image
                this.cv.width / 2 - this.img.width * this.camera.uiz / 2 / 2.4,
                this.cv.height / 2 - this.img.height * this.camera.uiz / 2 / 2.4,
                this.img.width * this.camera.uiz / 2.4,
                this.img.height * this.camera.uiz / 2.4);
        }

        if (this.gameSetting.isNaming || this.gameSetting.isConnecting) { // draw Black Alpha Panel
            this.panela = Math.min(this.panela + 0.05, 0.4); 

            this.ctx.save();
            this.ctx.globalAlpha = this.panela;
            this.ctx.fillStyle = "#000000";
            this.ctx.fillRect(0,0,this.cv.width,this.cv.height);
            this.ctx.restore();
        } else {
            this.panela = 0;
        }

        if (!this.gameSetting.isGaming) { // Connecting...
            drawText(this.ctx,this.cv.width / 2 / this.camera.uiz, this.cv.height / 2 / this.camera.uiz, this.camera.uiz, this.connectinga, new RGB("#FFFFFF"), "Connecting...", 60);
        }

        if (this.gameSetting.isNaming) { // draw Name input
            if (this.textinput.value) this.textinput.value = calByte.cutByteLength(this.textinput.value,15);

            let x = Math.floor(this.cv.width / 2 - 167 * this.camera.uiz),
            y = Math.floor((this.cv.height / 2 - 20 * this.camera.uiz) * (1-this.textinputanime)),
            w = Math.floor(333 * this.camera.uiz),
            h = Math.floor(41 * this.camera.uiz);

            this.textinputanime *= 0.95;
            this.connectinga = Math.max(this.connectinga - 0.2,0);

            this.textinputcontainer.style.left = window['unscale'](x) + "px";
            this.textinputcontainer.style.top = window['unscale'](y) + "px";

            this.textinput.style.width = window['unscale'](w) + "px";
            this.textinput.style.height = window['unscale'](h) + "px";
            this.textinput.style.fontSize = this.textinput.style.lineHeight = window['unscale'](h - 0.4) + "px";

            drawText(this.ctx, (x + w / 2) / this.camera.uiz, y / this.camera.uiz - 11, this.camera.uiz, 1, new RGB("#FFFFFF"), "This is the tale of...", 20);
            drawText(this.ctx, (x + w / 2) / this.camera.uiz, y / this.camera.uiz + 57, this.camera.uiz, 1, new RGB("#FFFFFF"), "(press enter to spawn)", 11.7);

            this.ctx.save();

            this.ctx.beginPath();

            this.ctx.lineCap = "round";
            this.ctx.lineJoin = "round";
            this.ctx.lineWidth = 4.5 * this.camera.uiz;
            this.ctx.fillStyle = "#FFFFFF";
            this.ctx.strokeStyle = "#000000";

            this.ctx.moveTo(x, y + 1);
            this.ctx.lineTo(x + w, y + 1);
            this.ctx.lineTo(x + w, y + h + 1);
            this.ctx.lineTo(x, y + h + 1);
            this.ctx.lineTo(x, y + 1);

            this.ctx.fill();
            this.ctx.stroke();

            this.ctx.closePath();

            this.ctx.restore();
        } else {
            this.textinputanime = 1;
        }

        requestAnimationFrame(this.loop.bind(this));
    }
}
