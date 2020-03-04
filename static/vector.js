var Vector = function(x,y){ // 벡터
  this.x = x;
  this.y = y;
}

Vector.prototype.add = function(v){ // 덧셈
  this.x += v.x;
  this.y += v.y;
}

Vector.prototype.sub = function(v){ // 뺄셈
  this.x -= v.x;
  this.y -= v.y;
}

Vector.prototype.mul = function(n){ // 곱셈
  this.x *= n;
  this.y *= n;
}

Vector.prototype.div = function(n){ // 나눗셈
  this.x /= n;
  this.y /= n;
}

Vector.prototype.mag = function(){ // 벡터 크기 구하기
  return Math.sqrt(this.x*this.x + this.y*this.y);
}

Vector.prototype.normalize = function(){ // 벡터 정규화
  var m = this.mag();
  if (m>0){
    this.div(m);
  }
}

Vector.prototype.limit = function(n){ // 벡터 길이 제한
  var m = this.mag();
  if (n<m){
    this.normalize();
    this.mul(n);
  }
}
