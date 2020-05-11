import { RGB, getPolygonRadius, getObjectPoint, getTextWidth } from '../lib/util';
import { colorList, colorType, gunList } from '../data/index';
import { drawC, drawObj, drawText } from '../lib/draw';

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

    this.isEnable = true;
    this.isDelete = false;

    this.cv = document.createElement("canvas");
    this.ctx = this.cv.getContext("2d");

    this.hitTime = 0;

    this.Animate = function (tick) {
        if (this.isDead) {
            this.opacity = Math.max(this.opacity - 0.1 * tick * 0.05, 0);
            this.r += this.r * 0.03 * tick * 0.05;

            if (this.opacity == 0) {
                this.isDelete = true;
                return;
            }
        }

        this.guns.forEach((g) => g.Animate());
    }

    this.ObjSet = function (data) {
        this.x = data.x;
        this.y = data.y;
        if (system.playerSetting.id !== this.id) {
            if (this.dir) {
                let ccw = Math.cos(data.dir)*Math.sin(this.dir)-Math.sin(data.dir)*Math.cos(this.dir);
                let a = -((Math.cos(data.dir)*Math.cos(this.dir)) + (Math.sin(data.dir)*Math.sin(this.dir))-1) * Math.PI / 2;

                if (ccw > 0) {
                    this.dir -= a / 3;
                } else if (ccw < 0) {
                    this.dir += a / 3;
                }
            } else {
                this.dir = data.dir;
            }
        }
        if (!this.isDead){
            this.r = data.r;
            this.h = data.h;
            this.mh = data.mh;
            this.opacity = data.opacity;
            this.score = data.score;
        }
        this.isDead = data.isDead;
        this.hitTime = data.hitTime;

        this.name = data.name;
        if (this.type !== data.type) {
            this.guns = (gunList[data.type] == undefined)?[]:gunList[data.type];
        }
        this.type = data.type;
        this.color = colorType(data.type,data.team);

        this.isEnable = true;
    };

    this.DrawSet = function (camera) {
        let c = colorList[this.color];
        if (this.hitTime > 60) {
            c = c.getLightRGB((this.hitTime - 60) / 70);
        } else if (this.hitTime > 0) {
            c = c.getRedRGB(this.hitTime / 60);
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
        let rr = r * getPolygonRadius(getObjectPoint(t));
        var size = {x: rr * z * 2, y: rr * z * 2,};
        var pos = {x: rr * z, y: rr * z,};
        this.guns.forEach((g) => g.SetCanvasSize(camera, size, pos, rr, dir));
        this.cv.width = size.x + 4 * camera.z;
        this.cv.height = size.y + 4 * camera.z;
        pos.x += 2 * camera.z;
        pos.y += 2 * camera.z;
        this.ctx.lineWidth = 2 * camera.z;
        this.ctx.imageSmoothingEnabled = false;
        return {
            ctxx: this.ctx,
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
            var s = this.DrawSet(camera);
            this.guns.forEach((g) => {
                if (!g.isFront) {
                    g.Draw(ctxx, camera, x / z, y / z, r, dir, this.hitTime);
                }
            });
            drawObj(ctxx,
                x / z + (s.x * z - x) - Math.floor(s.x * z - x),
                y / z + (s.y * z - y) - Math.floor(s.y * z - y),
            z, r, dir, t, 1, c);
            this.guns.forEach((g) => {
                if (g.isFront) {
                    g.Draw(ctxx, camera, x / z, y / z, r, dir, this.hitTime);
                }
            });
            ctx.save();
            ctx.globalAlpha = o;
            ctx.imageSmoothingEnabled = false;
            ctx.drawImage(this.cv,Math.floor(s.x * z - x),Math.floor(s.y * z - y));
            ctx.restore();
        } else {
            var {x, y, z, t, c, r, dir, o} = this.DrawSet(camera);
            this.guns.forEach((g) => {
                if (!g.isFront) {
                    g.Draw(ctx, camera, x, y, r, dir, this.hitTime);
                }
            });
            drawObj(ctx, x, y, z, r, dir, t, o, c);
            this.guns.forEach((g) => {
                if (g.isFront) {
                    g.Draw(ctx, camera, x, y, r, dir, this.hitTime);
                }
            });
        }
    }

    this.DrawName = function (ctx, camera) {
        ctx.save();

        const {x, y, z, r, o} = this.DrawSet(camera);
        
        if (this.name){
            if (this.score){
                this.cv.width = getTextWidth(this.name, "bold " + 0.8 * r * z + "px Ubuntu") + 5 * z;
                this.cv.height = 5 * 0.8 * r * z;
                this.ctx.imageSmoothingEnabled = false;
                drawText(this.ctx, this.cv.width / 2 / z, this.cv.height / 2 / z, z, 1, new RGB("#FFFFFF"), this.name, 0.8 * r);
                ctx.globalAlpha = o;
                ctx.drawImage(this.cv, x * z - this.cv.width / 2, (y - r * 1.8) * z - this.cv.height / 2);

                this.cv.width = getTextWidth(this.score, "bold " + 0.6 * r * z + "px Ubuntu") + 5 * z;
                this.cv.height = 5 * 0.4 * r * z;
                this.ctx.imageSmoothingEnabled = false;
                drawText(this.ctx, this.cv.width / 2 / z, this.cv.height / 2 / z, z, 1, new RGB("#FFFFFF"), this.score, 0.6 * r);
                ctx.globalAlpha = o;
                ctx.drawImage(this.cv, x * z - this.cv.width / 2, (y - r * 1.2) * z - this.cv.height / 2);
            } else {
                this.cv.width = getTextWidth(this.name, "bold " + 0.8 * r * z + "px Ubuntu") + 5 * z;
                this.cv.height = 5 * 0.8 * r * z;
                this.ctx.imageSmoothingEnabled = false;
                drawText(this.ctx, this.cv.width / 2 / z, this.cv.height / 2 / z, z, 1, new RGB("#FFFFFF"), this.name, 0.8 * r);
                ctx.globalAlpha = o;
                ctx.drawImage(this.cv, x * z - this.cv.width / 2, (y - r * 1.5) * z - this.cv.height / 2);
            }
        }

        ctx.restore();
    }

    this.hpBarP = 1; // hp bar Percent
    this.hpBarO = 0; // hp bar Opacity

    this.DrawHPBar = function(ctx, camera) {
        let healthPercent = this.h/this.mh;

        this.hpBarP -= (this.hpBarP - healthPercent) / 4;
    
        if (healthPercent < 1) {
            this.hpBarO = Math.min(this.hpBarO+0.4,1);
        } else {
            this.hpBarO = Math.max(this.hpBarO-0.2,0);
        }

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