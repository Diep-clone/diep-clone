export const RGB = function(r,g,b) {

    if (b) {
        this.r = r;
        this.g = g;
        this.b = b;
    } else if (r[0] === "#") {
        let f = function (x) {
            if (x<="9" && x>="0") return x - "0";
            else if (x<="Z" && x>="A") return x - "A" + 11;
            else if (x<="z" && x>="a") return x - "a" + 11;
            return 0;
        }
        this.b = f(r[6])*16 + f(r[5]);
        this.g = f(r[4])*16 + f(r[3]);
        this.r = f(r[1])*16 + f(r[2]);
    }

    this.setRGB = (r, g, b) => {
        if (b){
            this.r = r;
            this.g = g;
            this.b = b;
        } else if (r[0] === "#") {
            let f = function (x) {
                if (x<="9" && x>="0") return x - "0";
                else if (x<="Z" && x>="A") return x - "A" + 11;
                else if (x<="z" && x>="a") return x - "a" + 11;
                return 0;
            }
            this.b = f(r[6])*16 + f(r[5]);
            this.g = f(r[4])*16 + f(r[3]);
            this.r = f(r[1])*16 + f(r[2]);
        }
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

