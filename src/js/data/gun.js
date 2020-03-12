import { colorList } from '../data/index';
import { drawC } from '../lib/draw';

export const Gun = function (paths, dir, color, isFront) {
    'use strict';

    this.paths = paths;
    this.dir = dir;
    this.color = colorList[color];
    this.isFront = isFront || false;
    this.back = 0;
    this.list = [];

    this.Animate = function (tick) {
        for (let i = 0;i < this.list.length;i ++){
            if (this.list[i] < 1){
                this.back += this.list[i];
                this.list[i] += 0.2;
            }
            else this.list.splice(i--,1);
        }
    }

    this.SetCanvasSize = function (camera, size, pos, r, dir) {
        let xx = pos.x;
        let yy = pos.y;

        this.point.forEach((p) => {
            let x = (p[0] * Math.cos(dir - Math.PI/2 + this.dir) * camera.z * r + p[1] * Math.cos(dir + this.dir) * camera.z * r + Math.cos(dir + this.dir) * this.back + xx);
            let y = (p[0] * Math.sin(dir - Math.PI/2 + this.dir) * camera.z * r + p[1] * Math.sin(dir + this.dir) * camera.z * r + Math.sin(dir + this.dir) * this.back + yy);
            
            if (x < 0) {
                size.x += -x;
                pos.x += -x;
                xx = pos.x;
            } else if (x > size.x) {
                size.x = x;
            }

            if (y < 0) {
                size.y += -y;
                pos.y += -y;
                yy = pos.y;
            } else if (y > size.y) {
                size.y = y;
            }
        });
    }

    this.Draw = function (ctx, camera, x, y, r, dir) {

    }

    this.Shot = function (){
        this.list.push(-0.9);
    }
}