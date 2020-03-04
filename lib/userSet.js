'use strict';

const util = require('./utility');
const utilObj = require('./objectSet');

exports.healTank = function (t){
  if (Date.now()-30000>t.hitTime){
    utilObj.healObject(t);
  }
  else{
    t.health=Math.min(t.health+util.isF(t.maxHealth)/60/30*(0.03 + 0.12 * t.stats[0]),util.isF(t.maxHealth));
  }
}

exports.afkTank = function(t){
  t.health-=util.isF(t.maxHealth)/60/10;
}

exports.setUserLevel = function(user){
  let statList = "011111111111111111111111111101001001001001001001";
  let expList = [
    0,
    4,
    13,
    28,
    50,
    78,
    113,
    157,
    211,
    275,
    350,
    437,
    538,
    655,
    787,
    938,
    1109,
    1301,
    1516,
    1757,
    2026,
    2325,
    2658,
    3026,
    3433,
    3883,
    4379,
    4925,
    5525,
    6184,
    6907,
    7698,
    8537,
    9426,
    10368,
    11367,
    12426,
    13549,
    14739,
    16000,
    17337,
    18754,
    20256,
    21849,
    23536
  ];
  for (let i=user.level;expList[i]<=user.exp && i<45;i++){
    user.stat+=statList[i] * 1;
    user.level=i+1;
  }
  return expList[user.level];
}

exports.setUserSight = function(user){
  //let sight = 1.786 - user.level * 0.008;
  let sight = 16/9*Math.pow(0.995,user.level-1);
  switch(user.type){
    case 6: // Sniper
    case 11: // Overseer
    case 12: // Overload
    case 17: // Necromanser
    case 26: // Manager
    case 30: // Trapper
    case 31: // GunnerTrapper
    case 32: // OverTrapper
    case 33: // MegaTrapper
    case 34: // TriTrapper
    case 35: // Smasher
    case 36: // Landmine
    case 42: // AutoTrapper
    case 46: // BattleShip
    case 48: // AutoSmasher
    case 49: // Spike
    case 50: // Factory
    case 51: // Skimmer
    case 52: // Rocketeer
    return sight / 1.11;
    break;
    case 19: // Hunter
    case 28: // Predator
    case 41: // Streamliner
    return sight / 1.176;
    break;
    case 15: // Assasin
    case 21: // Stakler
    return sight / 1.25;
    break;
    case 22: // Ranger
    return sight / 1.428;
    break;
    default:
    return sight;
    break;
  }
}

function isUpgradeTank(user,socket){

}

function setStatsToMaxStats(user){
  for (let i=0;i<8;i++){
    if (user.stats[i]>user.maxStats[i]){
      user.stat+=user.stats[i]-user.maxStats[i];
      user.stats[i]=user.maxStats[i];
    }
  }
}

/*
    {
      bulletType:int,
      speed:%,
      damage:%,
      health:%,
      rotate:pi,
      rotateDistance:pi,
      radius:%,
      bound:%,
      coolTime:%,
      shotTime:0.0f,
      shotPTime:%,
      autoShot:true/false,
      life:int,
      pos:{x:%,y:%},
      dir:{rotate:pi,detect:pi,distance:%}
    }
*/
/*
function getGun(data){
  let gun = {
    objType:"bullet",
    bulletType: 1,
    speed: 1,
    damage: 1,
    health: 1,
    rotate: 0,
    rotateDistance: Math.PI/72,
    guns: [],
    radius: 1,
    bound: 1,
    bulletBound: 1,
    coolTime: 1, // 클타임 비율
    shotTime: 0, // 총알 쏘기까지 남은 시간
    shotPTime: 0, // 얼마나 눌러야 발사되는가
    clickTime: 0, // 총알을 쏴야 한다고 판단하기 시작한 후의 경과시간
    autoShot: false,
    life: 3,
    bulletCount: 0,
    bulletLimit: 0,
    isOwnCol: 0,
    bulletAi: function(b,mapSize,p,target){
      b.dx+=Math.cos(b.rotate) * b.speed;
      b.dy+=Math.sin(b.rotate) * b.speed;
    },
    pos: {x:0,y:1.88},
    dir: {rotate:null,detect:Math.PI,distance:0}
  };
  for (let key in data){
    gun[key]=data[key];
  }
  return gun;
}*/

function getBasicGun(data){
  let gun = {
    gunType: "basic",
    objType: "bullet",
    bulletType: 1,
    speed: 1,
    damage: 1,
    health: 1,
    rotate: 0,
    rotateDistance: Math.PI/36,
    guns: [],
    radius: 1,
    absoluteRadius: false,
    bound: 1,
    bulletBound: 1,
    bulletStance: 0.1,
    coolTime: 1,
    shotTime: 0,
    shotPTime: 0,
    clickTime: 0,
    isOwnCol: 0,
    autoShot: false,
    life: 3,
    bullets:[],
    bulletAi: function(b){
      b.dx+=Math.cos(b.rotate) * b.speed;
      b.dy+=Math.sin(b.rotate) * b.speed;
    },
    event:{
      killEvent:function(u,deader){
        u.owner.exp += deader.exp;
        return true;
      },
      deadEvent:function(u,killer){
        return true;
      }
    },
    pos: {x:0,y:1.88}
  };
  for (let key in data){
    gun[key]=data[key];
  }
  return gun;
}

