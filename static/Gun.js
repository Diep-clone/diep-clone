function Gun(paths, dir, r, g, b){
    'use strict';

    this.paths = paths;
    this.dir = dir;
    this.color = new RGB(r,g,b);
    this.back = 0;

    this.Animate = function (tick){

    }

    this.Draw = function (ctx, camera){

    }
}