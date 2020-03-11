import { RGB } from '../lib/util';
import { drawCircle, drawPolygon, drawSC } from '../lib/draw';

export const Obj = function(id,name,type,guns){
    'use strict';

    this.id = id;
    this.name = name;
    this.type = type;
    this.color;
    this.point;

    this.guns = guns;

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

    this.Animate = function (tick){
        if (this.isDead){
            this.opacity = Math.max(this.opacity - 0.13 * tick * 0.05, 0);
            this.radius += 0.4 * tick * 0.05;
            if (this.opacity == 0){
                //system.removeObject(this.id);
                return;
            }
        }
    }
    
    this.ObjSet = function (x,y,r,dir,h,mh,a,score,isDead){
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

    this.GunSet = function (guns) {
        this.guns = guns;
    }

    this.DrawSet = function (camera) {
        return {
            x:this.x - camera.x,
            y:this.y - camera.y,
            z:camera.z,
            r:this.radius,
            dir:this.dir,
            o:this.opacity,
        };
    }

    this.SetCanvasSize = function (camera){

    }

    this.Draw = function (ctx,camera){
        const {x,y,z,r,dir,o} = this.DrawSet(camera);
        ctx.save();
        ctx.globalAlpha = o;
        const ctxx = ctx;
        if (this.guns.length>0){
            this.SetCanvasSize(camera);
            ctxx = this.ctx;
        } else {
            
        }
        

        switch (this.type){
            case "Bullet":
                drawCircle(ctxx, x, y, z, (r + 4));
                drawCircle(ctxx, x, y, z, r);
                break;
            case "Polygon":
                drawPolygon(ctxx, x, y, z, r, dir, this.point);
                break;
            default:
                break;
        }
        ctx.restore();
    }

    this.DrawStatus = function (ctx,camera){
        if (this.name) this.drawName(ctx,camera);
        this.drawHPBar(ctx,camera);
    }

    this.DrawName = function (ctx,camera){
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

    this.DrawHPBar = function(ctx,camera){
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

            const {x,y,z,r} = this.DrawSet(camera);

            ctx.beginPath();
            ctx.moveTo((x + r) * z, (y + r * 5 / 3) * z);
            ctx.lineTo((x - r) * z, (y + r * 5 / 3) * z);
            ctx.closePath();
            drawSC(ctx,"#444444");
            ctx.lineWidth = 4.1 * z;
            ctx.stroke();
        
            ctx.beginPath();
            ctx.moveTo((x - r) * z, (y + r * 5 / 3) * z);
            ctx.lineTo((x - r + this.hpBarP * r * 2) * z, (y + r * 5 / 3) * z);
            ctx.closePath();
            drawSC(ctx,"#86e27f");
            ctx.lineWidth = 2.6 * z;
            ctx.stroke();

            ctx.restore();
        }
    }
}