'use strict';

const util = require('./utility');
const utilObj = require('./objectSet');

function createGun(obj){
  var clone = {};
  for(var i in obj) {
      if(i!=="owner" && typeof(obj[i])==="object" && obj[i] !== null){
          if (Array.isArray(obj[i])) clone[i] = createGunArray(obj[i]);
          else clone[i] = createGun(obj[i]);
      }
      else
          clone[i] = obj[i];
  }
  return clone;
}

function createGunArray(arr){
  var clone = [];
  for (let i=0;i<arr.length;i++){
    if (Array.isArray(arr[i])) clone[i] = createGunArray(arr[i]);
    else clone[i] = createGun(arr[i]);
  }
  return clone;
}

exports.gunSet = function(user,userIndex,io){
  user.isShot = false;
  for (let i=0;i<user.guns.length;i++){
    if (!user.guns[i]) continue;
    for (let j=0;j<user.guns[i].bullets.length;j++){
      if (user.guns[i].bullets[j].isDead){
        user.guns[i].bullets.splice(j--,1);
      }
    }
    if ((user.owner && user.owner.mouse && user.owner.mouse.left) || user.guns[i].autoShot){
      user.guns[i].clickTime += 1000/60;
      let object = null;
      let bulletOwner = (user.guns[i].owner)?user.guns[i].owner:user;
      switch (user.guns[i].gunType){
        case "basic":
          user.isShot = true;
          if (user.guns[i].shotTime <= 0 && user.guns[i].shotPTime < user.guns[i].clickTime/((0.6 - 0.04 * bulletOwner.stats[6]) / user.guns[i].coolTime * 1000)){
            let rotate = user.rotate+user.guns[i].rotate+util.randomRange(-user.guns[i].rotateDistance,user.guns[i].rotateDistance);
            object = {
              objType: 'bullet',
              type: user.guns[i].bulletType,
              owner: bulletOwner,
              id: objID(),
              exp: 0,
              team: bulletOwner.team,
              guns: createGunArray(user.guns[i].guns),
              x: user.x + Math.cos(user.rotate+user.guns[i].rotate-Math.PI/2) * user.guns[i].pos.x * util.isF(user.radius) + Math.cos(user.rotate+user.guns[i].rotate) * user.guns[i].pos.y * util.isF(user.radius),
              y: user.y + Math.sin(user.rotate+user.guns[i].rotate-Math.PI/2) * user.guns[i].pos.x * util.isF(user.radius) + Math.sin(user.rotate+user.guns[i].rotate) * user.guns[i].pos.y * util.isF(user.radius),
              rotate: rotate,
              dx: Math.cos(rotate) * (4 + 0.4 * bulletOwner.stats[3]) * user.guns[i].speed,
              dy: Math.sin(rotate) * (4 + 0.4 * bulletOwner.stats[3]) * user.guns[i].speed,
              speed: (0.056 + 0.02 * bulletOwner.stats[3]) * user.guns[i].speed,
              health: (8 + 6 * bulletOwner.stats[4]) * user.guns[i].health,
              maxHealth: (8 + 6 * bulletOwner.stats[4]) * user.guns[i].health,
              lastHealth: (8 + 6 * bulletOwner.stats[4]) * user.guns[i].health,
              damage: (7 + 3 * bulletOwner.stats[5]) * user.guns[i].damage,
              radius: (user.guns[i].absoluteRadius)?user.guns[i].radius:0.4 * user.guns[i].radius * util.isF(user.radius),
              bound: user.guns[i].bulletBound,
              stance: user.guns[i].bulletStance,
              invTime: -1,
              opacity: 1,
              moveAi: user.guns[i].bulletAi,
              variable: createGun(user.guns[i].variable),
              spawnTime: Date.now(),
              hitTime: Date.now(),
              deadTime: -1,
              time: 1000 * user.guns[i].life,
              event: createGun(user.guns[i].event),
              hitObject: null,
              isBorder : false,
              isCollision: false,
              isDead: false,
              isMove: true,
              isOwnCol: user.guns[i].isOwnCol * 1000
            };
            user.dx -= Math.cos(user.rotate+user.guns[i].rotate) * 0.1 * user.guns[i].bound;
            user.dy -= Math.sin(user.rotate+user.guns[i].rotate) * 0.1 * user.guns[i].bound;
            user.guns[i].shotTime = (0.6 - 0.04 * bulletOwner.stats[6]) / user.guns[i].coolTime * 1000;
          }
        break;
        case "drone":
          if (user.guns[i].shotTime <= 0  && user.guns[i].shotPTime < user.guns[i].clickTime/((0.6 - 0.04 * bulletOwner.stats[6]) / user.guns[i].coolTime * 1000) && user.guns[i].bulletLimit>0){
            user.guns[i].bulletLimit--;
            let rotate = user.rotate+user.guns[i].rotate+util.randomRange(-user.guns[i].rotateDistance,user.guns[i].rotateDistance);
            object = {
              objType: 'drone',
              type: user.guns[i].bulletType,
              owner: bulletOwner,
              id: objID(),
              team: bulletOwner.team,
              exp: 0,
              guns: createGunArray(user.guns[i].guns),
              x: user.x + Math.cos(user.rotate+user.guns[i].rotate-Math.PI/2) * user.guns[i].pos.x * util.isF(user.radius) + Math.cos(user.rotate+user.guns[i].rotate) * user.guns[i].pos.y * util.isF(user.radius),
              y: user.y + Math.sin(user.rotate+user.guns[i].rotate-Math.PI/2) * user.guns[i].pos.x * util.isF(user.radius) + Math.sin(user.rotate+user.guns[i].rotate) * user.guns[i].pos.y * util.isF(user.radius),
              rotate: rotate,
              dx: Math.cos(rotate) * 2.5 * user.guns[i].speed,
              dy: Math.sin(rotate) * 2.5 * user.guns[i].speed,
              speed: (0.056 + 0.02 * bulletOwner.stats[3]) * user.guns[i].speed,
              health: (8 + 6 * bulletOwner.stats[4]) * user.guns[i].health,
              maxHealth: (8 + 6 * bulletOwner.stats[4]) * user.guns[i].health,
              lastHealth: (8 + 6 * bulletOwner.stats[4]) * user.guns[i].health,
              damage: (7 + 3 * bulletOwner.stats[5]) * user.guns[i].damage,
              radius: 0.4 * user.guns[i].radius * util.isF(user.radius),
              bound: user.guns[i].bulletBound,
              stance: user.guns[i].bulletStance,
              invTime: -1,
              opacity: 1,
              moveAi: user.guns[i].bulletAi,
              event: createGun(user.guns[i].event),
              variable: createGun(user.guns[i].variable),
              goTank: false,
              goEnemy: null,
              spawnTime: Date.now(),
              hitTime: Date.now(),
              deadTime: -1,
              hitObject: null,
              isBorder : true,
              isCollision: false,
              isDead: false,
              isMove: true,
              isOwnCol: true
            };
            user.dx -= Math.cos(user.rotate+user.guns[i].rotate) * 0.1 * user.guns[i].bound;
            user.dy -= Math.sin(user.rotate+user.guns[i].rotate) * 0.1 * user.guns[i].bound;
            user.guns[i].shotTime = (0.6 - 0.04 * bulletOwner.stats[6]) / user.guns[i].coolTime * 1000;
          }
        break;
        case "auto":
        break;
        default:
        break;
      }
      if (object){
        if (object.type === 2 || object.type === 3){
          user.guns[i].bullets.push(object);
        }
        objects[userIndex]=object;
        userIndex=objects.length;
        objects.push(user);
        io.emit("shot",user.id,i);
      }
    }
    else{
      user.guns[i].clickTime = 0;
    }
    user.guns[i].shotTime = Math.max(user.guns[i].shotTime-1000/60,0);
  }
}
