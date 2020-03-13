import { colorList } from '../data/index';
import { drawC } from '../lib/draw';

export const Gun = function (paths, dir, color, isFront) {
    'use strict';

    this.paths = paths;
    this.dir = dir;
    this.color = color;
    this.isFront = isFront || false;
    this.back = 0;

    this.hitTime = 0;

    this.Animate = function (tick) {
        this.hitTime = Math.max(this.hitTime - tick,0);
    }

    this.SetCanvasSize = function (camera, size, pos, r, dir) {
        this.point.forEach((p) => {
            let x = (Math.cos(dir - Math.PI/2 + this.dir) * p[0] * camera.z * r
            + Math.cos(dir + this.dir) * (p[1] * camera.z * r + this.back) + pos.x);
            let y = (Math.sin(dir - Math.PI/2 + this.dir) * p[0] * camera.z * r
            + Math.sin(dir + this.dir) * (p[1] * camera.z * r + this.back) + pos.y);
            
            if (x < 0) {
                size.x -= x;
                pos.x -= x;
            } else if (x > size.x) {
                size.x = x;
            }

            if (y < 0) {
                size.y -= y;
                pos.y -= y;
            } else if (y > size.y) {
                size.y = y;
            }
        });
    }

    this.Hit = function () {
        this.hitTime = 100;
    }

    this.DrawSet = function (camera, x, y) { 
        let c = colorList[this.color];
        if (this.hitTime > 60) {
            c.getLightRGB(1 - (this.hitTime - 60) / 40);
        } else if (this.hitTime > 0) {
            c.getRedRGB(1 - this.hitTime / 60);
        }
        return {
            x: x - camera.x,
            y: y - camera.y,
            z: camera.z,
            ddir: this.dir,
            b: this.back,
            c: c,
        }
    }

    this.Draw = function (ctx, camera, x, y, r, dir) {
        const {x, y, z, ddir, b, c} = this.DrawSet(camera, x, y);

        ctx.save();
        drawC(ctx,c,c.getDarkRGB());
        let first = true;
        ctx.beginPath();
        this.point.forEach((p) => {
            let xx = Math.cos(dir - Math.PI/2 + ddir) * p[0] * z * r + Math.cos(dir + ddir) * (p[1] * z * r + b) + x;
            let yy = Math.sin(dir - Math.PI/2 + ddir) * p[0] * z * r + Math.sin(dir + ddir) * (p[1] * z * r + b) + y;
            first?ctx.moveTo(xx,yy):ctx.lineTo(xx,yy);
            first = false;
        });
        ctx.fill();
        ctx.stroke();
        ctx.closePath();
        ctx.restore();
    }

    this.Shot = function (){
        
    }
}