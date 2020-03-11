export const drawCircle = function (ctx,x,y,z,r){
    ctx.beginPath();
    ctx.arc(x * z, y * z,r * z, 0, Math.PI * 2);
    ctx.fill();
    ctx.closePath();
}

export const drawPolygon = function (ctx,x,y,z,r,dir,p){
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

export const drawFC = function (ctx,color){
    ctx.fillStyle = color;
}

export const drawSC = function (ctx,color){
    ctx.strokeStyle = color;
}
