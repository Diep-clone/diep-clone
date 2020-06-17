import { Gun } from '../data/gun.js';

export const RGB = function(r, g, b) {
    /*

    */
    this.r;
    this.g;
    this.b;

    this.setRGB = (r, g, b) => {
        if (b){
            this.r = r;
            this.g = g;
            this.b = b;
        } else if (r[0] === "#") { // hex to RGB
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

    this.setRGB(r, g, b);
    
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

export const calByte = { // string's byte calculation functions
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
    if (p === 0) return 1; // circle
    return Math.sqrt(Math.PI / (Math.sin(Math.PI / p) * Math.cos(Math.PI / p) * p)); // polygon
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


export const colorType = function(type, team){
    switch (team){            
        case "shape":
            switch (type){
                case "Square":
                    return 8;
                default:
                    return 7;
            }
        case "red":
            return 4;
        case "blue":
            return 3;
        case "purple":
            return 5;
        case "green":
            return 6;
        case system.playerSetting.team:
            switch (type){
                case "NecroSquare":
                    return 16;
                default:
                    return 2;
            }
        default:
            switch (type){
                case "AutoGun":
                    return 1;
                case "NecroSquare":
                    return 16;
                default:
                    return 15;
            }
    }
}

export const gunList = { // Returns the width of a given text typed.
    "Basic":[
        new Gun([[0.42,0],[0.42,1.88],[-0.42,1.88],[-0.42, 0]]),
    ],
    "Twin": [
        new Gun([[0.9,0],[0.9,1.88],[0.1,1.88],[0.1, 0]]),
        new Gun([[-0.1,0],[-0.1,1.88],[-0.9,1.88],[-0.9, 0]])
    ],
    "Triplet": [
        new Gun([[0.9,0],[0.9,1.6],[0.1,1.6],[0.1, 0]]),
        new Gun([[-0.1,0],[-0.1,1.6],[-0.9,1.6],[-0.9, 0]]),
        new Gun([[0.4,0],[0.4,1.88],[-0.4,1.88],[-0.4, 0]])
    ],
    "TripleShot": [
        new Gun([[0.4,0],[0.4,1.88],[-0.4,1.88],[-0.4, 0]]),
        new Gun([[0.4,0],[0.4,1.88],[-0.4,1.88],[-0.4, 0]],-Math.PI / 4),
        new Gun([[0.4,0],[0.4,1.88],[-0.4,1.88],[-0.4, 0]],Math.PI / 4)   
    ],
    "QuadTank": [

    ],
    "OctoTank": [
        new Gun([[0.4,0],[0.4,1.88],[-0.4,1.88],[-0.4, 0]]),
        new Gun([[0.4,0],[0.4,1.88],[-0.4,1.88],[-0.4, 0]],-Math.PI / 2),
        new Gun([[0.4,0],[0.4,1.88],[-0.4,1.88],[-0.4, 0]],Math.PI / 2),
        new Gun([[0.4,0],[0.4,1.88],[-0.4,1.88],[-0.4, 0]],Math.PI),
        new Gun([[0.4,0],[0.4,1.88],[-0.4,1.88],[-0.4, 0]],Math.PI / 4),
        new Gun([[0.4,0],[0.4,1.88],[-0.4,1.88],[-0.4, 0]],-Math.PI / 4 * 3),
        new Gun([[0.4,0],[0.4,1.88],[-0.4,1.88],[-0.4, 0]],Math.PI / 4 * 3),
        new Gun([[0.4,0],[0.4,1.88],[-0.4,1.88],[-0.4, 0]],-Math.PI / 4)
    ],
    "Sniper": [
        new Gun([[0.4,0],[0.4,2.2],[-0.4,2.2],[-0.4, 0]])
    ],
    "MachineGun": [
        new Gun([[0.4,0],[0.8,1.88],[-0.8,1.88],[-0.4, 0]])
    ],
    "FlankGuard": [
        new Gun([[0.4,0],[0.4,1.88],[-0.4,1.88],[-0.4, 0]]),
        new Gun([[0.4,0],[0.4,1.6],[-0.4,1.6],[-0.4, 0]],Math.PI)
    ],
    "Tri-Angle": [

    ],
    "Destroyer": [
        new Gun([[0.7,0],[0.7,1.88],[-0.7,1.88],[-0.7, 0]])
    ],
    "Overseer": [

    ],
    "Overlord": [
        new Gun([[0.43,0],[0.73,1.4],[-0.73,1.4],[-0.43,0]]),
        new Gun([[0.43,0],[0.73,1.4],[-0.73,1.4],[-0.43,0]],Math.PI/2),
        new Gun([[0.43,0],[0.73,1.4],[-0.73,1.4],[-0.43,0]],Math.PI),
        new Gun([[0.43,0],[0.73,1.4],[-0.73,1.4],[-0.43,0]],-Math.PI/2),
    ],
    "TwinFlank": [

    ],
    "PentaShot": [

    ],
    "Assasin": [

    ],
    "ArenaCloser": [

    ],
    "Necromancer":[
        new Gun([[0.43,0],[0.73,1.4],[-0.73,1.4],[-0.43,0]],Math.PI/2),
        new Gun([[0.43,0],[0.73,1.4],[-0.73,1.4],[-0.43,0]],-Math.PI/2),
    ],
    "TripleTwin":[

    ],
    "Hunter": [
        new Gun([[0.38,0],[0.38,2.23],[-0.38,2.23],[-0.54, 0]]),
        new Gun([[0.54,0],[0.54,1.91],[-0.54,1.91],[-0.54, 0]])
    ],
    "Gunner": [

    ],
    "Stalker": [
        new Gun([[0.72,0],[0.42,2.4],[-0.42,2.4],[-0.72, 0]])
    ],
    "Ranger": [

    ],
    "Booster": [

    ],
    "Fighter": [

    ],
    "Hybrid": [

    ],
    "Manager": [

    ],
    "Mothership":[

    ],
    "Predator": [
        new Gun([[0.42,0],[0.42,2.2],[-0.42,2.2],[-0.42, 0]]),
        new Gun([[0.565,0],[0.565,1.9],[-0.565,1.9],[-0.565, 0]]),
        new Gun([[0.713,0],[0.713,1.59],[-0.713,1.59],[-0.713, 0]])
    ],
    "Sprayer": [

    ],
    "Trapper": [
        new Gun([[0.42,0],[0.42,1.193],[0.73,1.578],[-0.73,1.578],[-0.42,1.193],[0.42,1.193],[-0.42,1.193],[-0.42,0]])
    ],
    "GunnerTrapper": [

    ],
    "OverTrapper": [

    ],
    "MegaTrapper": [

    ],
    "TriTrapper": [
        new Gun([[0.42,0],[0.42,1.193],[0.73,1.578],[-0.73,1.578],[-0.42,1.193],[0.42,1.193],[-0.42,1.193],[-0.42,0]]),
        new Gun([[0.42,0],[0.42,1.193],[0.73,1.578],[-0.73,1.578],[-0.42,1.193],[0.42,1.193],[-0.42,1.193],[-0.42,0]],Math.PI / 3 * 2),
        new Gun([[0.42,0],[0.42,1.193],[0.73,1.578],[-0.73,1.578],[-0.42,1.193],[0.42,1.193],[-0.42,1.193],[-0.42,0]],-Math.PI / 3 * 2)   
    ],
    "Smasher": [

    ],
    "Landmine": [

    ],
    "AutoGunner": [

    ],
    "Auto5": [

    ],
    "Auto3": [

    ],
    "SpreadShot": [

    ],
    "Streamliner": [

    ],
    "AutoTrapper": [

    ],
    "BasicDominator": [

    ],
    "GunnerDominator": [

    ],
    "TrapperDominator": [

    ],
    "Battleship": [
        new Gun([[0.9,0],[0.68,1.5],[0.1,1.5],[-0.16,0]],Math.PI / 2),
        new Gun([[-0.9,0],[-0.68,1.5],[-0.1,1.5],[0.16,0]],Math.PI / 2),
        new Gun([[0.9,0],[0.68,1.5],[0.1,1.5],[-0.16,0]],-Math.PI / 2),
        new Gun([[-0.9,0],[-0.68,1.5],[-0.1,1.5],[0.16,0]],-Math.PI / 2),
    ],
    "Annihilator": [

    ],
    "AutoSmasher": [

    ],
    "Spike": [
        new Gun(spikeGun(), 0, 0, Math.PI / 40)
    ],
    "Factory": [

    ],
    "Skimmer": [
        new Gun([[0.3,0],[0.59,1.84],[-0.59,1.84],[-0.3, 0]]),
        new Gun([[0.71,0],[0.71,1.59],[-0.71,1.59],[-0.71, 0]])
    ],
    "Rocketeer": [
        new Gun([[0.3,0],[0.59,1.84],[-0.59,1.84],[-0.3, 0]]),
        new Gun([[0.92,0],[0.52,1.59],[-0.52,1.59],[-0.92, 0]])
    ],
    "SkimmerBullet": [
        new Gun([[7,0],[7,1.4],[-7,1.4],[-7,0]], 0, -1, 0, false, true),
        new Gun([[7,0],[7,1.4],[-7,1.4],[-7,0]], Math.PI, -1, 0, false, true)
    ],
    "RocketeerBullet": [
        new Gun([[6,0],[10,1.5],[-10,1.5],[-6,0]],Math.PI, -1, 0, false, true)
    ],
}

function spikeGun() { // spike's guns
    let list = [];
    let dir = -Math.PI / 12 * 11;
    for (let i=0;i<24;i++){
        list.push([Math.cos(dir)*0.9,Math.sin(dir)*0.9]);
        dir += Math.PI / 12;
        list.push([Math.cos(dir)*1.3,Math.sin(dir)*1.3]);
        dir += Math.PI / 12;
    }
    return list;
}

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