'use strict';

exports.randomRange = function(x,y){ // x~y 사이의 랜덤 난수 생성.
  if (x>y){
    let im = x;
    x=y;
    y=im;
  }
  return Math.random() * (y-x) + x;
}

exports.findIndex = function(arr,id){ // 배열에서 id 가 똑같은 인덱스 찾기
  let len = arr.length;

  while (len--){
    if (arr[len].id === id)
    return len;
  }
  return -1;
}

exports.isF = function(value){
  if (typeof value === "function")
    return value();
  else
    return value;
}

exports.floor = function(n,c){
  return Math.floor(n*Math.pow(10,c))/Math.pow(10,c);
}

exports.ccw = function(x1,y1,x2,y2,x3,y3){
  let temp = x1*y2+x2*y3+x3*y1;
  temp -= y1*x2-y2*x3-y3*x1;
  if (temp > 0){
    return 1;
  } else if (temp < 0){
    return -1;
  } else {
    return 0;
  }
}
