'use strict'

const util = require('./utility');

const MAX_OBJ = 4;

function Quadtree(x, y, w, h, level){
  this.pos = {x:x, y:y};
  this.siz = {w:w, h:h};
  this.level = level || 1;
  this.objects = [];
  this.nodes = [];

  // console.log("level:"+level+" ("+x+", "+y+") size: "+w+" x "+h);
};

Quadtree.prototype.split = function (){
  var level = this.level + 1,
      w     = this.siz.w / 2,
      h     = this.siz.h / 2;

  /* index
   *  0 1
   *  2 3
   */

  this.nodes[0] = new Quadtree(this.pos.x, this.pos.y, w, h, level);
  this.nodes[1] = new Quadtree(this.pos.x + w, this.pos.y, w, h, level);
  this.nodes[2] = new Quadtree(this.pos.x, this.pos.y + h, w, h, level);
  this.nodes[3] = new Quadtree(this.pos.x + w, this.pos.y + h, w, h, level);
};

Quadtree.prototype.getIndex = function(obj,isArea) {
  var x = this.pos.x + this.siz.w / 2, // 중앙 경계선
      y = this.pos.y + this.siz.h / 2; // 중앙 경계선
  if (isArea){
    if (obj.x >= x && obj.x2 <= x)
  	  return -1;
    if (obj.y >= y && obj.y2 <= y)
  	  return -1;
  }
  else{
    var r = util.isF(obj.radius);

    if (obj.x + r >= x && obj.x - r <= x)
  	  return -1;
    if (obj.y + r >= y && obj.y - r <= y)
  	  return -1;
  }
  return (obj.x > x) + (obj.y > y) * 2;
};

Quadtree.prototype.insert = function(obj) {
  // if (obj.id == 1)
  // console.log("insert (id:" + obj.id + ")("+obj.x + ", " + obj.y + ") at level " + this.level);
  var i = 0,
      index;

  // obj가 this의 관할이 아닐 때
  if (obj.x + obj.radius < this.pos.x || obj.x - obj.radius > this.pos.x + this.siz.w)
    return;
  if (obj.y + obj.radius < this.pos.y || obj.y - obj.radius > this.pos.y + this.siz.h)
    return;

  // 자식 노드가 있을 때
  if (typeof this.nodes[0] !== 'undefined') {
    index = this.getIndex(obj);
    if (index !== -1) {
      this.nodes[index].insert(obj);
      return;
    }
  }

  this.objects.push(obj);

  if (this.objects.length > MAX_OBJ) {
    if (typeof this.nodes[0] === 'undefined') {
      this.split();
    }

    while( i < this.objects.length ) {

      index = this.getIndex( this.objects[ i ] );

      if( index !== -1 ) {
        this.nodes[index].insert( this.objects.splice(i, 1)[0] );
      } else {
        i = i + 1;
      }
    }
  }
};

Quadtree.prototype.retrieve = function(obj,isArea) {
  var index,returnObjects;
  if (isArea){
    index = this.getIndex(obj,true);
    returnObjects = this.objects;
  }
  else{
  	index = this.getIndex(obj);
		returnObjects = this.objects.map((o)=>{
      if (!o.isDead && obj.owner!==o.owner || (obj.owner===o.owner && obj.isOwnCol>0 && o.isOwnCol>0)) return o;
    }).filter(function(f) { return f; });
  }

  if( typeof this.nodes[0] !== 'undefined' ) {
    if( index !== -1 ) {
      returnObjects = returnObjects.concat( this.nodes[index].retrieve(obj,isArea) );
    } else {
      for( var i=0; i < this.nodes.length; i=i+1 ) {
        returnObjects = returnObjects.concat( this.nodes[i].retrieve(obj,isArea) );
      }
    }
  }

	return returnObjects;
};

Quadtree.prototype.clear = function() {
	this.objects = [];

	for( var i=0; i < this.nodes.length; i+=1 ) {
		if( typeof this.nodes[i] !== 'undefined' ) {
			this.nodes[i].clear();
		}
	}

	this.nodes = [];
};

module.exports = Quadtree;
/*
 * 노드에 오브젝트들이 5개 이상이 된다면, 자식 노드 4개를 생성한다.
 * 노드의 크기는 노드의 "중심" 좌표 (pos), 노드의 너비와 높이 (siz) 로 한다.
 */
