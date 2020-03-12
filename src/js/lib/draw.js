export const drawCircle = function (ctx, x, y, z, r){
    ctx.beginPath();
    ctx.arc(x * z, y * z,r * z, 0, Math.PI * 2);
    ctx.fill();
    ctx.closePath();
}

export const drawPolygon = function (ctx, x, y, z, r, dir, p){
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

export const drawC = function (ctx, color, color2){
    ctx.fillStyle = color.getRGB();
    ctx.strokeStyle = (color2)?color2.getRGB():color.getRGB();
}

export const drawObj = function (ctx, x, y, z, r, dir, p, o, c){
    ctx.save();
    ctx.globalAlpha = o;
    if (p < 3) {
        drawC(ctx,c.getDarkRGB());
        drawCircle(ctx, x, y, z, (r + 4));
        drawC(ctx,c);
        drawCircle(ctx, x, y, z, r);
    } else {
        drawC(ctx,c,c.getDarkRGB());
        drawPolygon(ctx, x, y, z, r, dir, p);
    }
    ctx.restore();
}