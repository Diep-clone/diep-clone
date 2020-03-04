function BasicBullet(radius,rotate){
  "use strict";

  DynamicObject.apply(this, arguments);
  this.radius = radius;
  this.rotate = rotate;
  this.color = new RGB(0,176,225);
  this.isDead = false;
  this.canvas = document.createElement("canvas");
  this.ctx = this.canvas.getContext("2d");
  this.canvasSize = {x:0,y:0};
  this.canvasPos = {x:0,y:0};
  this.animate = function(tick){
    if (this.isDead){
      this.opacity = Math.max(this.opacity - 0.13 * tick * 0.05, 0);
      this.radius += 0.4 * tick * 0.05;
      if (this.opacity === 0){
        system.removeObject(this.id);
        return;
      }
    }
  }
  this.setColor = function (color){
    this.color = color;
  }
  this.setDead = function(dead){
    this.isDead = dead;
  }
  this.hit = function(){

  }
  this.gunAnime = function(gun){

  }
  this.setCanvasSize = function(camera){
    let xx = ((this.x - this.dx - camera.x) * camera.z) - Math.floor((this.x - this.dx - camera.x) * camera.z);
    let yy = ((this.y - this.dy - camera.y) * camera.z) - Math.floor((this.y - this.dy - camera.y) * camera.z);
    this.canvasSize.x = ((this.radius * 2) * camera.z);
    this.canvasSize.y = ((this.radius * 2) * camera.z);
    this.canvasPos = {x:(this.radius * camera.z) + xx,y:(this.radius * camera.z) + yy};
    this.canvas.width = this.canvasSize.x + 4 * camera.z + 6;
    this.canvas.height = this.canvasSize.y + 4 * camera.z + 6;
    this.canvasPos.x += 2 * camera.z + 3;
    this.canvasPos.y += 2 * camera.z + 3;
    this.ctx.imageSmoothingEnabled = false;
    this.ctx.lineWidth = 2 * camera.z;
    this.ctx.lineCap = "round";
    this.ctx.lineJoin = "round";
  }
  this.draw = function(ctx,camera){
    if (this.opacity>=1){
      ctx.save();
      ctx.strokeStyle = this.color.getDarkRGB().getRGBValue(); // 몸체 그리기
      ctx.fillStyle = this.color.getRGBValue();
      ctx.lineWidth = 2 * camera.z;
      ctx.lineCap = "round";
      ctx.lineJoin = "round";
      ctx.globalAlpha = 1;
      ctx.beginPath();
      ctx.arc((this.x - camera.x) * camera.z,(this.y - camera.y) * camera.z,this.radius * camera.z,0,Math.PI * 2);
      ctx.fill();
      ctx.stroke();
      ctx.closePath();
      ctx.restore();
    }
    else{
      this.setCanvasSize(camera);

      this.ctx.strokeStyle = this.color.getDarkRGB().getRGBValue(); // 몸체 그리기
      this.ctx.fillStyle = this.color.getRGBValue();
      this.ctx.beginPath();
      this.ctx.arc(this.canvasPos.x,this.canvasPos.y,this.radius * camera.z,0,Math.PI * 2);
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
BasicBullet.prototype = new DynamicObject();
BasicBullet.prototype.constructor = BasicBullet;
