function SkimmerBullet(radius,rotate){
  "use strict";

  DynamicObject.apply(this, arguments);
  this.radius = radius;
  this.rotate = rotate;
  this.color = new RGB(0,176,225);
  this.isDead = false;
  this.canvas = document.createElement("canvas");
  this.ctx = this.canvas.getContext("2d");
  this.guns = [new BulletGun(this.radius,[[5,-2],[5,4.5],[-5,4.5],[-5,-2]],0),new BulletGun(this.radius,[[5,-2],[5,4.5],[-5,4.5],[-5,-2]],Math.PI)];
  this.imRotate = this.rotate;
  this.showRadius = this.radius;
  this.canvasSize = {x:0,y:0};
  this.canvasPos = {x:0,y:0};
  this.r = 0;
  this.w = 0;
  this.animate = function(tick){
    if (this.isDead){
      this.opacity = Math.max(this.opacity - 0.13 * tick * 0.05, 0);
      this.radius += 0.4 * tick * 0.05;
      if (this.opacity === 0){
        system.removeObject(this.id);
        return;
      }
    }
    for (let i=0;i<this.guns.length;i++){
      this.guns[i].animate();
    }
    this.imRotate = this.rotate;
    this.showRadius = this.radius;
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
    if (gun<this.guns.length) this.guns[gun].shot();
  }
  this.setCanvasSize = function(camera){
    let xx = ((this.x - this.dx - camera.x) * camera.z) - Math.floor((this.x - this.dx - camera.x) * camera.z);
    let yy = ((this.y - this.dy - camera.y) * camera.z) - Math.floor((this.y - this.dy - camera.y) * camera.z);
    this.canvasSize.x = ((this.showRadius * 2) * camera.z);
    this.canvasSize.y = ((this.showRadius * 2) * camera.z);
    this.canvasPos = {x:(this.showRadius * camera.z) + xx,y:(this.showRadius * camera.z) + yy};
    for (let i=0;i<this.guns.length;i++){
      this.guns[i].setParentCanvasSize(this,camera);
    }
    this.canvas.width = this.canvasSize.x + 4 * camera.z + 4;
    this.canvas.height = this.canvasSize.y + 4 * camera.z + 4;
    this.canvasPos.x += 2 * camera.z + 2;
    this.canvasPos.y += 2 * camera.z + 2;
    this.ctx.imageSmoothingEnabled = false;
    this.ctx.lineWidth = 2 * camera.z;
    this.ctx.lineCap = "round";
    this.ctx.lineJoin = "round";
  }
  this.draw = function(ctx,camera){
    if (this.opacity>=1){
      this.canvasPos = {x:(this.x - this.dx - camera.x) * camera.z,y:(this.y - this.dy - camera.y) * camera.z};
      ctx.save();
      ctx.strokeStyle = this.color.getDarkRGB().getRGBValue(); // 몸체 그리기
      ctx.fillStyle = this.color.getRGBValue();
      ctx.lineWidth = 2 * camera.z;
      ctx.lineCap = "round";
      ctx.lineJoin = "round";
      ctx.globalAlpha = 1;
      for (let i=0;i<this.guns.length;i++){
        this.guns[i].drawGun(this,ctx,camera);
      }
      ctx.beginPath();
      ctx.arc((this.x - camera.x) * camera.z,(this.y - camera.y) * camera.z,this.showRadius * camera.z,0,Math.PI * 2);
      ctx.fill();
      ctx.stroke();
      ctx.closePath();
      ctx.restore();
    }
    else{
      this.setCanvasSize(camera);

      for (let i=0;i<this.guns.length;i++){
        this.guns[i].drawGun(this,this.ctx,camera);
      }

      this.ctx.strokeStyle = this.color.getDarkRGB().getRGBValue(); // 몸체 그리기
      this.ctx.fillStyle = this.color.getRGBValue();
      this.ctx.beginPath();
      this.ctx.arc(this.canvasPos.x,this.canvasPos.y,this.showRadius * camera.z,0,Math.PI * 2);
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
SkimmerBullet.prototype = new DynamicObject();
SkimmerBullet.prototype.constructor = SkimmerBullet;

function RocketeerBullet(radius,rotate){
  "use strict";

  DynamicObject.apply(this, arguments);
  this.radius = radius;
  this.rotate = rotate;
  this.color = new RGB(0,176,225);
  this.isDead = false;
  this.canvas = document.createElement("canvas");
  this.ctx = this.canvas.getContext("2d");
  this.guns = [new BulletGun(this.radius,[[6.5,-2],[8.5,4.6],[-8.5,4.6],[-6.5,-2]],Math.PI)];
  this.imRotate = this.rotate;
  this.showRadius = this.radius;
  this.canvasSize = {x:0,y:0};
  this.canvasPos = {x:0,y:0};
  this.r = 0;
  this.w = 0;
  this.animate = function(tick){
    if (this.isDead){
      this.opacity = Math.max(this.opacity - 0.13 * tick * 0.05, 0);
      this.radius += 0.4 * tick * 0.05;
      if (this.opacity === 0){
        system.removeObject(this.id);
        return;
      }
    }
    for (let i=0;i<this.guns.length;i++){
      this.guns[i].animate();
    }
    this.imRotate = this.rotate;
    this.showRadius = this.radius;
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
    if (gun<this.guns.length) this.guns[gun].shot();
  }
  this.setCanvasSize = function(camera){
    let xx = ((this.x - this.dx - camera.x) * camera.z) - Math.floor((this.x - this.dx - camera.x) * camera.z);
    let yy = ((this.y - this.dy - camera.y) * camera.z) - Math.floor((this.y - this.dy - camera.y) * camera.z);
    this.canvasSize.x = ((this.showRadius * 2) * camera.z);
    this.canvasSize.y = ((this.showRadius * 2) * camera.z);
    this.canvasPos = {x:(this.showRadius * camera.z) + xx,y:(this.showRadius * camera.z) + yy};
    for (let i=0;i<this.guns.length;i++){
      this.guns[i].setParentCanvasSize(this,camera);
    }
    this.canvas.width = this.canvasSize.x + 4 * camera.z + 4;
    this.canvas.height = this.canvasSize.y + 4 * camera.z + 4;
    this.canvasPos.x += 2 * camera.z + 2;
    this.canvasPos.y += 2 * camera.z + 2;
    this.ctx.imageSmoothingEnabled = false;
    this.ctx.lineWidth = 2 * camera.z;
    this.ctx.lineCap = "round";
    this.ctx.lineJoin = "round";
  }
  this.draw = function(ctx,camera){
    if (this.opacity>=1){
      this.canvasPos = {x:(this.x - this.dx - camera.x) * camera.z,y:(this.y - this.dy - camera.y) * camera.z};
      ctx.save();
      ctx.strokeStyle = this.color.getDarkRGB().getRGBValue(); // 몸체 그리기
      ctx.fillStyle = this.color.getRGBValue();
      ctx.lineWidth = 2 * camera.z;
      ctx.lineCap = "round";
      ctx.lineJoin = "round";
      ctx.globalAlpha = 1;
      for (let i=0;i<this.guns.length;i++){
        this.guns[i].drawGun(this,ctx,camera);
      }
      ctx.beginPath();
      ctx.arc((this.x - camera.x) * camera.z,(this.y - camera.y) * camera.z,this.showRadius * camera.z,0,Math.PI * 2);
      ctx.fill();
      ctx.stroke();
      ctx.closePath();
      ctx.restore();
    }
    else{
      this.setCanvasSize(camera);

      for (let i=0;i<this.guns.length;i++){
        this.guns[i].drawGun(this,this.ctx,camera);
      }

      this.ctx.strokeStyle = this.color.getDarkRGB().getRGBValue(); // 몸체 그리기
      this.ctx.fillStyle = this.color.getRGBValue();
      this.ctx.beginPath();
      this.ctx.arc(this.canvasPos.x,this.canvasPos.y,this.showRadius * camera.z,0,Math.PI * 2);
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
RocketeerBullet.prototype = new DynamicObject();
RocketeerBullet.prototype.constructor = RocketeerBullet;
