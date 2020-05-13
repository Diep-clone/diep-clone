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

export const calByte = {
    getByteLength : function(s) {
		if (s == null || s.length == 0) {
			return 0;
		}
		var size = 0;

		for ( var i = 0; i < s.length; i++) {
			size += this.charByteSize(s.charAt(i));
		}

		return size;
	},	
	cutByteLength : function(s, len) {
		if (s == null || s.length == 0) {
			return 0;
		}
		var size = 0;
		var rIndex = s.length;

		for ( var i = 0; i < s.length; i++) {
			size += this.charByteSize(s.charAt(i));
			if( size == len ) {
				rIndex = i + 1;
				break;
			} else if( size > len ) {
				rIndex = i;
				break;
			}
		}

		return s.substring(0, rIndex);
	},
	charByteSize : function(ch) {
		if (ch == null || ch.length == 0) {
			return 0;
		}

		var charCode = ch.charCodeAt(0);

		if (charCode <= 0x00007F) {
			return 1;
		} else if (charCode <= 0x0007FF) {
			return 2;
		} else if (charCode <= 0x00FFFF) {
			return 3;
		} else {
			return 4;
		}
	}
}; // https://zzznara2.tistory.com/458

export const getPolygonRadius = function (p) {
    if (p === 0) return 1;
    return Math.sqrt(Math.PI / (Math.sin(Math.PI / p) * Math.cos(Math.PI / p) * p));
    //return 2 / (1+Math.cos(Math.PI/p));
}

export const getObjectPoint = function (type) {
    switch (type) {
        case "Trap":
            return -3;
        case "Triangle":
        case "Drone":
            return 3;
        case "Square":
        case "NecroSquare":
        case "Necromancer":
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

export const getTextWidth = function(text, font) {
    var canvas = getTextWidth.canvas || (getTextWidth.canvas = document.createElement("canvas"));
    var context = canvas.getContext("2d");
    context.font = font;
    var metrics = context.measureText(text);
    return metrics.width || 0;
}