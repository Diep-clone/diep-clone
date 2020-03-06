import { RGB } from '../lib/util';

export const objectList = {};
export const objectOrder = [];
export const uiObjectList = [];

export const colorList = [
  new RGB(230,176,138),
  new RGB(228,102,233),
  new RGB(148,102,234),
  new RGB(103,144,234),
  new RGB(234,178,102),
  new RGB(231,103,98),
  new RGB(147,234,103),
  new RGB(103,233,233)
];

export const expList = [
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

export const tick = 0;
export const lastTime = new Date();
export const isControlRotate = true;

export const input = {
  ctrl: false,
  isMouseOverUi: false,

  shot: 0,
  rShot: 0,

  moveRotate: null,
  // moveVector: new Vector(0,0),

  space: false,

  leftMouse: false,
  rightMouse: false,

  w: false,
  a: false,
  s: false,
  d: false,

  num:[false,false,false,false,false,false,false,false],

  o: false,
  k: false,
  l: false,
  e: false,
  c: false,

  autoE: false,
  autoC: false,

  changeTank: false,

  target: {
    x:0,
    y:0
  }
};
