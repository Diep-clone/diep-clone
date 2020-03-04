function DynamicObject(){
  "use strict";

  this.x=0.0;
  this.y=0.0;
  this.dx=0.0;
  this.dy=0.0;
  this.radius = 13.0;
  this.rotate=0;
  this.maxHealth= 10;
  this.health= 10;
  this.opacity=1;
  this.id;
/*
  socket.on('objectSet',function (data){

  });*/

  this.setPosition = function(x,y){
    this.x = x;
    this.y = y;
  }
  this.setDPosition = function(x,y){
    this.dx = x;
    this.dy = y;
  }
  this.setRotate = function(rotate){
    this.rotate = rotate;
  }
  this.setRadius = function(radius){
    this.radius = radius;
  }
  this.setHealth = function(hp,mhp){
    this.health = hp;
    this.maxHealth = mhp;
  }
  this.setOpacity = function (opacity){
    this.opacity = opacity;
  }
  this.setId = function(id){
    this.id = id;
  }
}


function Shape(){
  "use strict";

  this.addRotate;
  this.moveRotate;

  this.animate = function(){
    this.rotate += this.addRotate;
    this.addForce(this.moveRotate);
  }
}
