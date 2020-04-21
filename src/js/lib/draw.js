import { backgroundColor } from '../data/index'
import { RGB, getPolygonRadius, getObjectPoint, getTextWidth } from './util';

export const drawCircle = function (ctx, x, y, z, r) {
    ctx.beginPath();
    ctx.arc(x * z, y * z,r * z, 0, Math.PI * 2);
    ctx.fill();
    ctx.closePath();
}

export const drawPolygon = function (ctx, x, y, z, r, dir, p) {
    ctx.lineWidth = 2 * z;
    ctx.lineCap = "round";
    ctx.lineJoin = "round";

    r *= getPolygonRadius(p);

    ctx.beginPath();
    let t = (p % 2)? 0: Math.PI/p;
    ctx.moveTo((x + Math.cos(t+dir) * r) * z,(y + Math.sin(t+dir) * r) * z);
    for (;t >= -Math.PI * 2;t -= Math.PI*2/p){
        ctx.lineTo((x + Math.cos(t+dir) * r) * z,(y + Math.sin(t+dir) * r) * z);
    }
    ctx.fill();
    ctx.stroke();
    ctx.closePath();
}

export const drawC = function (ctx, color, color2) {
    ctx.fillStyle = color.getRGB();
    ctx.strokeStyle = (color2)?color2.getRGB():color.getRGB();
}

export const drawObj = function (ctx, x, y, z, r, dir, t, o, c) {
    ctx.save();
    ctx.globalAlpha = o;
    let im = getObjectPoint(t);
    if (im == 0) {
        drawC(ctx, c.getDarkRGB());
        drawCircle(ctx, x, y, z, (r + 2));
        drawC(ctx,c);
        drawCircle(ctx, x, y, z, r);
    } else {
        drawC(ctx, c, c.getDarkRGB());
        drawPolygon(ctx, x, y, z, r, dir, im);
    }
    ctx.restore();
}

export const drawText = function (ctx, x, y, z, o, c, text, size, dir) {
    ctx.save();
    let ctxSet = function (ctx) {
        ctx.font = "bold " + size * z + "px Ubuntu";
        ctx.lineWidth = 2.5 * z;
        ctx.textAlign = "center";
        ctx.textBaseLine = "bottom";
        ctx.globalAlpha = o;
        drawC(ctx, c, new RGB("#000000"));
    }
    
    if (o < 1) {
        let cv = document.createElement("canvas");
        let cctx = cv.getContext("2d");

        cv.width = getTextWidth(text, "bold " + size * z + "px Ubuntu");
        cv.height = 10 * z * size;

        ctxSet(cctx);
        cctx.globalAlpha = 1;

        cctx.translate(cv.width / 2, cv.height / 2);

        cctx.strokeText(text, x * z - cv.width / 2 - Math.floor(x * z - cv.width / 2), y * z - cv.height / 2 - Math.floor(y * z - cv.height / 2));
        cctx.fillText(text, x * z - cv.width / 2 - Math.floor(x * z - cv.width / 2), y * z - cv.height / 2 - Math.floor(y * z - cv.height / 2));

        if (dir) {
            ctx.rotate(dir);
        }
        ctx.globalAlpha = o;
        ctx.drawImage(cv, Math.floor(x * z - cv.width / 2), Math.floor(y * z - cv.height / 2));
    } else {
        ctxSet(ctx);
        ctx.translate(Math.floor(x * z), Math.floor(y * z))

        if (dir) {
            ctx.rotate(dir);
        }
        ctx.strokeText(text, x * z - Math.floor(x * z), y * z - Math.floor(y * z));
        ctx.fillText(text, x * z - Math.floor(x * z), y * z - Math.floor(y * z));
    }

    ctx.restore();
}

export const drawBackground = function (ctx, x, y, z, w, h, area) {
    ctx.save();
    drawC(ctx, backgroundColor.getDarkRGB());
    ctx.fillRect(0,0,w,h);
    drawC(ctx, backgroundColor);
    ctx.fillRect((area[0].X - x) * z, (area[0].Y - y) * z, area[0].W * z, area[0].H * z);

    for (let i=1;i<area.length;i++){
        drawC(ctx, area[i].c);
        ctx.fillRect((area[i].X - x) * z, (area[i].Y - y) * z, area[i].W * z, area[i].H * z);
    }

    ctx.beginPath(); // draw grid
    for (let i = -x % 12.9 * z; i <= w; i += 12.9 * z){
        ctx.moveTo(i, 0);
        ctx.lineTo(i, h);
    }
    for (let i = -y % 12.9 * z; i <= h; i += 12.9 * z){
        ctx.moveTo(0, i);
        ctx.lineTo(w, i);
    }
    ctx.strokeStyle = "black";
    ctx.globalAlpha = 0.1;
    ctx.lineWidth = 0.4;
    ctx.stroke();
    ctx.restore();
}