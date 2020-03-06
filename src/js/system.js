import { objectList, objectOrder, uiObjectList, colorList, expList, tick, lastTime, isControlRotate } from './data';

export default class System {
  constructor({ name }) {
    this.name = name;

    this.objectList = objectList;
    this.objectOrder = objectOrder;
    this.uiObjectList = uiObjectList;

    this.colorList = colorList;
    this.expList = expList;

    this.tick = tick;
    this.lastTime = lastTime;
    this.isControlRotate = isControlRotate;
  }

  insertComma = (number) => number.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');


}