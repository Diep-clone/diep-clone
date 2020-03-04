function Pentagon(radius,rotate){
  "use strict";

  HealthShowObject.apply(this, arguments);
  this.radius = radius;
  this.rotate = rotate;
  this.color = new RGB(255,232,105);
  this.isDead = false;
  this.canvas = document.createElement("canvas");
  this.ctx = this.canvas.getContext("2d");
  this.canvasSize = {x:0,y:0};
  this.canvasPos = {x:0,y:0};
  this.imRotate = this.rotate;
  this.hitTime = 0;
  this.r = 0;
  this.w = -0.0001;

  this.animate = function(tick){
    if (this.isDead){
      this.opacity = Math.max(this.opacity - 0.1 * tick * 0.05, 0);
      this.radius += 0.3 * tick * 0.05;
      if (this.opacity === 0){
        system.removeObject(this.id);
        return;
      }
    }
    if (this.hitTime>0){ // hit effect
      this.hitTime -= 0.1 * tick * 0.05;
      this.w= 1;
    }
    else{
      this.hitTime = 0;
      if (this.w>0) this.w = Math.max(this.w - 0.7 * tick * 0.05,0);
      if (this.w==0 && this.r==0){
        this.r = 0.8;
        this.w = -0.0001;
      }
    }
    this.r= Math.max(this.r - 0.2 * tick * 0.05,0);

    let ccw = Math.cos(this.rotate)*Math.sin(this.imRotate)-Math.sin(this.rotate)*Math.cos(this.imRotate);
    let a = -((Math.cos(this.rotate)*Math.cos(this.imRotate)) + (Math.sin(this.rotate)*Math.sin(this.imRotate))-1) * Math.PI / 2;

    if (ccw > 0){
      this.imRotate-= a / 3;
    } else if (ccw < 0){
      this.imRotate+= a / 3;
    }
  }
  this.setColor = function (color){
    this.color = color;
  }
  this.setDead = function(dead){
    this.isDead = dead;
  }
  this.hit = function(){
    this.hitTime=0.1;
  }
  this.setCanvasSize = function(camera){
    let xx = ((this.x - this.dx - camera.x) * camera.z) - Math.floor((this.x - this.dx - camera.x) * camera.z);
    let yy = ((this.y - this.dy - camera.y) * camera.z) - Math.floor((this.y - this.dy - camera.y) * camera.z);
    this.canvasSize.x = ((this.radius * 3) * camera.z);
    this.canvasSize.y = ((this.radius * 3) * camera.z);
    this.canvasPos = {x:(this.radius * 1.5 * camera.z) + xx,y:(this.radius * 1.5 * camera.z) + yy};
    this.canvas.width = this.canvasSize.x + 4 * camera.z + 6;
    this.canvas.height = this.canvasSize.y + 4 * camera.z + 6;
    this.canvasPos.x += 2 * camera.z + 3;
    this.canvasPos.y += 2 * camera.z + 3;
    this.ctx.lineWidth = 2 * camera.z;
    this.ctx.lineCap = "round";
    this.ctx.lineJoin = "round";
  }
  this.draw = function(ctx,camera){
    if (this.opacity>=1){
      ctx.save();
      ctx.strokeStyle = this.color.getDarkRGB().getRedRGB(this.r).getLightRGB(this.w).getRGBValue(); // 몸체 그리기
      ctx.fillStyle = this.color.getRedRGB(this.r).getLightRGB(this.w).getRGBValue();
      ctx.lineWidth = 2 * camera.z;
      ctx.lineCap = "round";
      ctx.lineJoin = "round";
      ctx.globalAlpha = 1;
      ctx.beginPath();
      ctx.moveTo((this.x - camera.x + (Math.cos(this.imRotate) * this.radius * 1.5)) * camera.z,(this.y - camera.y + (Math.sin(this.imRotate) * this.radius * 1.5)) * camera.z);
      ctx.lineTo((this.x - camera.x + (Math.cos(this.imRotate + Math.PI / 5 * 2) * this.radius * 1.5)) * camera.z,(this.y - camera.y + (Math.sin(this.imRotate + Math.PI / 5 * 2) * this.radius * 1.5)) * camera.z);
      ctx.lineTo((this.x - camera.x + (Math.cos(this.imRotate + Math.PI / 5 * 4) * this.radius * 1.5)) * camera.z,(this.y - camera.y + (Math.sin(this.imRotate + Math.PI / 5 * 4) * this.radius * 1.5)) * camera.z);
      ctx.lineTo((this.x - camera.x + (Math.cos(this.imRotate - Math.PI / 5 * 4) * this.radius * 1.5)) * camera.z,(this.y - camera.y + (Math.sin(this.imRotate - Math.PI / 5 * 4) * this.radius * 1.5)) * camera.z);
      ctx.lineTo((this.x - camera.x + (Math.cos(this.imRotate - Math.PI / 5 * 2) * this.radius * 1.5)) * camera.z,(this.y - camera.y + (Math.sin(this.imRotate - Math.PI / 5 * 2) * this.radius * 1.5)) * camera.z);
      ctx.lineTo((this.x - camera.x + (Math.cos(this.imRotate) * this.radius * 1.5)) * camera.z,(this.y - camera.y + (Math.sin(this.imRotate) * this.radius * 1.5)) * camera.z);
      ctx.fill();
      ctx.stroke();
      ctx.closePath();
      ctx.restore();
    }
    else{
      this.setCanvasSize(camera);

      this.ctx.strokeStyle = this.color.getDarkRGB().getRedRGB(this.r).getLightRGB(this.w).getRGBValue(); // 몸체 그리기
      this.ctx.fillStyle = this.color.getRedRGB(this.r).getLightRGB(this.w).getRGBValue();
      this.ctx.beginPath();
      this.ctx.moveTo(this.canvasPos.x + (Math.cos(this.imRotate) * this.radius * 1.5) * camera.z,this.canvasPos.y + (Math.sin(this.imRotate) * this.radius * 1.5) * camera.z);
      this.ctx.lineTo(this.canvasPos.x + (Math.cos(this.imRotate + Math.PI / 5 * 2) * this.radius * 1.5) * camera.z,this.canvasPos.y + (Math.sin(this.imRotate + Math.PI / 5 * 2) * this.radius * 1.5) * camera.z);
      this.ctx.lineTo(this.canvasPos.x + (Math.cos(this.imRotate + Math.PI / 5 * 4) * this.radius * 1.5) * camera.z,this.canvasPos.y + (Math.sin(this.imRotate + Math.PI / 5 * 4) * this.radius * 1.5) * camera.z);
      this.ctx.lineTo(this.canvasPos.x + (Math.cos(this.imRotate - Math.PI / 5 * 4) * this.radius * 1.5) * camera.z,this.canvasPos.y + (Math.sin(this.imRotate - Math.PI / 5 * 4) * this.radius * 1.5) * camera.z);
      this.ctx.lineTo(this.canvasPos.x + (Math.cos(this.imRotate - Math.PI / 5 * 2) * this.radius * 1.5) * camera.z,this.canvasPos.y + (Math.sin(this.imRotate - Math.PI / 5 * 2) * this.radius * 1.5) * camera.z);
      this.ctx.lineTo(this.canvasPos.x + (Math.cos(this.imRotate) * this.radius * 1.5) * camera.z,this.canvasPos.y + (Math.sin(this.imRotate) * this.radius * 1.5) * camera.z);
      this.ctx.fill();
      this.ctx.stroke();
      this.ctx.closePath();

      ctx.save();
      ctx.globalAlpha = this.opacity;
      ctx.drawImage(this.canvas,(this.x - camera.x) * camera.z-this.canvasPos.x,(this.y - camera.y) * camera.z-this.canvasPos.y);
      ctx.restore();
    }
  }
}
Pentagon.prototype = new HealthShowObject();
Pentagon.prototype.constructor = Pentagon;
