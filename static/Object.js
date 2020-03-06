function Obj(id,name,paths,r,g,b){
    'use strict';

    this.id = id;
    
    this.name = name;
    this.paths = paths;
    this.color = new RGB(r,g,b);
    this.guns = [];

    this.x;
    this.y;
    this.r;
    this.dir;
    this.h;
    this.mh;
    this.a;
    this.score;
    this.isDead;

    this.animate = function (tick){
        if (this.isDead){
            this.opacity = Math.max(this.opacity - 0.13 * tick * 0.05, 0);
            this.radius += 0.4 * tick * 0.05;
            if (this.opacity == 0){
                system.removeObject(this.id);
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

    this.GraphicSet = function (name,vertex,r,g,b,guns){
        this.name = name;
        this.vertex = vertex;
        this.color.SetRGB(r,g,b);
        //this.guns = guns;
    }

    this.draw = function (ctx,camera){

    }

    this.drawStatus = function (ctx,camera){
        if (this.name) this.drawName(ctx,camera);
        this.drawHPBar(ctx,camera);
    }

    this.drawName = function (ctx,camera){
        ctx.save();
        ctx.font = "bold " + 0.8 * this.radius * camera.z + "px Ubuntu";
        ctx.lineWidth = 2.5 * camera.z;
        ctx.lineCap = "round";
        ctx.lineJoin = "round";
        ctx.textAlign = "center";
        ctx.textBaseLine = "bottom";
        ctx.globalAlpha = this.opacity;
        ctx.strokeStyle = "#000000";
        ctx.fillStyle = "#ffffff";
        if (this.score <= 0){
            ctx.strokeText(this.name,(this.x - camera.x) * camera.z,(this.y - this.radius - 5 - camera.y) * camera.z);
            ctx.fillText(this.name,(this.x - camera.x) * camera.z,(this.y - this.radius - 5 - camera.y) * camera.z);
        }
        else{
            ctx.strokeText(this.name,(this.x - camera.x) * camera.z,(this.y - this.radius - 15 - camera.y) * camera.z);
            ctx.fillText(this.name,(this.x - camera.x) * camera.z,(this.y - this.radius - 15 - camera.y) * camera.z);
            ctx.font = "bold " + 0.4 * this.radius * camera.z + "px Ubuntu";
            ctx.strokeText(this.score,(this.x - camera.x) * camera.z,(this.y - this.radius - 5 - camera.y) * camera.z);
            ctx.fillText(this.score,(this.x - camera.x) * camera.z,(this.y - this.radius - 5 - camera.y) * camera.z);
        }
        ctx.restore();
    }

    this.showPercent = 1;

    this.drawHPBar = function(ctx,camera){
        let healthPercent = this.health/this.maxHealth;
    
        ctx.save();
        ctx.globalAlpha = this.opacity;
    
        this.showPercent -= (this.showPercent - healthPercent) / 3;
    
        if (healthPercent<1){
            ctx.beginPath();
            ctx.moveTo((this.x - camera.x + this.radius) * camera.z,(this.y - camera.y + this.radius * 5 / 3) * camera.z);
            ctx.lineTo((this.x - camera.x - this.radius) * camera.z,(this.y - camera.y + this.radius * 5 / 3) * camera.z);
            ctx.closePath();
            ctx.strokeStyle = "#444444";
            ctx.lineWidth = 4.1 * camera.z;
            ctx.stroke();
        
            ctx.beginPath();
            ctx.moveTo((this.x - camera.x - this.radius) * camera.z,(this.y - camera.y + this.radius * 5 / 3) * camera.z);
            ctx.lineTo((this.x - camera.x - this.radius + this.showPercent * this.radius * 2) * camera.z,(this.y - camera.y + this.radius * 5 / 3) * camera.z);
            ctx.closePath();
            ctx.strokeStyle = "#86e27f";
            ctx.lineWidth = 2.6 * camera.z;
            ctx.stroke();
        }
        ctx.restore();
    }
}