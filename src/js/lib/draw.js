import { backgroundColor } from '../data/index'

export const drawCircle = function (ctx, x, y, z, r) {
    ctx.beginPath();
    ctx.arc(x * z, y * z,r * z, 0, Math.PI * 2);
    ctx.fill();
    ctx.closePath();
}

export const drawPolygon = function (ctx, x, y, z, r, dir, p) {
    ctx.beginPath();
    let t = (dir % 2)? 0: Math.PI/p;
    ctx.moveTo((x + Math.cos(t) * r) * z,(y + Math.sin(t) * r) * z);
    for (;t >= -Math.PI * 2;t += Math.PI*2 /p){
        ctx.lineTo((x + Math.cos(t) * r) * z,(y + Math.sin(t) * r) * z);
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
    drawC(ctx,c,c.getDarkRGB());
    switch (t) {
        case "Triangle":
            drawPolygon(ctx, x, y, z, r, dir, 3);
            break;
        case "Shape":
            drawPolygon(ctx, x, y, z, r, dir, 4);
            break;
        case "Pentagon":
            drawPolygon(ctx, x, y, z, r, dir, 5);
            break;
        case "Mothership":
            drawPolygon(ctx, x, y, z, r, dir, 20);
            break;
        default:
            drawC(ctx,c.getDarkRGB());
            drawCircle(ctx, x, y, z, (r + 4));
            drawC(ctx,c);
            drawCircle(ctx, x, y, z, r);
            break;
    }
    ctx.restore();
}

export const drawBackground = function (ctx, x, y, z, w, h, area) {
    drawC(ctx, backgroundColor.getDarkRGB());
    ctx.fillRect(0,0,w,h);
    drawC(ctx, backgroundColor);
    ctx.fillRect((-area[0].x - x) * z, (-area[0].y - y) * z, area[0].x * z * 2, area[1].y * z * 2);

    for (let i=1;i<area.length;i++){
        drawC(ctx, area[i].c);
        ctx.fillRect((-area[i].x - x) * z, (-area[i].y - y) * z, area[i].x * z * 2, area[i].y * z * 2);  
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
}