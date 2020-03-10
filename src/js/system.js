import * as data from './data';

export default class System {
  constructor({ name }) {
    this.name = name;

    this.objectList = [];
    this.uiObjectList = [];

    this.colorList = data.colorList;
    this.expList = data.expList;

    this.tick = 0;
    this.lastTime = Date.now();
    this.isControlRotate = true;

    window.input = data.input;
  }

  insertComma = (number) => number.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');


}