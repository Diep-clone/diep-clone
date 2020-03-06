export const RGB = function(r,g,b) {

  this.r = r;
  this.g = g;
  this.b = b;

  this.setRGB = (r, g, b) => {
    this.r = r;
    this.g = g;
    this.b = b;
  }

  this.getRGBValue = () => `rgb(${Math.round(this.r)}, ${Math.round(this.g)}',${Math.round(this.b)})`;
  

  this.getDarkRGB = (per) => {
    if (!per) per = 0.25;

    let r = (0 - this.r) * per + this.r;
    let g = (0 - this.g) * per + this.g;
    let b = (0 - this.b) * per + this.b;

    return new RGB(r,g,b);
  }

  this.getLightRGB = (per) => {
    let r = (255 - this.r) * per + this.r;
    let g = (255 - this.g) * per + this.g;
    let b = (255 - this.b) * per + this.b;

    return new RGB(r,g,b);
  }

  this.getRedRGB = (per) => {
    let r = (255 - this.r) * per + this.r;
    let g = (0 - this.g) * per + this.g;
    let b = (0 - this.b) * per + this.b;
    
    return new RGB(r,g,b);
  }
}

