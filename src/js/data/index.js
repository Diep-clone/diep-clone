import { RGB } from '../lib/util';
import { Gun } from './gun';

export const statColor = [
    new RGB(230,176,138),
    new RGB(228,102,233),
    new RGB(148,102,234),
    new RGB(103,144,234),
    new RGB(234,178,102),
    new RGB(231,103,98),
    new RGB(147,234,103),
    new RGB(103,233,233)
];

export const colorList = [
    new RGB("#555555"), // smasher's bases color 0
    new RGB("#999999"), // Barrles's color 1
    new RGB("#00B1DE"), // FFA your's color 2
    new RGB("#00B1DE"), // Blue 3
    new RGB("#F14E54"), // Red 4
    new RGB("#BE7FF5"), // Purple 5
    new RGB("#00F46C"), // Green 6
    //new RGB("#D68163"), //  6 Brown (Old)
    new RGB("#89FF69"), // Shiny 7
    new RGB("#FFE869"), // Square 8
    new RGB("#FC7677"), // Triangle 9
    new RGB("#768DFC"), // Pentagon 10
    new RGB("#FF77DC"), // Crashers 11
    new RGB("#FFE869"), // AC 12
    new RGB("#44FFA0"), // Scoreboard 13
    new RGB("#BBBBBB"), // Maze Walls 14
    new RGB("#F14E54"), // FFA other's color 15
    new RGB("#fcc276"), // Necromanser's Drone color 16
    new RGB("#C0C0C0"), // Fallen 17
];

export const backgroundColor = new RGB("#CDCDCD");
export const minimapBackgroundColor = new RGB("#CDCDCD");
export const minimapBorderColor = new RGB("#555555");

export const colorType = function(type, team){
    switch (team){
        case "shape":
            switch (type){
                case "Square":
                    return 8;
                default:
                    return 2;
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
                case "NecroSquare":
                    return 16;
                default:
                    return 15;
            }
    }
}

export const gunList = {
    "Basic":[
        new Gun([[0,0],[0.42,0],[0.42,1.88],[-0.42,1.88],[-0.42, 0]],0,1),
    ],
    "Twin": [

    ],
    "Triplet": [

    ],
    "TripleShot": [

    ],
    "QuadTank": [

    ],
    "OctoTank": [

    ],
    "Sniper": [

    ],
    "MachineGun": [

    ],
    "Flank Guard": [

    ],
    "Tri-Angle": [

    ],
    "Destroyer": [

    ],
    "Overseer": [

    ],
    "Necromanser":[
        new Gun([[0.43,0],[0.73,1.4],[-0.73,1.4],[-0.43,0]],Math.PI/2,1),
        new Gun([[0.43,0],[0.73,1.4],[-0.73,1.4],[-0.43,0]],-Math.PI/2,1),
    ]
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