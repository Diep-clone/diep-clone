import { RGB } from '../lib/util';
import { colorList, colorType, gunList } from '../data/index';
import { drawC, drawObj } from '../lib/draw';
import { Socket } from '../system';

export const Obj = function(id) {
    'use strict';

    this.id = id;

    this.name;
    this.type;
    this.guns = [];
    this.color;

    this.x;
    this.y;
    this.r;
    this.dir;
    this.h;
    this.mh;
    this.opacity;
    this.score;
    this.isDead;

    this.cv = document.createElement("canvas");
    this.ctx = this.cv.getContext("2d");

    this.hitTime = 0;

    this.Animate = function (tick) {
        if (this.isDead) {
            this.opacity = Math.max(this.opacity - 0.13 * tick * 0.05, 0);
            this.r += 0.4 * tick * 0.05;

            if (this.opacity == 0) {
                system.RemoveObject(this.id);
                return;
            }
        }

        this.guns.forEach((g) => g.Animate());

        this.hitTime = Math.max(this.hitTime - tick,0);
    }

    this.ObjSet = function (data) {
        this.x = data.x;
        this.y = data.y;
        if (!this.isDead){
            this.r = data.r;
            this.dir = data.dir;
            this.h = data.h;
            this.mh = data.mh;
            this.opacity = data.opacity;
            this.score = data.score;
        }
        this.isDead = data.isDead;

        this.name = data.name;
        if (this.type !== data.type){
            if (gunList[data.type] == undefined){
                this.guns = [];
            } else {
                this.guns = gunList[data.type];
            } 
        }
        this.type = data.type;
        this.color = colorType(data.type,data.team);
    };

    this.DrawSet = function (camera) {
        let c = colorList[this.color];
        if (this.hitTime > 60) {
            c = c.getLightRGB(1 - (this.hitTime - 60) / 40);
        } else if (this.hitTime > 0) {
            c = c.getRedRGB(1 - this.hitTime / 60);
        }
        return {
            x: this.x - camera.x,
            y: this.y - camera.y,
            z: camera.z,
            t: this.type,
            c: c,
            r: this.r,
            dir: this.dir,
            o: this.opacity,
        };
    }

    this.SetCanvasSize = function (camera) {
        var {z, t, c, r, dir, o} = this.DrawSet(camera);
        var size = {x: r * z, y: r * z,};
        var pos = {x: r * z / 2, y: r * z / 2,};
        this.guns.forEach((g) => g.SetCanvasSize(camera, size, pos, r, dir));
        this.cv.width = size.x + 4 * camera.z;
        this.cv.height = size.y + 4 * camera.z;
        pos.x += 2 * camera.z;
        pos.y += 2 * camera.z;
        this.ctx.lineWidth = 2 * camera.z;
        this.ctx.imageSmoothingEnabled = false;
        return {
            ctx: this.ctx,
            x: pos.x,
            y: pos.y,
            z: z,
            t: t,
            c: c,
            r: r,
            dir: dir,
            o: o,
        }
    }

    this.Draw = function (ctx,camera) {
        if (this.guns.length > 0 && this.opacity < 1){
            var {ctxx, x, y, z, t, c, r, dir, o} = this.SetCanvasSize(camera);
            this.guns.forEach((g) => {
                if (!g.isFront) {
                    g.Draw(ctxx, camera, x, y, r, dir);
                }
            });
            drawObj(ctxx,
                x + s.x * z - x - Math.floor(s.x * z - x),
                y + s.y * z - y - Math.floor(s.y * z - y),
            z, r, dir, t, 1, c);
            this.guns.forEach((g) => {
                if (g.isFront) {
                    g.Draw(ctxx, camera, x, y, r, dir);
                }
            });
            var s = this.DrawSet(camera);
            ctx.drawImage(this.cv,Math.floor(s.x * z - x),Math.floor(s.y * z - y));
        } else {
            var {x, y, z, t, c, r, dir, o} = this.DrawSet(camera);
            this.guns.forEach((g) => {
                if (!g.isFront) {
                    g.Draw(ctx, camera, x, y, r, dir);
                }
            });
            drawObj(ctx, x, y, z, r, dir, t, o, c);
            this.guns.forEach((g) => {
                if (g.isFront) {
                    g.Draw(ctx, camera, x, y, r, dir);
                }
            });
        }
    }

    this.DrawName = function (ctx, camera) {
        /*ctx.save();

        const {x,y,z,r,o} = this.DrawSet(camera);

        ctx.font = "bold " + 0.8 * r * z + "px Ubuntu";
        ctx.lineWidth = 2.5 * z;
        ctx.textAlign = "center";
        ctx.textBaseLine = "bottom";
        ctx.globalAlpha = o;
        ctx.strokeStyle = "#000000";
        ctx.fillStyle = "#ffffff";
        if (this.score <= 0){
            ctx.strokeText(this.name,x * z,(y - r - 5) * z);
            ctx.fillText(this.name,x * z,(y - r - 5) * z);
        }
        else{
            ctx.strokeText(this.name,x * z,(y - r - 15) * z);
            ctx.fillText(this.name,x * z,(y - r - 15) * z);
            ctx.font = "bold " + 0.4 * r * z + "px Ubuntu";
            ctx.strokeText(this.score,x * z,(y - r - 5) * z);
            ctx.fillText(this.score,x * z,(y - r - 5) * z);
        }
        ctx.restore();*/
    }

    this.hpBarP = 1; // hp bar Percent
    this.hpBarO = 0; // hp bar Opacity

    this.DrawHPBar = function(ctx, camera) {
        let healthPercent = this.h/this.mh;

        this.hpBarP -= (this.hpBarP - healthPercent) / 10;
    
        if (healthPercent < 1) {
            this.hpBarO = Math.min(this.hpBarO+0.1,1);
        } else {
            this.hpBarO = Math.max(this.hpBarO-0.1,0);
        }

        //console.log(healthPercent,this.hpBarP,this.hpBarO);

        if (this.hpBarO > 0) {
            const {x, y, z, r, o} = this.DrawSet(camera);

            ctx.save();
            ctx.globalAlpha = o * this.hpBarO;
            ctx.lineCap = "round";
            ctx.lineJoin = "round";
            ctx.lineWidth = 4.1 * z;

            ctx.beginPath();
            ctx.moveTo((x + r) * z, (y + r * 5 / 3) * z);
            ctx.lineTo((x - r) * z, (y + r * 5 / 3) * z);
            ctx.closePath();
            drawC(ctx, new RGB("#444444"));
            ctx.stroke();
        
            if (this.hpBarP > 0){
                ctx.lineWidth = 2.6 * z;
                ctx.beginPath();
                ctx.moveTo((x - r) * z, (y + r * 5 / 3) * z);
                ctx.lineTo((x - r + this.hpBarP * r * 2) * z, (y + r * 5 / 3) * z);
                ctx.closePath();
                drawC(ctx, new RGB("#86e27f"));
                ctx.stroke();
            }

            ctx.restore();
        }
    }
}