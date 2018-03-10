const W3CWebSocket = require('websocket').w3cwebsocket;
export const client = new W3CWebSocket('ws://localhost:3000/ws');

client.onopen = () => { 
    console.log('Websocket is established', client.readyState);
}

export const send = function (message, callback) {
    waitForConnection(function () {
        client.send(JSON.stringify(message));
        if (typeof callback !== 'undefined') {
          callback();
        }
    }, 1000);
};

const waitForConnection = function (callback, interval) {
    if (client.readyState === 1) {
        callback();
    } else {
        setTimeout(function () {
            waitForConnection(callback, interval);
        }, interval);
    }
};
