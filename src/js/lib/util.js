export const RGB = function(r, g, b) {

    this.r;
    this.g;
    this.b;
    if (b!==undefined) {
        this.r = r;
        this.g = g;
        this.b = b;
    } else if (r[0] === "#") {
        let f = function (x) {
            if (x<="9" && x>="0") return x.charCodeAt(0) - 48;
            else if (x<="Z" && x>="A") return x.charCodeAt(0) - 65 + 10;
            else if (x<="z" && x>="a") return x.charCodeAt(0) - 97 + 10;
            return 0;
        }
        this.b = f(r[5])*16 + f(r[6]);
        this.g = f(r[3])*16 + f(r[4]);
        this.r = f(r[1])*16 + f(r[2]);
    }

    this.setRGB = (r, g, b) => {
        if (b){
            this.r = r;
            this.g = g;
            this.b = b;
        } else if (r[0] === "#") {
            let f = function (x) {
                if (x <= "9" && x >= "0") return x - "0";
                else if (x <= "Z" && x >= "A") return x - "A" + 11;
                else if (x <= "z" && x >= "a") return x - "a" + 11;
                return 0;
            }
            this.b = f(r[6])*16 + f(r[5]);
            this.g = f(r[4])*16 + f(r[3]);
            this.r = f(r[1])*16 + f(r[2]);
        }
    }
    
    this.getRGB = () => {return `rgb(${Math.round(this.r)}, ${Math.round(this.g)}, ${Math.round(this.b)})`;}

    this.getDarkRGB = (per = 0.25) => {
        let r = -this.r*per + this.r;
        let g = -this.g*per + this.g;
        let b = -this.b*per + this.b;

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
        let g = -this.g*per + this.g;
        let b = -this.b*per + this.b;
        
        return new RGB(r,g,b);
    }
}

export const getPolygonRadius = function (p) {
    let d = (90 - 180 / p) * (Math.PI * 2) / 360;
    return Math.sqrt(Math.PI / (Math.sin(d) * Math.cos(d) * p));
}

export const getObjectPoint = function (type) {
    switch (type) {
        case "Triangle":
            return 3;
        case "Square":
        case "NecroSquare":
        case "Necromanser":
        case "Factory":
            return 4;
        case "Pentagon":
            return 5;
        case "Mothership":
            return 20;
        default:
            return 0;
    }
}