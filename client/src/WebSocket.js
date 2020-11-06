import React, { createContext } from 'react'
import { useDispatch } from 'react-redux';
import {
    registerWebClientSuccess,
    registerWebClientError,
    addPageControlsSuccess,
    addPageControlsError,
    changeProps,
    cleanControl,
    removeControl
} from './features/page/pageSlice'
import ReconnectingWebSocket from 'reconnecting-websocket';

const WebSocketContext = createContext(null)

export { WebSocketContext }

export default ({ children }) => {
    let socket;
    let ws;

    const dispatch = useDispatch();

    const registerWebClient = (pageName) => {

        console.log("Call registerWebClient()")
        var msg = {
            action: "registerWebClient",
            payload: {
                pageName: pageName
            }
        }

        socket.send(JSON.stringify(msg));
    }

    const pageEventFromWeb = (eventTarget, eventName, eventData) => {

        console.log("Call pageEventFromWeb()")
        var msg = {
            action: "pageEventFromWeb",
            payload: {
                eventTarget: eventTarget,
                eventName: eventName,
                eventData: eventData
            }
        }

        socket.send(JSON.stringify(msg));
    }

    const updateControlProps = (props) => {

        console.log("Call updateControlProps()")
        var msg = {
            action: "updateControlProps",
            payload: {
                props
            }
        }

        socket.send(JSON.stringify(msg));
    }

    if (!socket) {
        const wsProtocol = document.location.protocol === "https:" ? "wss:" : "ws:";
        socket = new ReconnectingWebSocket(`${wsProtocol}//${document.location.host}/ws`);

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
            var data = JSON.parse(evt.data);
            console.log(data);

            if (data.action === "registerWebClient") {
                if (data.payload.error) {
                    dispatch(registerWebClientError(data.payload.error));
                } else {
                    dispatch(registerWebClientSuccess(data.payload.session));
                }
            } else if (data.action === "addPageControls") {
                if (data.payload.error) {
                    dispatch(addPageControlsError(data.payload.error));
                } else {
                    dispatch(addPageControlsSuccess(data.payload.controls));
                }
            } else if (data.action === "updateControlProps") {
                dispatch(changeProps(data.payload.props));
            } else if (data.action === "cleanControl") {
                dispatch(cleanControl(data.payload));
            } else if (data.action === "removeControl") {
                dispatch(removeControl(data.payload));
            }
        };

        ws = {
            socket: socket,
            registerWebClient,
            pageEventFromWeb,
            updateControlProps
        }
    }

    return (
        <WebSocketContext.Provider value={ws}>
            {children}
        </WebSocketContext.Provider>
    )
}
