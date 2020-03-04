function Gun(paths,dir){
  "use strict";

  this.point = paths;
  this.addRotate = dir;
  this.color = new RGB(153,153,153);
  this.backGun = 0;
  this.list = [];

  this.animate = function () {
    for (let i=0;i<this.list.length;i++){
      if (this.list[i]<1){
        this.backGun+=this.list[i];
        this.list[i]+=0.2;
      }
      else this.list.splice(i,1);
    }
  }

  this.setParentCanvasSize = function (tank,camera){
    let rotate = tank.imRotate;
    let radius = tank.showRadius;
    let xx = tank.canvasPos.x;
    let yy = tank.canvasPos.y;

    for (let i=0;i<this.point.length;i++){
      let x = Math.floor(this.point[i][0]*Math.cos(rotate-Math.PI/2+this.addRotate)*camera.z*radius+this.point[i][1]*Math.cos(rotate+this.addRotate)*camera.z*radius+Math.cos(rotate+this.addRotate)*this.backGun+xx);
      let y = Math.floor(this.point[i][0]*Math.sin(rotate-Math.PI/2+this.addRotate)*camera.z*radius+this.point[i][1]*Math.sin(rotate+this.addRotate)*camera.z*radius+Math.sin(rotate+this.addRotate)*this.backGun+yy);

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
    let rotate = tank.imRotate;
    let radius = tank.showRadius;
    let xx = tank.canvasPos.x;
    let yy = tank.canvasPos.y;

    ctx.save();
    ctx.strokeStyle = this.color.getDarkRGB().getRedRGB(tank.r).getLightRGB(tank.w).getRGBValue();
    ctx.fillStyle = this.color.getRedRGB(tank.r).getLightRGB(tank.w).getRGBValue();

    ctx.beginPath();
    let x = this.point[0][0]*Math.cos(rotate-Math.PI/2+this.addRotate)*camera.z*radius+this.point[0][1]*Math.cos(rotate+this.addRotate)*camera.z*radius+Math.cos(rotate+this.addRotate)*this.backGun+xx;
    let y = this.point[0][0]*Math.sin(rotate-Math.PI/2+this.addRotate)*camera.z*radius+this.point[0][1]*Math.sin(rotate+this.addRotate)*camera.z*radius+Math.sin(rotate+this.addRotate)*this.backGun+yy;
    ctx.moveTo(x,y);
    for (let i=0;i<this.point.length;i++){
      let x = this.point[i][0]*Math.cos(rotate-Math.PI/2+this.addRotate)*camera.z*radius+this.point[i][1]*Math.cos(rotate+this.addRotate)*camera.z*radius+Math.cos(rotate+this.addRotate)*this.backGun+xx;
      let y = this.point[i][0]*Math.sin(rotate-Math.PI/2+this.addRotate)*camera.z*radius+this.point[i][1]*Math.sin(rotate+this.addRotate)*camera.z*radius+Math.sin(rotate+this.addRotate)*this.backGun+yy;
      ctx.lineTo(x,y);
    }
    ctx.fill();
    ctx.stroke();
    ctx.closePath();
    ctx.restore();
  }

  this.shot = function (){
    this.list.push(-0.9);
  }
}
