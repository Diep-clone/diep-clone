import { RGB } from '../lib/util';
import { colorList } from '../data/index';
import { drawC, drawObj } from '../lib/draw';

export const Obj = function(id, name, type, color, guns) {
    'use strict';

    this.id = id;

    this.name = name;
    this.guns = guns;

    this.type = type;
    this.color = color;

    this.x;
    this.y;
    this.r;
    this.dir;
    this.h;
    this.mh;
    this.a;
    this.score;
    this.isDead;

    this.cv = document.createElement("canvas");
    this.ctx = cv.getContext("2d");

    this.r = 0;
    this.w = 0;

    this.Animate = function (tick) {
        if (this.isDead){
            this.opacity = Math.max(this.opacity - 0.13 * tick * 0.05, 0);
            this.radius += 0.4 * tick * 0.05;
            if (this.opacity == 0){
                //system.removeObject(this.id);
                return;
            }
        }
    }

    this.ObjChange = function (name, type, color, guns) {
        this.name = name;
        this.guns = guns;
    
        this.type = type;
        this.color = color;
    }
    
    this.ObjSet = function (x, y, r, dir, h, mh, a, score, isDead) {
        this.x = x;
        this.y = y;
        this.r = r;
        this.dir = dir;
        this.h = h;
        this.mh = mh;
        this.a = a;
        this.score = score;
        this.isDead = isDead;
    }

    this.DrawSet = function (camera) {
        return {
            x: this.x - camera.x,
            y: this.y - camera.y,
            z: camera.z,
            t: this.type,
            c: colorList[this.color],
            r: this.radius,
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
        this.ctx.lineCap = "round";
        this.ctx.lineJoin = "round";
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
        if (this.guns.length>0){
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
                    g.Draw(ctxx, camera, x, y, r, dir);
                }
            });
            drawObj(ctx, x, y, z, r, dir, t, o, c);
            this.guns.forEach((g) => {
                if (g.isFront) {
                    g.Draw(ctxx, camera, x, y, r, dir);
                }
            });
        }
    }

    this.DrawStatus = function (ctx, camera) {
        if (this.name) this.drawName(ctx, camera);
        this.drawHPBar(ctx, camera);
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

    this.hpBarP = 1;
    this.hpBarO = 0;

    this.DrawHPBar = function(ctx, camera) {
        let healthPercent = this.health/this.maxHealth;

        this.hpBarP -= (this.hpBarP - healthPercent) / 3;
    
        if (healthPercent < 1){
            this.hpBarO = Math.max(this.hpBarO-0.1,0);
        }else{
            this.hpBarO = Math.min(this.hpBarO+0.1,1);
        }

        if (this.hpBarO > 0){
            ctx.save();
            ctx.globalAlpha = this.opacity * this.hpBarO;

            const {x, y, z, r} = this.DrawSet(camera);

            ctx.beginPath();
            ctx.moveTo((x + r) * z, (y + r * 5 / 3) * z);
            ctx.lineTo((x - r) * z, (y + r * 5 / 3) * z);
            ctx.closePath();
            drawC(ctx,"#444444");
            ctx.lineWidth = 4.1 * z;
            ctx.stroke();
        
            ctx.beginPath();
            ctx.moveTo((x - r) * z, (y + r * 5 / 3) * z);
            ctx.lineTo((x - r + this.hpBarP * r * 2) * z, (y + r * 5 / 3) * z);
            ctx.closePath();
            drawC(ctx,"#86e27f");
            ctx.lineWidth = 2.6 * z;
            ctx.stroke();

            ctx.restore();
        }
    }
}