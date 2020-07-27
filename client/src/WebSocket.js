import React, { createContext } from 'react'
import { useDispatch } from 'react-redux';
import ReconnectingWebSocket from 'reconnecting-websocket';

const WebSocketContext = createContext(null)

export { WebSocketContext }

export default ({ children }) => {
    let socket;
    let ws;

    const dispatch = useDispatch();

    const registerWebClient = (pageName) => {
        var msg = {
            action: "registerWebClient",
            payload: {
                pageName: pageName
            }
        }

        socket.send(JSON.stringify(msg));
    }

    if (!socket) {
        socket = new ReconnectingWebSocket(`ws://${document.location.host}/ws`);

        socket.onopen = function (evt) {
            console.log("WebSocket connection opened");
            console.log(evt);
        };
        socket.onclose = function (evt) {
            console.log("WebSocket connection closed");
            console.log(evt);
        };

        socket.onmessage = function (evt) {
            console.log("WebSocket onmessage");
            console.log(evt);
        };

        ws = {
            socket: socket,
            registerWebClient
        }
    }

    return (
        <WebSocketContext.Provider value={ws}>
            {children}
        </WebSocketContext.Provider>
    )
}
