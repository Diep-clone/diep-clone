var socket = io.connect();

let system;
function start() {
    if (event.keyCode == 13) {
        document.getElementById('textInput').style.display = 'none';
        document.getElementById('canvas').style.display = 'block';
        
        system = new System(document.getElementById('textInput').value);
        system.drawObject.resize();
    }
}
