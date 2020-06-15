import { colorList } from './console';
import { drawC } from '../lib/draw';

export const Gun = function (paths, dir, color, isMoveDir, isFront, isStatic) {
    'use strict';

    this.paths = paths;
    this.dir = dir || 0;
    this.color = color;
    if (this.color == undefined) this.color = 1;
    this.isMoveDir = isMoveDir || 0;
    this.isFront = isFront || false;
    this.isStatic = isStatic || false;
    this.shotTime = [];
    this.back = 0;

    this.Animate = function (tick) {
        this.back = 0;
        for (let i = 0; i < this.shotTime.length;){
            this.shotTime[i] += tick;
            let s = this.shotTime[i];
            if (s < 220) {
                this.back += 3 / 24200 * s * s - 3 / 110 * s;
                i++;
            } else {
                this.shotTime.splice(i, 1);
            }
        }
        if (this.isMoveDir) {
            this.dir += this.isMoveDir * tick / 20
        }
    }

    this.SetCanvasSize = function (camera, size, pos, r, dir) {
        if (this.isMoveDir) {
            dir = 0;
        }
        const {z, ddir, b} = this.DrawSet(camera, r);

        this.paths.forEach((p) => {
            const x = Math.floor(Math.cos(dir - Math.PI/2 + ddir) * p[0] * z * ((this.isStatic)?1:r)
            + Math.cos(dir + ddir) * (p[1] * z * r + b) + pos.x);
            const y = Math.floor(Math.sin(dir - Math.PI/2 + ddir) * p[0] * z * ((this.isStatic)?1:r)
            + Math.sin(dir + ddir) * (p[1] * z * r + b) + pos.y);
            
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

    this.DrawSet = function (camera, c, hitTime) {
        if (this.color != -1) {
            c = colorList[this.color];
            if (this.hitTime > 60) {
                c.getLightRGB(1 - (hitTime - 60) / 70);
            } else if (this.hitTime > 0) {
                c.getRedRGB(1 - hitTime / 60);
            }
        }
        return {
            z: camera.z,
            ddir: this.dir,
            b: this.back,
            c: c,
        }
    }

    this.Draw = function (ctx, camera, x, y, r, cc, dir, hitTime) {
        if (isMoveDir) {
            dir = 0;
        }

        const {z, ddir, b, c} = this.DrawSet(camera, cc, hitTime);

        ctx.save();
        drawC(ctx,c,c.getDarkRGB());
        ctx.lineWidth = 2 * z;
        ctx.lineCap = "round";
        ctx.lineJoin = "round";
        let first = true;
        ctx.beginPath();
        this.paths.forEach((p) => {
            let xx = (x + Math.cos(dir - Math.PI/2 + ddir) * p[0] * ((this.isStatic)?1:r) + Math.cos(dir + ddir) * (p[1] * r + b)) * z;
            let yy = (y + Math.sin(dir - Math.PI/2 + ddir) * p[0] * ((this.isStatic)?1:r) + Math.sin(dir + ddir) * (p[1] * r + b)) * z;
            first?ctx.moveTo(xx,yy):ctx.lineTo(xx,yy);
            first = false;
        });
        ctx.fill();
        ctx.stroke();
        ctx.closePath();
        ctx.restore();
    }

    this.Shot = function (){
        this.shotTime.push(0);
    }
}