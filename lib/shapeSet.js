'use strict';

const util = require('./utility');
const utilObj = require('./objectSet');

let maxShapes = 200;
let shapeCount = 0;

exports.spawnShape = function (mapSize){
  while (maxShapes>shapeCount){
    let xx = util.randomRange(-mapSize.x,mapSize.x);
    let yy = util.randomRange(-mapSize.y,mapSize.y);
    let type;
    let ran = util.randomRange(1,100);
    if (xx<mapSize.x/3 && xx>-mapSize.x/3 && yy<mapSize.y/3 && yy>-mapSize.y/3){
      if (ran<70) type = 2;
      else if (ran<85) type = 4;
      else if (ran<95) type = 5;
      else type = 3;
    }
    else{
      if (ran<70) type = 0;
      else if (ran<90) type = 1;
      else type = 2;
    }
    let name = [
      "square",
      "triangle",
      "pentagon",
      "alphaPentagon",
      "triangle",
      "triangle"
    ];
    let exp = [10,25,130,3000,15,25];
    let health = [10,30,100,3000,10,30];
    let damage = [8,8,12,20,8,8];
    let radius = [10,10,12,30,5,8];
    let stance = [0.2,0.2,0.04,0.002,0.2,0.2];

    objects.push({
      objType:'shape',
      type:type,
      owner: null,
      id:objID(),
      team: 'ene',
      x:xx,
      y:yy,
      dx:0,
      dy:0,
      exp:exp[type],
      health:health[type],
      maxHealth:health[type],
      lastHealth:health[type],
      damage:damage[type],
      radius:radius[type],
      rotate:util.randomRange(-Math.PI,Math.PI),
      bound:1,
      stance:stance[type],
      opacity: 1,
      name: name[type],
      isOwnCol: true,
      spawnTime: Date.now(),
      hitTime: Date.now(),
      deadTime: -1,
      event:{
        killEvent:function(deader){
          return true;
        },
        deadEvent:function(killer){
          shapeCount--;
          return true;
        }
      },
      hitObject: null,
      moveAi:null,
      isBorder : true,
      isCollision:false,
      isMove: true,
      isDead:false
    });
    shapeCount++;
  }
}

exports.extendMaxShape = function (n){
  maxShapes+=n;
}

exports.moveShape = function (s,mapSize,p,target){
  s.moveAi(s,mapSize,p,target);
  utilObj.moveObject(s);
}