function getDroneGun(data){
  let gun = {
    gunType:"drone",
    objType:"drone",
    bulletType: 2,
    speed: 1,
    damage: 1,
    health: 1,
    rotate: 0,
    rotateDistance: Math.PI/36,
    guns:[],
    radius: 1,
    absoluteRadius: false,
    bound: 1,
    bulletBound: 1,
    bulletStance: 1,
    coolTime: 1, // 클타임 비율
    shotTime: 0, // 총알 쏘기까지 남은 시간
    shotPTime: 0, // 얼마나 눌러야 발사되는가
    clickTime: 0, // 총알을 쏴야 한다고 판단하기 시작한 후의 경과시간
    autoShot: true,
    isOwnCol: 1,
    bulletLimit: 4,
    bullets:[],
    bulletAi: function(b){
      let p = b.owner.owner;
      let target = null;
      if (p && p.mouse.right){
        b.rotate = Math.atan2(b.y-(b.owner.y+p.target.y),b.x-(b.owner.x+p.target.x));
        b.dx += Math.cos(b.rotate) * b.speed;
        b.dy += Math.sin(b.rotate) * b.speed;
        b.goTank = false;
      }
      else if (p && p.mouse.left){
        b.rotate = Math.atan2((b.owner.y+p.target.y)-b.y,(b.owner.x+p.target.x)-b.x);
        b.dx += Math.cos(b.rotate) * b.speed;
        b.dy += Math.sin(b.rotate) * b.speed;
        b.goTank = false;
      }
      else{
        let dis = Math.sqrt((b.owner.x-b.x)*(b.owner.x-b.x)+(b.owner.y-b.y)*(b.owner.y-b.y));

        if (b.goTank){
          if (dis<60) b.goTank = false;
          else{
            b.rotate = Math.atan2(b.owner.y-b.y,b.owner.x-b.x);
            b.dx += Math.cos(b.rotate) * b.speed;
            b.dy += Math.sin(b.rotate) * b.speed;
          }
        }
        else if (b.goEnemy){
          if (b.goEnemy.isDead){
            b.goEnemy = target;
          }
          else if (dis>417.96){
            b.goEnemy = null;
          }
          else{
            b.rotate = Math.atan2(b.goEnemy.y-b.y,b.goEnemy.x-b.x);
            b.dx += Math.cos(b.rotate) * b.speed;
            b.dy += Math.sin(b.rotate) * b.speed;
          }
        }
        else if (dis<util.randomRange(139.32,150)){ // 플레이어 주변에 있으면
          let dir = Math.atan2(b.owner.y-b.y,b.owner.x-b.x) - Math.PI / 2;
          b.dx += Math.cos(dir) * 0.02;
          b.dy += Math.sin(dir) * 0.02;
          b.rotate=Math.atan2(b.dy,b.dx);
          if (target) b.goEnemy = target;
        }
        else{
          b.goTank = true;
        }
      }
    },
    event:{
      killEvent:function(u,deader){
        u.owner.exp += deader.exp;
        return true;
      },
      deadEvent:function(u,killer){
        gun.bulletLimit++;
        return true;
      }
    },
    pos: {x:0,y:1.88}
  };
  for (let key in data){
    gun[key]=data[key];
  }
  return gun;
}

function getAutoGun(data){
  let gun = {
    gunType:"auto",
    bulletType: 1,
    speed: 1,
    damage: 1,
    health: 1,
    rotate: 0,
    rotateDistance: Math.PI/36,
    guns: [],
    radius: 1,
    absoluteRadius: false,
    bound: 1,
    bulletBound: 1,
    bulletStance: 0.1,
    coolTime: 1, // 클타임 비율
    shotTime: 0, // 총알 쏘기까지 남은 시간
    shotPTime: 0, // 얼마나 눌러야 발사되는가
    clickTime: 0, // 총알을 쏴야 한다고 판단하기 시작한 후의 경과시간
    autoShot: true,
    life: 3,
    bulletCount: 0,
    bulletLimit: 0,
    isOwnCol: 0,
    bullets:[],
    bulletAi: function(b){
      b.dx+=Math.cos(b.rotate) * b.speed;
      b.dy+=Math.sin(b.rotate) * b.speed;
    },
    pos: {x:0,y:1},
    dir: {rotate:null,detect:Math.PI,distance:0}
  };
  for (let key in data){
    gun[key]=data[key];
  }
  return gun;
}

