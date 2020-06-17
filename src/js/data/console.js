import { RGB } from '../lib/util';

/*
    다이피오의 콘솔에서 바꿀 수 있는 값들을 적어놓는 파일입니다.
    지금은 값들만 있지만, 콘솔 관련 함수들도 여기에 추가해야 합니다.
*/

export let statColor = [
    new RGB(230,176,138),
    new RGB(228,102,233),
    new RGB(148,102,234),
    new RGB(103,144,234),
    new RGB(234,178,102),
    new RGB(231,103,98),
    new RGB(147,234,103),
    new RGB(103,233,233)
];

export let colorList = [
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

export let backgroundColor = new RGB("#CDCDCD");
export let minimapBackgroundColor = new RGB("#CDCDCD");
export let minimapBorderColor = new RGB("#555555");
export let uiOpacity = 0.7;