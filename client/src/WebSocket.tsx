import React from 'react'
import { useDispatch } from 'react-redux';
import {
    registerWebClientSuccess,
    registerWebClientError,
    addPageControlsSuccess,
    addPageControlsError,
    changeProps,
    cleanControl,
    removeControl
} from './slices/pageSlice'
import ReconnectingWebSocket from 'reconnecting-websocket';

export interface IWebSocket {
    socket: ReconnectingWebSocket;
    registerWebClient(pageName: string): void;
    pageEventFromWeb(eventTarget: string, eventName: string, eventData: string): void;
    updateControlProps(props: any): void;
}

const WebSocketContext = React.createContext<IWebSocket>(undefined!)

export { WebSocketContext }

export const WebSocketProvider: React.FC<React.ReactNode> = ({children}) => {
    let socket : ReconnectingWebSocket | null = null;

    const dispatch = useDispatch();

    if (socket == null) {
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
            console.log("WebSocket onmessage:", evt.data);
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
    }

    const registerWebClient = (pageName: string) => {

        console.log("Call registerWebClient()")
        var msg = {
            action: "registerWebClient",
            payload: {
                pageName: pageName
            }
        }

        socket!.send(JSON.stringify(msg));
    }

    const pageEventFromWeb = (eventTarget: string, eventName: string, eventData: string) => {

        console.log("Call pageEventFromWeb()")
        var msg = {
            action: "pageEventFromWeb",
            payload: {
                eventTarget: eventTarget,
                eventName: eventName,
                eventData: eventData
            }
        }

        socket!.send(JSON.stringify(msg));
    }

    const updateControlProps = (props: any) => {

        console.log("Call updateControlProps()")
        var msg = {
            action: "updateControlProps",
            payload: {
                props
            }
        }

        socket!.send(JSON.stringify(msg));
    }

    const ws: IWebSocket = {
        socket: socket,
        registerWebClient,
        pageEventFromWeb,
        updateControlProps
    }

    return (
        <WebSocketContext.Provider value={ws}>
            {children}
        </WebSocketContext.Provider>
    )
}