import '../styles/index.css';

import System from './system';

const system = new System();

system.loop();

/*
const run = (e) => {
    console.log(e ,15)
    if (e.keyCode == 13) { // 13 = Enter Key 
        document.getElementById('name-input').style.display = 'none';
        document.getElementById('canvas').style.display = 'block';

        const name = document.getElementById('name-input').value;
        const system = new System(name);

        // system.drawObject.resize();
    }
};

document.getElementById('name-input').onkeypress = run;
*/