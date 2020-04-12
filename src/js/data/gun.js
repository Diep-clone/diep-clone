import { colorList } from '../data/index';
import { drawC } from '../lib/draw';

export const Gun = function (paths, dir, color, isFront) {
    'use strict';

    this.paths = paths;
    this.dir = dir;
    this.color = color;
    this.isFront = isFront || false;
    this.back = 0;

    this.Animate = function (tick) {
        
    }

    this.SetCanvasSize = function (camera, size, pos, r, dir) {
        this.paths.forEach((p) => {
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

    this.DrawSet = function (camera, hitTime) { 
        let c = colorList[this.color];
        if (this.hitTime > 60) {
            c.getLightRGB(1 - (hitTime - 60) / 40);
        } else if (this.hitTime > 0) {
            c.getRedRGB(1 - hitTime / 60);
        }
        return {
            z: camera.z,
            ddir: this.dir,
            b: this.back,
            c: c,
        }
    }

    this.Draw = function (ctx, camera, x, y, r, dir, hitTime) {
        const {z, ddir, b, c} = this.DrawSet(camera, hitTime);

        ctx.save();
        drawC(ctx,c,c.getDarkRGB());
        ctx.lineWidth = 2 * z;
        ctx.lineCap = "round";
        ctx.lineJoin = "round";
        let first = true;
        ctx.beginPath();
        this.paths.forEach((p) => {
            let xx = (x + Math.cos(dir - Math.PI/2 + ddir) * p[0] * r + Math.cos(dir + ddir) * (p[1] * r + b)) * z;
            let yy = (y + Math.sin(dir - Math.PI/2 + ddir) * p[0] * r + Math.sin(dir + ddir) * (p[1] * r + b)) * z;
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