exports.setUserTank = function(user){ // 유저 총구 설정.
  let droneAi = function(b){
    let target = null;
    if (b.owner){
      let dis = Math.sqrt((b.owner.x-b.x)*(b.owner.x-b.x)+(b.owner.y-b.y)*(b.owner.y-b.y));

      if (b.goTank){
        if (dis<60) b.goTank = false;
        else{
          b.rotate = Math.atan2(b.owner.y-b.y,b.owner.x-b.x);
          b.dx += Math.cos(b.rotate) * b.speed;
          b.dy += Math.sin(b.rotate) * b.speed;
        }
      }
      else if (b.goEnemy){
        if (b.goEnemy.isDead){
          b.goEnemy = target;
        }
        else if (dis>417.96){
          b.goEnemy = null;
        }
        else{
          b.rotate = Math.atan2(b.goEnemy.y-b.y,b.goEnemy.x-b.x);
          b.dx += Math.cos(b.rotate) * b.speed;
          b.dy += Math.sin(b.rotate) * b.speed;
        }
      }
      else if (dis<util.randomRange(139.32,150)){ // 플레이어 주변에 있으면
        let dir = Math.atan2(b.owner.y-b.y,b.owner.x-b.x) - Math.PI / 2;
        b.dx += Math.cos(dir) * 0.02;
        b.dy += Math.sin(dir) * 0.02;
        b.rotate=Math.atan2(b.dy,b.dx);
        if (target) b.goEnemy = target;
      }
      else{
        b.goTank = true;
      }
    }
  };

  for (let i=0;i<user.guns.length;i++){
    if (!user.guns[i]) continue;
    for (let j=0;j<user.guns[i].bullets.length;j++){
      user.guns[i].bullets[j].isDead = true;
      if (user.guns[i].bullets[j].event)
        user.guns[i].bullets[j].event.deadEvent = null;
    }
  }
  user.speed = function (){return (0.07 + (0.007 * user.stats[7])) * Math.pow(0.985,user.level-1);};
  user.damage = function (){return 20 + user.stats[2] * 4;};
  user.radius = function (){return 13*Math.pow(1.01,(user.level-1));};
  user.maxHealth = function (){return 48 + user.level * 2 + user.stats[1] * 20;};
  user.bound = 1;
  user.stance = 1;
  user.guns = [];
  user.maxStats = [7,7,7,7,7,7,7,7];
  user.invTime=-1;
  user.isOwnCol=false;
  user.isCanDir=true;
  user.event = { // 여기 있는 값들은 모두 "함수" 입니다.
    rightDownEvent:function(){}, // 우클릭을 시작했을 때의 이벤트
    rightEvent:function(){}, // 우클릭을 진행할 때의 이벤트
    rightUpEvent:function(){}, // 우클릭을 끝낼 때의 이벤트
    killEvent:function(u,deader){ // 오브젝트를 죽였을 때의 이벤트
      u.exp += deader.exp;
      return true;
    },
    deadEvent:function(u,killer){
      return true;
    } // 자신의 오브젝트가 죽었을 때의 이벤트
  };


  switch(user.type){
      case 0: // Basic
      user.guns = [getBasicGun()];
      break;
      case 1: // Twin
      user.guns = [getBasicGun({
        damage:0.65,
        bound:0.75,
        pos:{x:0.5,y:1.88}
      }),getBasicGun({
        damage:0.65,
        bound:0.75,
        shotPTime:0.5,
        pos:{x:-0.5,y:1.88}
      })];
      break;
      case 2: // Triplet
      user.guns = [getBasicGun({
        damage:0.6,
        health:0.7,
        bound:0.5,
        shotPTime:0.5,
        pos:{x:0.5,y:1.6}
      }),getBasicGun({
        damage:0.6,
        health:0.7,
        bound:0.5,
        shotPTime:0.5,
        pos:{x:-0.5,y:1.6}
      }),getBasicGun({
        damage:0.6,
        health:0.7,
        bound:0.5
      })];
      break;
      case 3: // TripleShot
      user.guns = [getBasicGun({
        damage:0.7
      }),getBasicGun({
        damage:0.7,
        rotate:Math.PI / 4
      }),getBasicGun({
        damage:0.7,
        rotate:-Math.PI / 4
      })];
      break;
      case 4: // QuadTank
      user.guns = [];
      for (let i=-Math.PI/2;i<=Math.PI;i+=Math.PI/2){
        user.guns.push(getBasicGun({
          damage:0.75,
          rotate:i
        }));
      }
      break;
      case 5: // OctoTank
      user.guns = [];
      for (let i=-Math.PI/2;i<=Math.PI;i+=Math.PI/2){
        user.guns.push(getBasicGun({
          damage:0.65,
          rotate:i
        }));
      }
      for (let i=-Math.PI/4*3;i<=Math.PI;i+=Math.PI/2){
        user.guns.push(getBasicGun({
          damage:0.65,
          rotate:i,
          shotPTime:0.5
        }));
      }
      break;
      case 6: // Sniper
      user.guns = [getBasicGun({
          speed:1.5,
          bound:3,
          rotateDistance:Math.PI/120,
          coolTime:0.667,
          pos:{x:0,y:2.2}
      })];
      break;
      case 7: // MachineGun
      user.guns = [getBasicGun({
          damage:0.7,
          coolTime:2,
          rotateDistance:Math.PI/8,
          pos:{x:0,y:1.6}
      })];
      break;
      case 8: // FlankGuard
      user.guns = [getBasicGun(),getBasicGun({
          rotate:Math.PI,
          pos:{x:0,y:1.6},
      })];
      break;
      case 9: // TriAngle
      user.guns = [getBasicGun({
          bound:0.2
      }),getBasicGun({
          damage:0.2,
          rotate:Math.PI / 6 * 5,
          bound:2.5,
          shotPTime:0.5,
          life:1.5,
          pos:{x:0,y:1.6}
      }),getBasicGun({
          damage:0.2,
          rotate:-Math.PI / 6 * 5,
          bound:2.5,
          shotPTime:0.5,
          life:1.5,
          pos:{x:0,y:1.6}
      })];
      break;
      case 10: // Destroyer
      user.guns = [getBasicGun({
          speed:0.75,
          damage:3,
          health:2,
          radius:1.75,
          bound:15,
          bulletBound:0.1,
          coolTime:0.25,
          pos:{x:0,y:1.88}
      })];
      break;
      case 11: // Overseer
        user.guns = [getDroneGun({
          bulletType:2,
          speed:0.6,
          damage:0.7,
          health:2,
          rotate:Math.PI/2,
          bulletBound:0.4,
          coolTime:0.167,
          autoShot:true,
          bulletLimit:4,
          isOwnCol:true,
          pos:{x:0,y:1.6}
        }),getDroneGun({
          bulletType:2,
          speed:0.6,
          damage:0.7,
          health:2,
          rotate:-Math.PI/2,
          bulletBound:0.4,
          coolTime:0.167,
          autoShot:true,
          bulletLimit:4,
          isOwnCol:true,
          pos:{x:0,y:1.6}
        })];
      break;
      case 12: // Overload
        let dir = -Math.PI/2;
        user.guns = [];
        for (let i=0;i<4;i++){
          user.guns.push(getDroneGun({
            objType:"drone",
            bulletType:2,
            speed:0.6,
            damage:0.7,
            health:2,
            rotate:dir,
            bulletBound:0.4,
            coolTime:0.167,
            autoShot:true,
            bulletLimit:2,
            isOwnCol:true,
            pos:{x:0,y:1.6}
          }));
          dir+=Math.PI/2;
        }
      break;
      case 13: // TwinFlank
      user.guns = [getBasicGun({
        damage:0.5,
        pos:{x:0.5,y:1.88}
      }),getBasicGun({
        damage:0.5,
        shotPTime:0.5,
        pos:{x:-0.5,y:1.88}
      }),getBasicGun({
        damage:0.5,
        rotate:Math.PI,
        pos:{x:0.5,y:1.88}
      }),getBasicGun({
        damage:0.5,
        rotate:Math.PI,
        shotPTime:0.5,
        pos:{x:-0.5,y:1.88}
      })];
      break;
      case 14: // PentaShot
      user.guns = [getBasicGun({
        damage:0.6,
        bound:0.7,
        shotPTime:0.667,
        rotate:Math.PI / 4,
        pos:{x:0,y:1.6}
      }),getBasicGun({
        damage:0.6,
        bound:0.7,
        shotPTime:0.667,
        rotate:-Math.PI / 4,
        pos:{x:0,y:1.6}
      }),getBasicGun({
        damage:0.6,
        bound:0.7,
        shotPTime:0.333,
        rotate:Math.PI / 8,
        pos:{x:0,y:1.9}
      }),getBasicGun({
        damage:0.6,
        bound:0.7,
        shotPTime:0.333,
        rotate:-Math.PI / 8,
        pos:{x:0,y:1.9}
      }),getBasicGun({
        damage:0.6,
        bound:0.7,
        pos:{x:0,y:2.25}
      })];
      break;
      case 15: // Assasin
      user.guns = [getBasicGun({
          speed:1.5,
          bound:3,
          rotateDistance:Math.PI/120,
          coolTime:0.5,
          pos:{x:0,y:2.45}
      })];
      break;
      case 16: // ArenaCloser
      user.guns = [getBasicGun({
        speed:2,
        damage:20,
        health:100,
        bulletBound:50,
        bulletStance:0.001
      })];
      break;
      case 17: // Necromanser
      user.droneCount = 0;
      user.guns = [getDroneGun({
        bulletLimit:0
      })];
      user.variable.droneCount=0;
      user.event = {
        rightDownEvent:function(u){},
        rightEvent:function(u){},
        rightUpEvent:function(u){},
        killEvent:function(u,deader){
          if (deader.objType==="shape" && deader.type===0 && deader.team==="ene" && u.variable.droneCount <= 22 + 2 * u.stats[6]){
            u.variable.droneCount++;
            deader.team = u.team;
            deader.exp = 0;
            deader.owner = u;
            deader.hitObject = null;
            deader.health = deader.maxHealth = (8 + 6 * u.stats[4]) * 0.5;
            deader.damage = (7 + 3 * u.stats[5]) * 1.68;
            deader.speed = (0.056 + 0.01 * u.stats[3]) * 0.5;
            deader.isDead = false;
            deader.deadTime = -1;
            let imF = deader.event.deadEvent;
            deader.event = {
              killEvent:u.event.killEvent,
              deadEvent:function(u,killer){
                imF();
                u.variable.droneCount--;
                return true;
              }
            }
            deader.moveAi = function(b){
              let p = b.owner.owner;
              let target = null;
              if (p && p.mouse.right){
                b.rotate = Math.atan2(b.y-(b.owner.y+p.target.y),b.x-(b.owner.x+p.target.x));
                b.dx += Math.cos(b.rotate) * b.speed;
                b.dy += Math.sin(b.rotate) * b.speed;
                b.goTank = false;
              }
              else if (p && p.mouse.left){
                b.rotate = Math.atan2((b.owner.y+p.target.y)-b.y,(b.owner.x+p.target.x)-b.x);
                b.dx += Math.cos(b.rotate) * b.speed;
                b.dy += Math.sin(b.rotate) * b.speed;
                b.goTank = false;
              }
              else if (b.owner){
                let dis = Math.sqrt((b.owner.x-b.x)*(b.owner.x-b.x)+(b.owner.y-b.y)*(b.owner.y-b.y));

                if (b.goTank){
                  if (dis<70) b.goTank = false;
                  else{
                    b.rotate = Math.atan2(b.owner.y-b.y,b.owner.x-b.x);
                    b.dx += Math.cos(b.rotate) * b.speed;
                    b.dy += Math.sin(b.rotate) * b.speed;
                  }
                }
                else if (b.goEnemy){
                  if (b.goEnemy.isDead || b.goEnemy.owner === b.owner){
                    b.goEnemy = target;
                  }
                  else if (dis>417.96){
                    b.goEnemy = null;
                  }
                  else{
                    b.rotate = Math.atan2(b.goEnemy.y-b.y,b.goEnemy.x-b.x);
                    b.dx += Math.cos(b.rotate) * b.speed;
                    b.dy += Math.sin(b.rotate) * b.speed;
                  }
                }
                else if (dis<util.randomRange(80,100)){ // 플레이어 주변에 있으면
                  let dir = Math.atan2(b.owner.y-b.y,b.owner.x-b.x) - Math.PI / 2;
                  b.dx += Math.cos(dir) * 0.01;
                  b.dy += Math.sin(dir) * 0.01;
                  b.rotate=Math.atan2(b.dy,b.dx);
                  if (target) b.goEnemy = target;
                }
                else{
                  b.goTank = true;
                }
              }
            };
            u.guns[0].bullets.push(deader);
            return false;
          }
          else{
            u.exp += deader.exp;
            return true;
          }
        },
        deadEvent:function(u,killer){
          return true;
        }
      };
      break;
      case 18: // TripleTwin
      user.guns = [getBasicGun({
        damage:0.5,
        pos:{x:0.5,y:1.88}
      }),getBasicGun({
        damage:0.5,
        shotPTime:0.5,
        pos:{x:-0.5,y:1.88}
      }),getBasicGun({
        damage:0.5,
        rotate:Math.PI / 3 * 2,
        pos:{x:0.5,y:1.88}
      }),getBasicGun({
        damage:0.5,
        rotate:Math.PI / 3 * 2,
        shotPTime:0.5,
        pos:{x:-0.5,y:1.88}
      }),getBasicGun({
        damage:0.5,
        rotate:-Math.PI / 3 * 2,
        pos:{x:0.5,y:1.88}
      }),getBasicGun({
        damage:0.5,
        rotate:-Math.PI / 3 * 2,
        shotPTime:0.5,
        pos:{x:-0.5,y:1.88}
      })];
      break;
      case 19: // Hunter
      user.guns = [getBasicGun({
          speed:1.25,
          damage:0.75,
          rotateDistance:Math.PI/120,
          radius:0.7,
          bound:0.03,
          coolTime:0.4,
          pos:{x:0,y:2.23}
      }),getBasicGun({
          speed:1.25,
          damage:0.75,
          rotateDistance:Math.PI/120,
          bound:0.03,
          shotPTime:0.25,
          coolTime:0.4,
          pos:{x:0,y:2.23}
      })];
      break;
      case 20: // Gunner
      user.guns = [getBasicGun({
        damage:0.5,
        health:0.45,
        radius:0.6,
        bound:0.2,
        shotPTime:0.5,
        pos:{x:0.67,y:1.3}
      }),getBasicGun({
        damage:0.5,
        health:0.45,
        radius:0.6,
        bound:0.2,
        shotPTime:0.75,
        pos:{x:-0.67,y:1.3}
      }),getBasicGun({
        damage:0.5,
        health:0.45,
        radius:0.6,
        bound:0.2,
        shotPTime:0,
        pos:{x:0.35,y:1.75}
      }),getBasicGun({
        damage:0.5,
        health:0.45,
        radius:0.6,
        bound:0.2,
        shotPTime:0.25,
        pos:{x:-0.35,y:1.75}
      })];
      break;
      case 21: // Stalker
      user.guns = [getBasicGun({
          speed:1.5,
          bound:3,
          rotateDistance:Math.PI/120,
          coolTime:0.5,
          pos:{x:0,y:2.45}
      })];
      user.invTime=1.5;
      break;
      case 22: // Ranger
      user.guns = [getBasicGun({
          speed:1.5,
          bound:3,
          rotateDistance:Math.PI/120,
          coolTime:0.5,
          pos:{x:0,y:2.47}
      })];
      break;
      case 23: // Booster
      user.guns = [getBasicGun({
        bound:0.2
      }),getBasicGun({
        damage:0.2,
        bound:0.85,
        life:1.5,
        shotPTime:0.75,
        rotate:Math.PI / 4 * 3,
        pos:{x:0,y:1.45}
      }),getBasicGun({
        damage:0.2,
        bound:0.85,
        life:1.5,
        shotPTime:0.75,
        rotate:-Math.PI / 4 * 3,
        pos:{x:0,y:1.45}
      }),getBasicGun({
        damage:0.2,
        rotate:Math.PI / 6 * 5,
        bound:2.5,
        shotPTime:0.5,
        life:1.5,
        pos:{x:0,y:1.65}
      }),getBasicGun({
        damage:0.2,
        rotate:-Math.PI / 6 * 5,
        bound:2.5,
        shotPTime:0.5,
        life:1.5,
        pos:{x:0,y:1.65}
      })];
      break;
      case 24: // Fighter
      user.guns = [getBasicGun({
          bound:0.2
      }),getBasicGun({
          damage:0.2,
          rotate:Math.PI / 6 * 5,
          bound:2.5,
          shotPTime:0.5,
          life:1.5,
          pos:{x:0,y:1.6}
      }),getBasicGun({
          damage:0.2,
          rotate:-Math.PI / 6 * 5,
          bound:2.5,
          shotPTime:0.5,
          life:1.5,
          pos:{x:0,y:1.6}
      }),getBasicGun({
        damage:0.8,
        rotate:Math.PI / 2,
        coolTime:0.667
      }),getBasicGun({
        damage:0.8,
        rotate:-Math.PI / 2,
        coolTime:0.667
      })];
      break;
      case 25: // Hybrid
      user.guns = [getBasicGun({
        speed:0.75,
        damage:3,
        health:2,
        radius:1.75,
        bound:15,
        bulletBound:0.1,
        coolTime:0.25,
        pos:{x:0,y:1.88}
      }),getDroneGun({
        bulletType:2,
        speed:0.6,
        damage:0.7,
        health:1.4,
        rotate:Math.PI,
        bulletBound:0.4,
        coolTime:0.167,
        autoShot:true,
        bulletLimit:2,
        isOwnCol:true,
        bulletAi:droneAi,
        life:-1,
        pos:{x:0,y:1.6}
      })];
      break;
      case 26: // Manager
      user.guns = [getDroneGun({
        objType:"drone",
        bulletType:2,
        speed:0.6,
        damage:0.7,
        health:2,
        bulletBound:0.4,
        coolTime:0.333,
        autoShot:true,
        bulletLimit:8,
        isOwnCol:true,
        pos:{x:0,y:1.6}
      })];
      user.invTime=1.5;
      break;
      case 27: // MotherShip
      user.guns = [];
      for (let i=-Math.PI/16*15;i<=Math.PI;i+=Math.PI/8){
        user.guns.push(getDroneGun({
          objType:"drone",
          bulletType:2,
          speed:0.6,
          damage:0.7,
          health:2,
          rotate:i,
          radius:0.25,
          bulletBound:0.4,
          coolTime:0.167,
          autoShot:true,
          bulletLimit:2,
          isOwnCol:true,
          pos:{x:0,y:1.1}
        }));
        i+=Math.PI/8;
        user.guns.push(getDroneGun({
          objType:"drone",
          bulletType:2,
          speed:0.6,
          damage:0.7,
          health:2,
          rotate:i,
          radius:0.25,
          bulletBound:0.4,
          coolTime:0.167,
          autoShot:true,
          bulletLimit:2,
          isOwnCol:true,
          bulletAi:droneAi,
          pos:{x:0,y:1.1}
        }));
      }
      break;
      case 28: // Predator
      user.guns = [getBasicGun({
          speed:1.25,
          damage:0.75,
          rotateDistance:Math.PI/120,
          radius:0.7,
          bound:0.03,
          coolTime:0.4,
          pos:{x:0,y:2.23}
      }),getBasicGun({
          speed:1.25,
          damage:0.75,
          radius:0.95,
          rotateDistance:Math.PI/120,
          bound:0.03,
          shotPTime:0.25,
          coolTime:0.4,
          pos:{x:0,y:2.23}
      }),getBasicGun({
          speed:1.25,
          damage:0.75,
          rotateDistance:Math.PI/120,
          radius:1.25,
          bound:0.03,
          shotPTime:0.5,
          coolTime:0.4,
          pos:{x:0,y:2.23}
      })];
      user.event = { // 여기 있는 값들은 모두 "함수" 입니다.
        rightDownEvent:function(u){
          u.variable.setCamera = {
            x: u.x + Math.cos(u.rotate) * 400,
            y: u.y + Math.sin(u.rotate) * 400
          };
        }, // 우클릭을 시작했을 때의 이벤트
        rightEvent:function(u){
          u.owner.camera.x = u.variable.setCamera.x;
          u.owner.camera.y = u.variable.setCamera.y;
        }, // 우클릭을 진행할 때의 이벤트
        rightUpEvent:function(u){}, // 우클릭을 끝낼 때의 이벤트
        killEvent:function(u,deader){ // 오브젝트를 죽였을 때의 이벤트
          u.exp += deader.exp;
        },
        deadEvent:function(u,killer){} // 자신의 오브젝트가 죽었을 때의 이벤트
      }
      user.variable = {setCamera:{x:0,y:0}};
      break;
      case 29: // Sprayer
      user.guns = [getBasicGun({
        damage:0.1,
        radius:0.7,
        bound:0.05,
        shotPTime:0.5,
        pos:{x:0,y:2.25}
      }),getBasicGun({
        damage:0.7,
        rotateDistance:Math.PI/8,
        coolTime:2,
        pos:{x:0,y:1.6}
      })];
      break;
      case 30: // Trapper
      user.guns = [getBasicGun({
        bulletType:0,
        health:2,
        coolTime:0.667,
        radius:0.8,
        life:24,
        isOwnCol:3,
        bulletAi:null,
        pos:{x:0,y:1.2}
      })];
      break;
      case 31: // GunnerTrapper
      user.guns = [getBasicGun({
        bulletType:1,
        damage:0.5,
        radius:0.5,
        shotPTime:1,
        bound:0.4,
        pos:{x:0.3,y:1.5}
      }),getBasicGun({
        bulletType:1,
        damage:0.5,
        radius:0.5,
        shotPTime:1.5,
        bound:0.4,
        pos:{x:-0.3,y:1.5}
      }),getBasicGun({
        bulletType:0,
        health:2,
        rotate:Math.PI,
        coolTime:0.333,
        life:24,
        isOwnCol:3,
        bulletAi:null,
        pos:{x:0,y:1.2}
      })];
      break;
      case 32: // OverTrapper
      user.guns = [getBasicGun({
        bulletType:0,
        health:2,
        coolTime:0.667,
        radius:0.8,
        bound:0.8,
        life:24,
        isOwnCol:3,
        bulletAi:null,
        pos:{x:0,y:1.2}
      }),getDroneGun({
        objType:"drone",
        bulletType:2,
        speed:0.6,
        damage:0.7,
        health:1.4,
        rotate:Math.PI/3*2,
        bulletBound:0.4,
        coolTime:0.167,
        autoShot:true,
        bulletLimit:1,
        isOwnCol:true,
        bulletAi:droneAi,
        pos:{x:0,y:1.6}
      }),getDroneGun({
        objType:"drone",
        bulletType:2,
        speed:0.6,
        damage:0.7,
        health:1.4,
        rotate:-Math.PI/3*2,
        bulletBound:0.4,
        coolTime:0.167,
        autoShot:true,
        bulletLimit:1,
        isOwnCol:true,
        bulletAi:droneAi,
        pos:{x:0,y:1.6}
      })];
      break;
      case 33: // MegaTrapper
      user.guns = [getBasicGun({
        bulletType:0,
        health:3.2,
        damage:1.6,
        coolTime:0.333,
        radius:1.6,
        life:24,
        isOwnCol:3,
        bulletAi:null,
        pos:{x:0,y:1.2}
      })];
      break;
      case 34: // TriTrapper
      user.guns = [getBasicGun({
        bulletType:0,
        health:2,
        coolTime:0.667,
        radius:0.8,
        life:10,
        isOwnCol:3,
        bulletAi:null,
        pos:{x:0,y:1.2}
      }),getBasicGun({
        bulletType:0,
        health:2,
        rotate:Math.PI/3*2,
        coolTime:0.667,
        radius:0.8,
        life:10,
        isOwnCol:3,
        bulletAi:null,
        pos:{x:0,y:1.2}
      }),getBasicGun({
        bulletType:0,
        health:2,
        rotate:-Math.PI/3*2,
        coolTime:0.667,
        radius:0.8,
        life:10,
        isOwnCol:3,
        bulletAi:null,
        pos:{x:0,y:1.2}
      })];
      break;
      case 35: // Smasher
      user.maxStats = [10,10,10,0,0,0,0,10];
      break;
      case 36: // Landmine
      user.maxStats = [10,10,10,0,0,0,0,10];
      user.invTime=13.5;
      break;
      case 37: // AutoGunner
      user.guns = [getBasicGun({
        damage:0.5,
        health:0.45,
        radius:0.6,
        bound:0.2,
        shotPTime:0.5,
        pos:{x:0.67,y:1.3}
      }),getBasicGun({
        damage:0.5,
        health:0.45,
        radius:0.6,
        bound:0.2,
        shotPTime:0.75,
        pos:{x:-0.67,y:1.3}
      }),getBasicGun({
        damage:0.5,
        health:0.45,
        radius:0.6,
        bound:0.2,
        shotPTime:0,
        pos:{x:0.35,y:1.75}
      }),getBasicGun({
        damage:0.5,
        health:0.45,
        radius:0.6,
        bound:0.2,
        shotPTime:0.25,
        pos:{x:-0.35,y:1.75}
      })];
      break;
      case 38: // Auto5
      user.guns = [];
      user.isCanDir=false;
      break;
      case 39: // Auto3
      user.guns = [getAutoGun({
        speed:1.2,
        damage:0.4,
        radius:0.7,
        bound:0.333,
        shotPTime:0.5,
        pos:{x:0,y:1},
        dir:{rotate:0,detect:Math.PI/2,distance:0.5}
      }),getAutoGun({
        speed:1.2,
        damage:0.4,
        radius:0.7,
        rotate:Math.PI/3*2,
        bound:0.333,
        shotPTime:0.5,
        pos:{x:0,y:1},
        dir:{rotate:Math.PI/3*2,detect:Math.PI/2,distance:0.5}
      }),getAutoGun({
        speed:1.2,
        damage:0.4,
        radius:0.7,
        rotate:-Math.PI/3*2,
        bound:0.333,
        shotPTime:0.5,
        pos:{x:0,y:1},
        dir:{rotate:-Math.PI/3*2,detect:Math.PI/2,distance:0.5}
      })];
      user.isCanDir=false;
      break;
      case 40: // SpreadShot
      for (let i=5;i>=1;i--){
        user.guns.push(getBasicGun({
          damage:0.6,
          coolTime:0.5,
          rotate:Math.PI/12*i,
          radius:0.7,
          bound:0.095,
          shotPTime:0.167*i,
          pos:{x:0,y:1.92-i*0.12}
        }));
        user.guns.push(getBasicGun({
          damage:0.6,
          coolTime:0.5,
          rotate:-Math.PI/12*i,
          radius:0.7,
          bound:0.095,
          shotPTime:0.167*i,
          pos:{x:0,y:1.92-i*0.12}
        }));
      }
      user.guns.push(getBasicGun({
        coolTime:0.5,
        bound:0.095
      }));
      break;
      case 41: // Streamliner
      user.guns = [getBasicGun({
          speed:0.9,
          damage:0.2,
          rotateDistance:Math.PI/120,
          radius:0.7,
          bound:0.2,
          shotPTime:0,
          pos:{x:0,y:2.25}
      }),getBasicGun({
          speed:0.9,
          damage:0.2,
          rotateDistance:Math.PI/120,
          radius:0.7,
          bound:0.2,
          shotPTime:0.2,
          pos:{x:0,y:2.25}
      }),getBasicGun({
          speed:0.9,
          damage:0.2,
          rotateDistance:Math.PI/120,
          radius:0.7,
          bound:0.2,
          shotPTime:0.4,
          pos:{x:0,y:2.25}
      }),getBasicGun({
          speed:0.9,
          damage:0.2,
          rotateDistance:Math.PI/120,
          radius:0.7,
          bound:0.2,
          shotPTime:0.6,
          pos:{x:0,y:2.25}
      }),getBasicGun({
          speed:0.9,
          damage:0.2,
          rotateDistance:Math.PI/120,
          radius:0.7,
          bound:0.2,
          shotPTime:0.8,
          pos:{x:0,y:2.25}
      })];
      break;
      case 42: // AutoTrapper
      user.guns = [getBasicGun({
        bulletType:0,
        health:2,
        coolTime:0.667,
        radius:0.8,
        life:24,
        isOwnCol:3,
        bulletAi:null,
        pos:{x:0,y:1.2}
      })];
      break;
      case 43: // BasicDominator
      user.guns = [null,getBasicGun({
        damage:10,
        health:50,
        coolTime:0.333,
        bound:0,
        bulletStance:0
      })];
      user.maxStats = [0,0,0,0,0,0,0,0];
      user.maxHealth = function (){return 6000 + user.level * 2;};
      break;
      case 44: // GunnerDominator
      user.guns = [null,getBasicGun({
        health:5,
        radius:0.3,
        coolTime:3,
        shotPTime:0.666,
        bound:0,
        pos:{x:-0.2,y:1.6}
      }),getBasicGun({
        health:5,
        radius:0.3,
        coolTime:3,
        shotPTime:0.333,
        bound:0,
        pos:{x:0.2,y:1.6}
      }),getBasicGun({
        health:5,
        radius:0.3,
        coolTime:3,
        bound:0
      })];
      user.maxStats = [0,0,0,0,0,0,0,0];
      user.maxHealth = function (){return 6000 + user.level * 2;};
      break;
      case 45: // TrapperDominator
      user.guns = [null];
      for (let i=-Math.PI/4*3;i<=Math.PI;i+=Math.PI/4){
        user.guns.push(getBasicGun({
          bulletType:0,
          damage:3,
          health:20,
          coolTime:0.667,
          radius:0.5,
          rotate:i,
          life:8,
          isOwnCol:3,
          autoShot:true,
          bulletAi:null,
          pos:{x:0,y:1.2}
        }));
      }
      user.maxStats = [0,0,0,0,0,0,0,0];
      user.maxHealth = function (){return 6000 + user.level * 2;};
      break;
      case 46: // BattleShip
      let r = [Math.PI/2,-Math.PI/2];
      let pos = [{x:0.25,y:1.5},{x:-0.25,y:1.5}];
      let battleShipDroneAi = [
        function(b){
          let p = b.owner.owner;
          let target = null;
          if (p && p.mouse.right){
            b.rotate = Math.atan2(b.y-(b.owner.y+p.target.y),b.x-(b.owner.x+p.target.x));
            b.dx += Math.cos(b.rotate) * b.speed;
            b.dy += Math.sin(b.rotate) * b.speed;
            b.goTank = false;
          }
          else if (p && p.mouse.left){
            b.rotate = Math.atan2((b.owner.y+p.target.y)-b.y,(b.owner.x+p.target.x)-b.x);
            b.dx += Math.cos(b.rotate) * b.speed;
            b.dy += Math.sin(b.rotate) * b.speed;
            b.goTank = false;
          }
          else{
            if (b.goEnemy){
              if (b.goEnemy.isDead) b.goEnemy = null;
              else{
                b.rotate = Math.atan2(b.goEnemy.y-b.y,b.goEnemy.x-b.x);
                b.dx+=Math.cos(b.rotate) * b.speed;
                b.dy+=Math.sin(b.rotate) * b.speed;
              }
              b.goTank = false;
            }
            else {
              let dir;
              if (b.goTank)
                dir = b.rotate;
              else
                dir = b.rotate - Math.PI / 2;
              b.goTank = true;
              if (target) b.goEnemy = target;
              b.dx+=Math.cos(dir) * b.speed;
              b.dy+=Math.sin(dir) * b.speed;
              b.rotate=Math.atan2(b.dy,b.dx);
            }
          }
        },
        function (b){
          let target = null;
          if (b.goEnemy){
            if (b.goEnemy.isDead) b.goEnemy = null;
            else{
              b.rotate = Math.atan2(b.goEnemy.y-b.y,b.goEnemy.x-b.x);
              b.dx+=Math.cos(b.rotate) * b.speed;
              b.dy+=Math.sin(b.rotate) * b.speed;
            }
            b.goTank = false;
          }
          else {
            let dir;
            if (b.goTank)
              dir = b.rotate;
            else
              dir = b.rotate - Math.PI / 2;
            b.goTank = true;
            if (target) b.goEnemy = target;
            b.dx+=Math.cos(dir) * b.speed;
            b.dy+=Math.sin(dir) * b.speed;
            b.rotate=Math.atan2(b.dy,b.dx);
          }
        }
      ];
      for (let i=0;i<2;i++){
        for (let j=0;j<2;j++){
          user.guns.push(getBasicGun({
            bulletType:2,
            damage:0.15,
            radius:0.5,
            rotate:r[i],
            bulletBound:3.4,
            life:4,
            bulletAi:battleShipDroneAi[j],
            pos:pos[j]
          }));
        }
      }
      break;
      case 47: // Annihilator
      user.guns = [getBasicGun({
        speed:0.75,
        damage:3,
        health:2,
        radius:2.3,
        bound:17,
        bulletBound:0.1,
        coolTime:0.25,
        pos:{x:0,y:1.88}
      })];
      break;
      case 48: // AutoSmasher
      user.maxStats = [10,10,10,10,10,10,10,10];
      break;
      case 49: // Spike
      user.maxStats = [10,10,10,0,0,0,0,10];
      user.damage = function (){return 28 + user.stats[2] * 4;};
      break;
      case 50: // Factory
      user.guns = [getDroneGun({
        bulletType:3,
        speed:0.334,
        health:4,
        damage:0.7,
        coolTime:0.333,
        bulletLimit:6,
        bulletBound:0.4,
        radius:1.25,
        life:-1,
        pos:{x:0,y:1.6},
        guns:[getBasicGun({
          owner:user,
          speed:0.8,
          health:0.4,
          damage:0.4,
          radius:1.2,
          pos:{x:0,y:2.4}
        })],
        autoShot: true,
        bulletAi:function(b){
          let p = b.owner.owner;
          let target = null;
          if (p && p.mouse.right){
            b.guns[0].autoShot = true;
            b.rotate = Math.atan2(b.y-(b.owner.y+p.target.y),b.x-(b.owner.x+p.target.x));
            let dis = Math.sqrt((b.x-(b.owner.x+p.target.x))*(b.x-(b.owner.x+p.target.x))+(b.y-(b.owner.y+p.target.y))*(b.y-(b.owner.y+p.target.y)));
            if (dis<60){
              b.dx -= Math.cos(b.rotate) * b.speed;
              b.dy -= Math.sin(b.rotate) * b.speed;
            }
            else if (dis<250){
              b.dx += Math.cos(b.rotate + Math.PI / 2) * b.speed;
              b.dy += Math.sin(b.rotate + Math.PI / 2) * b.speed;
            }
            else{
              b.dx += Math.cos(b.rotate) * b.speed;
              b.dy += Math.sin(b.rotate) * b.speed;
            }
            b.goTank = false;
          }
          else if (p && p.mouse.left){
            b.guns[0].autoShot = true;
            b.rotate = Math.atan2((b.owner.y+p.target.y)-b.y,(b.owner.x+p.target.x)-b.x);
            let dis = Math.sqrt(((b.owner.x+p.target.x)-b.x)*((b.owner.x+p.target.x)-b.x)+((b.owner.y+p.target.y)-b.y)*((b.owner.y+p.target.y)-b.y));
            if (dis<150){
              b.dx += Math.cos(b.rotate - Math.PI / 2) * b.speed;
              b.dy += Math.sin(b.rotate - Math.PI / 2) * b.speed;
            }
            else{
              b.dx += Math.cos(b.rotate) * b.speed;
              b.dy += Math.sin(b.rotate) * b.speed;
            }
            b.goTank = false;
          }
          else{
            let dis = Math.sqrt((b.owner.x-b.x)*(b.owner.x-b.x)+(b.owner.y-b.y)*(b.owner.y-b.y));

            if (b.goTank){
              b.autoShot = false;
              if (dis<120) b.goTank = false;
              else{
                b.rotate = Math.atan2(b.owner.y-b.y,b.owner.x-b.x);
                b.dx += Math.cos(b.rotate) * b.speed;
                b.dy += Math.sin(b.rotate) * b.speed;
              }
            }
            else if (b.goEnemy){
              b.guns[0].autoShot = true;
              if (b.goEnemy.isDead){
                b.goEnemy = target;
              }
              else if (dis>417.96){
                b.goEnemy = null;
              }
              else{
                b.rotate = Math.atan2(b.goEnemy.y-b.y,b.goEnemy.x-b.x);
                let dis = Math.sqrt((b.goEnemy.x-b.x)*(b.goEnemy.x-b.x)+(b.goEnemy.y-b.y)*(b.goEnemy.y-b.y));
                if (dis<150){
                  b.dx += Math.cos(b.rotate - Math.PI / 2) * b.speed;
                  b.dy += Math.sin(b.rotate - Math.PI / 2) * b.speed;
                }
                else{
                  b.dx += Math.cos(b.rotate) * b.speed;
                  b.dy += Math.sin(b.rotate) * b.speed;
                }
              }
            }
            else if (dis<util.randomRange(139.32,150)){ // 플레이어 주변에 있으면
              b.guns[0].autoShot = false;
              let dir = Math.atan2(b.owner.y-b.y,b.owner.x-b.x) - Math.PI / 2;
              b.dx += Math.cos(dir) * 0.02;
              b.dy += Math.sin(dir) * 0.02;
              b.rotate=Math.atan2(b.dy,b.dx);
              if (target) b.goEnemy = target;
            }
            else{
              b.goTank = true;
              b.guns[0].autoShot = false;
            }
          }
        }
      })];
      break;
      case 51: // Skimmer
      user.guns = [null,getBasicGun({
          bulletType:4,
          speed:0.6,
          health:3,
          radius:1.75,
          bound:3,
          bulletBound:0.2,
          coolTime:0.25,
          life:4,
          pos:{x:0,y:1.88},
          guns:[getBasicGun({
            owner:user,
            health:0.6,
            damage:0.6,
            radius:5,
            absoluteRadius:true,
            coolTime:3,
            life:0.5,
            bound:1.85,
            autoShot: true,
            pos:{x:0,y:1.8}
          }),getBasicGun({
            owner:user,
            health:0.6,
            damage:0.6,
            radius:5,
            absoluteRadius:true,
            coolTime:3,
            life:0.5,
            bound:1.85,
            rotate:Math.PI,
            autoShot: true,
            pos:{x:0,y:1.8}
          })],
          bulletAi: function(b){
            if (b.variable.isSpawn){
              b.variable.isSpawn = false;
              b.variable.moveRotate = b.rotate;
              b.variable.turnRight = user.owner.mouse.right;
            }
            b.dx+=Math.cos(b.variable.moveRotate) * b.speed;
            b.dy+=Math.sin(b.variable.moveRotate) * b.speed;
            if (b.variable.turnRight) b.rotate-=0.05;
            else b.rotate+=0.05;
          },
          variable: {
            moveRotate: 0,
            turnRight: false,
            isSpawn: true
          }
      })];
      break;
      case 52: // Rocketeer
      user.guns = [null,getBasicGun({
          bulletType:5,
          speed:0.75,
          health:5,
          radius:1.333,
          bound:3,
          bulletBound:0.2,
          coolTime:0.25,
          pos:{x:0,y:1.88},
          guns:[getBasicGun({
            owner:user,
            speed:1.25,
            health:0.6,
            damage:0.3,
            radius:4.5,
            absoluteRadius:true,
            coolTime:5,
            shotPTime:8,
            life:0.5,
            rotate:Math.PI,
            rotateDistance:Math.PI/8,
            autoShot: true,
            pos:{x:0,y:1.8}
          })],
          bulletAi: function(b){
            b.dx+=Math.cos(b.rotate) * b.speed;
            b.dy+=Math.sin(b.rotate) * b.speed;
          }
      })];
      break;
      case 53: // Bumper
      user.maxStats = [8,8,8,0,0,0,8,8];
      user.bound=100;
      user.stance=150;
      user.damage = function (){return 32 + user.stats[2] * 4;};
      user.guns = [getBasicGun({
          speed:0.75,
          damage:0.1,
          health:2,
          radius:1.75,
          bound:50,
          bulletBound:0.1,
          coolTime:0.25,
          pos:{x:0,y:1.88}
      })];
      break;
      case 54: // Dispersion
      user.guns = [];
      for (let i=0;i<7;i++){
        user.guns.push(getBasicGun({
          health:0.7,
          damage:0.5,
          rotateDistance:Math.PI/9,
        }))
      }
      break;
      case 55: // Diffusion
      let list = [];
      for (let i=0;i<8;i++){
        list.push(getBasicGun({
          owner:user,
          damage:0.4,
          life:2,
          pos:{x:0,y:0},
          rotateDistance:Math.PI
        }));
      }
      user.guns = [getBasicGun({
        health:2,
        damage:1.75,
        coolTime:0.25,
        radius:1.75,
        bound:3,
        guns:list,
        event:{
          killEvent:function(u,deader){
            u.owner.exp += deader.exp;
            return true;
          },
          deadEvent:function(u,killer){
            for (let i=0;i<8;i++){
              u.guns[i].autoShot = true;
            }
            return true;
          }
        }
      })];
      break;
      default:
      user.guns = [];
      break;
  }

  setStatsToMaxStats(user);
}
