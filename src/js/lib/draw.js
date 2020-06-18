import { backgroundColor, minimapBackgroundColor, minimapBorderColor } from '../data/console'
import { RGB, getPolygonRadius, getObjectPoint } from './util';

/* 
    전체적인 그래픽을 담당하고 있는 파일입니다.
*/

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

    r *= getPolygonRadius(Math.abs(p));

    ctx.beginPath();
    let t = (p % 2)? 0: Math.PI/p;
    for (;t >= -Math.PI * 2; t -= Math.PI*2/Math.abs(p)) {
        if (p < 0) { // draw Trap object
            ctx.lineTo((x + Math.cos(t+dir-Math.PI/p) * r * 0.4) * z,(y + Math.sin(t+dir-Math.PI/p) * r * 0.4) * z);
            ctx.lineTo((x + Math.cos(t+dir) * r) * z,(y + Math.sin(t+dir) * r) * z);
        } else {
            ctx.lineTo((x + Math.cos(t+dir) * r) * z,(y + Math.sin(t+dir) * r) * z);
        }
    }
    ctx.fill();
    ctx.stroke();
    ctx.closePath();
}

export const drawC = function (ctx, color, color2) { // set draw Color 1. one color 2. fill / stroke color
    ctx.fillStyle = color.getRGB();
    ctx.strokeStyle = (color2)?color2.getRGB():color.getRGB();
}

export const drawObj = function (ctx, x, y, z, r, dir, t, o, c) {
    ctx.save();

    ctx.globalAlpha = o;
    let im = getObjectPoint(t);
    if (im == 0) {
        drawC(ctx, c.getDarkRGB());
        drawCircle(ctx, x, y, z, (r + 1));
        drawC(ctx,c);
        drawCircle(ctx, x, y, z, (r - 1));
    } else {
        drawC(ctx, c, c.getDarkRGB());
        drawPolygon(ctx, x, y, z, r, dir, im);
    }
    
    ctx.restore();
}

export const drawText = function (ctx, x, y, z, o, c, text, size, dir) {
    ctx.save();
    ctx.font = "bold " + size * z + "px Ubuntu";
    ctx.lineWidth = 0.18 * z * size;
    ctx.textAlign = "center";
    ctx.textBaseLine = "bottom";
    ctx.globalAlpha = o;
    drawC(ctx, c, new RGB("#000000"));

    ctx.translate(x * z, y * z);

    if (dir) {
        ctx.rotate(dir);
    }

    ctx.strokeText(text, 0, 0);
    ctx.fillText(text, 0, 0);

    ctx.restore();
}

export const drawBar = function (ctx, x, y, r, l, z, o, p, c) {
    ctx.save();

    ctx.globalAlpha = o;
    ctx.lineCap = "round";
    ctx.lineJoin = "round";
    ctx.lineWidth = r * z;
    
    ctx.beginPath();
    ctx.moveTo(x * z, y * z);
    ctx.lineTo((x + l) * z, y * z);
    ctx.closePath();
    drawC(ctx, new RGB("#444444"));
    ctx.stroke();

    if (p) {
        ctx.lineWidth = (r - 1.5) * z;
        ctx.beginPath();
        ctx.moveTo(x * z, y * z);
        ctx.lineTo((x + l * p) * z, y * z);
        ctx.closePath();
        drawC(ctx, c);
        ctx.stroke();
    }

    ctx.restore();
}

export const drawMiniMap = function (ctx, x, y, z, w, h, area, pos, dir) { // not work
    ctx.save();

    ctx.lineCap = "round";
    ctx.lineJoin = "round";
    ctx.lineWidth = 4.5 * z;

    drawC(ctx, minimapBackgroundColor, minimapBorderColor);
    
    for (let i=4;i<area.length;){
        drawC(ctx, area[i++]);
        
    }

    

    ctx.restore();
}

export const drawBackground = function (ctx, x, y, z, w, h, area) {
    ctx.save();

    drawC(ctx, backgroundColor.getDarkRGB());
    ctx.fillRect(0,0,w,h);
    drawC(ctx, backgroundColor);
    ctx.fillRect((area[0] - x) * z, (area[1] - y) * z, area[2] * z, area[3] * z);

    for (let i=4;i<area.length;) {
        drawC(ctx, area[i++]);
        ctx.fillRect((area[i++] - x) * z, (area[i++] - y) * z, area[i++] * z, area[i++] * z);
    }

    ctx.beginPath(); // draw grid
    for (let i = -x % 12.9 * z; i <= w; i += 12.9 * z) {
        ctx.moveTo(i, 0);
        ctx.lineTo(i, h);
    }
    for (let i = -y % 12.9 * z; i <= h; i += 12.9 * z) {
        ctx.moveTo(0, i);
        ctx.lineTo(w, i);
    }
    ctx.strokeStyle = "black";
    ctx.globalAlpha = 0.1;
    ctx.lineWidth = 0.2;
    ctx.stroke();

    ctx.restore();
}