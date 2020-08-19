import { RGB, getPolygonRadius, getObjectPoint, getTextWidth, colorType, gunList } from '../lib/util';
import { colorList } from './console';
import { drawObj, drawText, drawBar } from '../lib/draw';

function deepClone(obj) {
    if(obj === null || typeof obj !== 'object') {
        return obj;
    }
    
    const result = Array.isArray(obj) ? [] : {};
    
    for(let key of Object.keys(obj)) {
        result[key] = deepClone(obj[key])
    }
    
    return result;
}

export const Obj = function(id) {
    'use strict';

    /*
        Most of these objects are contained in the objectList in the system.
        Most of the things that conflict detection is applied to within the game are included in this object.
    */

    this.id = id;

    this.name;
    this.type;
    this.guns = [];
    this.color;

    this.x;
    this.y;
    this.r; // radius
    this.dir;
    this.h; // health
    this.mh; // max health
    this.opacity;
    this.score;
    this.subobjs = [];
    this.isDead;

    this.cv = document.createElement("canvas");
    this.ctx = this.cv.getContext("2d");

    /*
        When an object has a gun, it is treated as an image if it is not opaque.
    */

    this.hitTime = 0;

    this.Animate = function (tick) {
        this.guns.forEach((g) => g.Animate(tick)); // gun animation (when their shot bullet)
        this.subobjs.forEach((o) => {
            if (o != null) {
                o.Animate(tick)
            }
        }); // sub objects animation
    }

    this.ObjSet = function (data) { // obj value setting
        this.x = data.x;
        this.y = data.y;
        if (system.playerSetting.id !== this.id) {
            /*
                If it's a self-controlled tank, the client sets the directional value.
                Otherwise, it will make the directional value smoother.
                I think the client handles the location of the direct-controlled tank a little bit.
                That part has not been implemented yet.
            */
            if (this.dir) {
                let ccw = Math.cos(data.dir)*Math.sin(this.dir)-Math.sin(data.dir)*Math.cos(this.dir);
                let a = -((Math.cos(data.dir)*Math.cos(this.dir)) + (Math.sin(data.dir)*Math.sin(this.dir))-1) * Math.PI / 2;

                if (ccw > 0) {
                    this.dir -= Math.min(a, 0.3);
                } else if (ccw < 0) {
                    this.dir += Math.min(a, 0.3);
                }
            } else {
                this.dir = data.dir;
            }
        }
        this.r = data.r;
        this.h = data.h;
        this.mh = data.mh;
        this.opacity = data.opacity;
        this.score = data.score;
        this.isDead = data.isDead;
        this.hitTime = data.hitTime;

        this.name = data.name;
        if (this.type !== data.type) {
            this.guns = [];
            if (gunList[data.type] != undefined) {
                this.guns = deepClone(gunList[data.type]); // clone gun object
            }
        }
        for (let i = 0; i < this.guns.length && i < data.guns.length; i++) {
            if (data.guns[i]) {
                this.guns[i].Shot();
            }
        }

        if (this.subobjs.length > 0) {
            for (let i = 0; i < this.subobjs.length && i < data.subobjs.length; i++) {
                if (data.subobjs[i] == null) continue;
                data.subobjs[i].isDead = this.isDead;
                this.subobjs[i].ObjSet(data.subobjs[i]);
            }
        } else {
            data.subobjs.forEach((d) => {
                if (d == null) {
                    this.subobjs.push(null);
                } else {
                    let obi = new Obj(d.id);
                    obi.ObjSet(d);
                    this.subobjs.push(obi);
                }
            });
        }

        this.type = data.type;
        this.color = colorType(data.type,data.team);
    };

    this.DrawSet = function (camera) {
        let c = colorList[this.color]; // get color value
        if (this.hitTime > 60) { // hit effect
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
        var {x, y, z, t, r, dir} = this.DrawSet(camera);
        let rr = r * getPolygonRadius(Math.abs(getObjectPoint(t)));
        var size = {x: rr * z * 2, y: rr * z * 2,};
        var pos = {x: rr * z, y: rr * z,};
        this.guns.forEach((g) => g.SetCanvasSize(camera, size, pos, rr, dir));
        let setSubobjCanvas = function (objs) {
            objs.forEach((o) => {
                if (o != null) {
                    let {px, py, w, h} = o.SetCanvasSize(camera);
                    let ax = -(x + camera.x - o.x) * z + pos.x + 2 * camera.z + 2;
                    let ay = -(y + camera.y - o.y) * z + pos.y + 2 * camera.z + 2;
                    let x1 = Math.floor(ax - px);
                    let x2 = Math.floor(ax - px + w);
                    let y1 = Math.floor(ay - py);
                    let y2 = Math.floor(ay - py + h);
        
                    if (x1 < 0) {
                        size.x -= x1;
                        pos.x -= x1;
                    } else if (x1 > size.x) {
                        size.x = x1;
                    }
        
                    if (y1 < 0) {
                        size.y -= y1;
                        pos.y -= y1;
                    } else if (y1 > size.y) {
                        size.y = y1;
                    }
        
                    if (x2 < 0) {
                        size.x -= x2;
                        pos.x -= x2;
                    } else if (x2 > size.x) {
                        size.x = x2;
                    }
        
                    if (y2 < 0) {
                        size.y -= y2;
                        pos.y -= y2;
                    } else if (y2 > size.y) {
                        size.y = y2;
                    }

                    setSubobjCanvas(o.subobjs);
                }
            });
        };
        setSubobjCanvas(this.subobjs);
        this.cv.width = size.x + 4 * camera.z + 4;
        this.cv.height = size.y + 4 * camera.z + 4;
        pos.x += 2 * camera.z + 2;
        pos.y += 2 * camera.z + 2;
        this.ctx.lineWidth = 2 * camera.z;
        this.ctx.imageSmoothingEnabled = false;
        return {
            ctxx: this.ctx,
            px: pos.x,
            py: pos.y,
            w: size.x,
            h: size.y,
        }
    }

    this.Draw = function (ctx, camera) {
        if ((this.guns.length > 0 || this.subobjs.length > 0) && this.opacity < 1){
            var {ctxx, px, py} = this.SetCanvasSize(camera);

            let {x, y, z, o} = this.DrawSet(camera);
            let lx = (x * z - px) - Math.floor(x * z - px) - x;
            let ly = (y * z - py) - Math.floor(y * z - py) - y;

            let drawObjb = function (obj) {
                let {x, y, z, t, c, r, dir} = obj.DrawSet(camera);
                let i = 0;
                for (; i < obj.subobjs.length; i++) {
                    if (obj.subobjs[i] === null) {
                        break;
                    }
                    drawObjb(obj.subobjs[i]);
                }
                obj.guns.forEach((g) => {
                    g.Draw(ctxx, camera, px / z + lx + x, py / z + ly + y, r, c, dir, obj.hitTime);
                });
                drawObj(ctxx, px / z + lx + x, py / z + ly + y, z, r, dir, t, 1, c);
                for (i++; i < obj.subobjs.length; i++) {
                    drawObjb(obj.subobjs[i]);
                }
            }
            drawObjb(this);

            ctx.save();
            ctx.globalAlpha = o;
            ctx.imageSmoothingEnabled = false;
            ctx.drawImage(this.cv, Math.floor(x * z - px), Math.floor(y * z - py));
            ctx.restore();
        } else if (this.opacity > 0) {
            let drawObjb = function(obj) {
                var {x, y, z, t, c, r, dir, o} = obj.DrawSet(camera);
                let i = 0;
                for (; i < obj.subobjs.length; i++) {
                    if (obj.subobjs[i] === null) {
                        break;
                    }
                    drawObjb(obj.subobjs[i]);
                }
                obj.guns.forEach((g) => {
                    g.Draw(ctx, camera, x, y, r, c, dir, obj.hitTime);
                });
                drawObj(ctx, x, y, z, r, dir, t, o, c);
                for (i++; i < obj.subobjs.length; i++) {
                    drawObjb(obj.subobjs[i]);
                }
            }
            drawObjb(this);
        }
    }

    this.DrawName = function (ctx, camera) {
        /*
            A function that draws the name and score of an object.
            This function is not perfect yet.
            Three things are incomplete in total.
            One is layer,
            One is the way to draw.
            The other is that the name changes to yellow when you use the cheats in sandbox mode.
        */

        ctx.save();

        const {x, y, z, r, o} = this.DrawSet(camera);
        
        if (this.score) {
            if (this.name) {
                this.cv.width = getTextWidth(this.name, "bold " + 0.8 * r * z + "px Ubuntu") + 5 * z;
                this.cv.height = 5 * 0.8 * r * z;
                this.ctx.imageSmoothingEnabled = false;
                drawText(this.ctx, this.cv.width / 2 / z, this.cv.height / 2 / z, z, 1, new RGB("#FFFFFF"), this.name, 0.8 * r);
                ctx.globalAlpha = o;
                ctx.drawImage(this.cv, x * z - this.cv.width / 2, (y - r * 1.8) * z - this.cv.height / 2);
            }

            this.cv.width = getTextWidth(this.score, "bold " + 0.6 * r * z + "px Ubuntu") + 5 * z;
            this.cv.height = 5 * 0.4 * r * z;
            this.ctx.imageSmoothingEnabled = false;
            drawText(this.ctx, this.cv.width / 2 / z, this.cv.height / 2 / z, z, 1, new RGB("#FFFFFF"), this.score, 0.6 * r);
            ctx.globalAlpha = o;
            ctx.drawImage(this.cv, x * z - this.cv.width / 2, (y - r * 1.2) * z - this.cv.height / 2);
        } else {
            if (this.name) {
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
        /*
            Function that draws the object's physical strength bar.
            This function is not perfect yet.
            Three things are incomplete in total.
            One is layer,
            One is the animation of the fitness bar.
            The other is the length of the fitness bar, depending on the size or type of object.
        */

        let healthPercent = this.h/this.mh;

        this.hpBarP -= (this.hpBarP - healthPercent) / 4;
    
        if (healthPercent < 1) {
            this.hpBarO = Math.min(this.hpBarO+0.4,1);
        } else {
            this.hpBarO = Math.max(this.hpBarO-0.2,0);
        }

        if (this.hpBarO > 0) {
            const {x, y, z, r, o} = this.DrawSet(camera);

            drawBar(ctx, x - r, y + r * 5 / 3, 4.1, r * 2, z, o * this.hpBarO, this.hpBarP, new RGB("#86e27f"));
        }
    }
}