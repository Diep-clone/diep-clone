import '../styles/index.css';

import System from './system';

const system = new System();

system.loop();

(function(){
    var scale = window.localStorage['no_retina'] ? 1.0 : window.devicePixelRatio;
    var canvas = document.getElementById('canvas');
    var typing = false;
    function preventDefault(e) {
        if (e.preventDefault)
            e.preventDefault();
    }
    function onMouseWheel(e) {
        e = e || window.event;
        preventDefault(e);
        if (!window['input'])
            return;
        window['input']['wheel'](e.wheelDelta / -120 || e.detail || 0);
    }
    function flushInputHooks() {
        if (window["input"] && window["input"]["flushInputHooks"])
            window["input"]["flushInputHooks"]();
    }
    if (/firefox/i.test(navigator.userAgent)) {
        document.addEventListener("DOMMouseScroll", onMouseWheel, false);
    } else {
        document.body['onmousewheel'] = onMouseWheel;
    }
    canvas.onmousemove = function(e) {
        e = e || window.event;
        if (!window['input'])
            return;
        window['input']['mouse'](e.clientX * scale, e.clientY * scale);
    }
    canvas.onmousedown = function(e) {
        flushInputHooks();
        e = e || window.event;
        if (!window['input'])
            return;
        preventDefault(e);
        if (!e.which && e.button !== undefined) {
            e.which = (e.button & 1 ? 1 : (e.button & 2 ? 3 : (e.button & 4 ? 2 : 0)));
        }
        if (e.which >= 1 && e.which <= 3)
            window['input']['keyDown'](e.which);
    }
    canvas.onmouseup = function(e) {
        flushInputHooks();
        e = e || window.event;
        if (!window['input'])
            return;
        preventDefault(e);
        if (!e.which && e.button !== undefined) {
            e.which = (e.button & 1 ? 1 : (e.button & 2 ? 3 : (e.button & 4 ? 2 : 0)));
        }
        if (e.which >= 1 && e.which <= 3)
            window['input']['keyUp'](e.which);
    }
    window.onkeydown = function(e) {
        flushInputHooks();
        e = e || window.event;
        if (!window['input'])
            return;
        if (e.keyCode >= 112 && e.keyCode <= 130 && e.keyCode != 113)
            return;
        window['input']['keyDown'](e.keyCode);
        if (e.keyCode == 9)
            preventDefault(e);
        if (!typing && !e.ctrlKey && !e.metaKey)
            preventDefault(e);
    }
    window.onkeyup = function(e) {
        flushInputHooks();
        e = e || window.event;
        if (!window['input'])
            return;
        if (e.keyCode >= 112 && e.keyCode <= 130 && e.keyCode != 113)
            return;
        window['input']['keyUp'](e.keyCode);
        if (e.keyCode == 9)
            preventDefault(e);
        if (!typing && !e.ctrlKey && !e.metaKey)
            preventDefault(e);
    }
    canvas.onclick = window.onclick = function(e) {
        flushInputHooks();
    }
    canvas.ondragstart = function(e) {
        e = e || window.event;
        preventDefault(e);
    }
    window.onblur = function(e) {
        e = e || window.event;
        if (!window['input'])
            return;
        window['input'].blur();
    }
    canvas.oncontextmenu = function(e) {
        e = e || window.event;
        if (!window['input'] || !window['input']['prevent_right_click'] || window['input']['prevent_right_click']())
            preventDefault(e);
    }
    window['setLoadingStatus'] = function(str) {
        document.getElementById('loading').innerText = str;
    }
    window['setTyping'] = function(v) {
        typing = v;
    }
    window['unscale'] = function(v) {
        return v / scale;
    }
    function onResize() {
        canvas.width = window.innerWidth * scale;
        canvas.height = window.innerHeight * scale;
    }
    window.onerror = function(message, source, lineno, colno, error) {
        window.onerror = null;
        if (error)
            error = error.toString();
        if (typeof Uint8Array == 'undefined')
            return;
        if (typeof WebSocket == 'undefined')
            return;
        if (typeof CanvasRenderingContext2D == 'undefined')
            return;
        if (!document.createElement('canvas').getContext('2d'))
            return;
        if (source != null && source.indexOf("//diep.io") == -1)
            return;
        if (source != null && source.indexOf(".js") == -1)
            return;
        if (source != null) {
            setTimeout(function() {
                alert("The game crashed, refreshing page to recover from error");
                window.onbeforeunload = null;
                window.location.reload(true);
            }, 1000);
        }
        if (error == null && lineno == 0 && colno == 0)
            return;
        if (error != null) {
            if (error.indexOf("renoTransGloRef") != -1)
                return;
        }
        var e = JSON.stringify({
            message: message,
            source: source,
            lineno: lineno,
            colno: colno,
            error: error
        });
        console.log("Uncaught error: " + e);
        var xhr = new XMLHttpRequest();
        xhr.open("POST", "//lb.diep.io/v1/error", true);
        xhr.send(e);
    }
    window.onbeforeunload = function() {
        if (window["input"] && window["input"]["should_prevent_unload"] && window["input"]["should_prevent_unload"]()) {
            return "Are you sure you wanna quit?";
        }
    }
    window.addEventListener("gamepadconnected", function(e) {
        console.log("Gamepad connected: " + e.id);
    });
    window.onresize = onResize;
    onResize();
})();