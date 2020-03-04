function Bolt(paths,dir){
  "use strict";

  this.point = paths;
  this.addRotate = dir;
  this.rotate = 0;
  this.color = new RGB(80,80,80);

  this.animate = function () {
    this.rotate+=this.addRotate;
  }

  this.setParentCanvasSize = function (tank,camera){
    let rotate = this.rotate;
    let radius = tank.showRadius;
    let xx = tank.canvasPos.x;
    let yy = tank.canvasPos.y;

    for (let i=0;i<this.point.length;i++){
      let x = Math.floor(this.point[i][0]*Math.cos(rotate-Math.PI/2+this.addRotate)*camera.z*radius+this.point[i][1]*Math.cos(rotate+this.addRotate)*camera.z*radius+xx);
      let y = Math.floor(this.point[i][0]*Math.sin(rotate-Math.PI/2+this.addRotate)*camera.z*radius+this.point[i][1]*Math.sin(rotate+this.addRotate)*camera.z*radius+yy);

      if (x<0){
        tank.canvasSize.x += -x;
        tank.canvasPos.x += -x;
        xx = tank.canvasPos.x;
      }
      else if (x>tank.canvasSize.x){
        tank.canvasSize.x = x;
      }
      if (y<0){
        tank.canvasSize.y += -y;
        tank.canvasPos.y += -y;
        yy = tank.canvasPos.y;
      }
      else if (y>tank.canvasSize.y){
        tank.canvasSize.y = y;
      }
    }
  }

  this.drawGun = function (tank,ctx,camera){
    let rotate = this.rotate;
    let radius = tank.showRadius;
    let xx = tank.canvasPos.x;
    let yy = tank.canvasPos.y;

    ctx.save();
    ctx.strokeStyle = this.color.getDarkRGB().getRedRGB(tank.r).getLightRGB(tank.w).getRGBValue();
    ctx.fillStyle = this.color.getRedRGB(tank.r).getLightRGB(tank.w).getRGBValue();

    ctx.beginPath();
    let x = this.point[0][0]*Math.cos(rotate-Math.PI/2+this.addRotate)*camera.z*radius+this.point[0][1]*Math.cos(rotate+this.addRotate)*camera.z*radius+xx;
    let y = this.point[0][0]*Math.sin(rotate-Math.PI/2+this.addRotate)*camera.z*radius+this.point[0][1]*Math.sin(rotate+this.addRotate)*camera.z*radius+yy;
    ctx.moveTo(x,y);
    for (let i=0;i<this.point.length;i++){
      let x = this.point[i][0]*Math.cos(rotate-Math.PI/2+this.addRotate)*camera.z*radius+this.point[i][1]*Math.cos(rotate+this.addRotate)*camera.z*radius+xx;
      let y = this.point[i][0]*Math.sin(rotate-Math.PI/2+this.addRotate)*camera.z*radius+this.point[i][1]*Math.sin(rotate+this.addRotate)*camera.z*radius+yy;
      ctx.lineTo(x,y);
    }
    ctx.lineTo(x,y);
    ctx.fill();
    ctx.stroke();
    ctx.closePath();
    ctx.restore();
  }
}
