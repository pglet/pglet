import React from 'react'
import { useDispatch } from 'react-redux';
import {
    registerWebClientSuccess,
    registerWebClientError,
    appBecomeInactive,
    sessionCrashed,
    addPageControlsSuccess,
    addPageControlsError,
    replacePageControlsSuccess,
    replacePageControlsError,    
    changeProps,
    appendProps,
    cleanControl,
    removeControl,
    pageControlsBatchSuccess,
    pageControlsBatchError
} from './slices/pageSlice'
import ReconnectingWebSocket from 'reconnecting-websocket';
import Cookies from 'universal-cookie';

const cookies = new Cookies();

export interface IWebSocket {
    socket: ReconnectingWebSocket;
    registerWebClient(pageName: string, pageHash: string): void;
    pageEventFromWeb(eventTarget: string, eventName: string, eventData: string): void;
    updateControlProps(props: any): void;
}

const WebSocketContext = React.createContext<IWebSocket>(undefined!)

export { WebSocketContext }

export const WebSocketProvider: React.FC<React.ReactNode> = ({children}) => {
    let socket : ReconnectingWebSocket | null = null;
    let _registeredPageName : string = "";
    let _registeredPageHash : string = "";
    let _subscribed: boolean = false;

    const dispatch = useDispatch();

    if (socket == null) {
        const wsProtocol = document.location.protocol === "https:" ? "wss:" : "ws:";
        socket = new ReconnectingWebSocket(`${wsProtocol}//${document.location.host}/ws`);

        socket.onopen = function () {
            console.log("WebSocket connection opened");
            if (!_subscribed && _registeredPageName !== "") {
                registerWebClient(_registeredPageName, _registeredPageHash);
            }
        };
        socket.onclose = function () {
            console.log("WebSocket connection closed");
            _subscribed = false;
        };

        socket.onmessage = function (evt) {
            //console.log("WebSocket onmessage:", evt.data);
            var data = JSON.parse(evt.data);
            console.log("WebSocket onmessage:", data);

            if (data.action === "registerWebClient") {
                if (data.payload.error) {
                    dispatch(registerWebClientError(data.payload.error));
                } else {
                    dispatch(registerWebClientSuccess({
                        pageName: _registeredPageName,
                        pageHash: _registeredPageHash,
                        session: data.payload.session
                    }));
                }
            } else if (data.action === "appBecomeInactive") {
                dispatch(appBecomeInactive(data.payload));
            } else if (data.action === "sessionCrashed") {
                dispatch(sessionCrashed(data.payload));                
            } else if (data.action === "addPageControls") {
                if (data.payload.error) {
                    dispatch(addPageControlsError(data.payload.error));
                } else {
                    dispatch(addPageControlsSuccess(data.payload));
                }
            } else if (data.action === "replacePageControls") {
                if (data.payload.error) {
                    dispatch(replacePageControlsError(data.payload.error));
                } else {
                    dispatch(replacePageControlsSuccess(data.payload));
                }                
            } else if (data.action === "pageControlsBatch") {
                if (data.payload.error) {
                    dispatch(pageControlsBatchError(data.payload.error));
                } else {
                    dispatch(pageControlsBatchSuccess(data.payload));
                }                
            } else if (data.action === "updateControlProps") {
                dispatch(changeProps(data.payload.props));
            } else if (data.action === "appendControlProps") {
                dispatch(appendProps(data.payload.props));
            } else if (data.action === "cleanControl") {
                dispatch(cleanControl(data.payload));
            } else if (data.action === "removeControl") {
                dispatch(removeControl(data.payload));
            }
        };
    }

    const registerWebClient = (pageName: string, pageHash: string) => {

        console.log("ws.registerWebClient()")
        _registeredPageName = pageName;
        _registeredPageHash = pageHash;
        _subscribed = true;

        var msg = {
            action: "registerWebClient",
            payload: {
                pageName: pageName,
                pageHash: pageHash,
                sessionID: cookies.get(`sid-${pageName}#${pageHash}`)
            }
        }

        console.log(msg);

        socket!.send(JSON.stringify(msg));
    }

    const pageEventFromWeb = (eventTarget: string, eventName: string, eventData: string) => {

        var msg = {
            action: "pageEventFromWeb",
            payload: {
                eventTarget: eventTarget,
                eventName: eventName,
                eventData: eventData
            }
        }
        console.log("ws.pageEventFromWeb()", msg.payload)

        socket!.send(JSON.stringify(msg));
    }

    const updateControlProps = (props: any) => {

        //console.log("ws.updateControlProps()")
        var msg = {
            action: "updateControlProps",
            payload: {
                props
            }
        }

        const msgJson = JSON.stringify(msg);

        //console.log("Call updateControlProps()", msgJson)
        socket!.send(msgJson);
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