function HealthShowObject(){
  "use strict";
  DynamicObject.apply(this, arguments);

  this.showPercent = 1;

  this.drawHPBar = function(ctx,camera){
    let healthPercent = this.health/this.maxHealth;

    ctx.save();
    ctx.globalAlpha = this.opacity;

    this.showPercent -= (this.showPercent - healthPercent) / 3;

    if (healthPercent<1){
      ctx.beginPath();
      ctx.moveTo((this.x - camera.x + this.radius) * camera.z,(this.y - camera.y + this.radius * 5 / 3) * camera.z);
      ctx.lineTo((this.x - camera.x - this.radius) * camera.z,(this.y - camera.y + this.radius * 5 / 3) * camera.z);
      ctx.closePath();
      ctx.strokeStyle = "#444444";
      ctx.lineWidth = 4.1 * camera.z;
      ctx.stroke();

      ctx.beginPath();
      ctx.moveTo((this.x - camera.x - this.radius) * camera.z,(this.y - camera.y + this.radius * 5 / 3) * camera.z);
      ctx.lineTo((this.x - camera.x - this.radius + this.showPercent * this.radius * 2) * camera.z,(this.y - camera.y + this.radius * 5 / 3) * camera.z);
      ctx.closePath();
      ctx.strokeStyle = "#86e27f";
      ctx.lineWidth = 2.6 * camera.z;
      ctx.stroke();
    }
    ctx.restore();
  }
}
HealthShowObject.prototype = new DynamicObject();
HealthShowObject.prototype.constructor = HealthShowObject;
