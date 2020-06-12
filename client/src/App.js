import React, { useEffect } from 'react';
//import logo from './logo.svg';
import './pglet.scss';
import { useSelector } from 'react-redux';
import Page from './components/Page'
import LoadingButton from './components/LoadingButton';
import User from './components/User'
import ReconnectingWebSocket from 'reconnecting-websocket';

const App = () => {

  const root = useSelector(state => state.page.controls[0]);

  useEffect(() => {
    console.log("Connecting WebSockets...");
    const conn = new ReconnectingWebSocket(`ws://${document.location.host}/ws`);
        conn.onopen = function (evt) {
          console.log("WebSocket connection opened");
          console.log(evt);
          conn.send("Hello!");
        };
        conn.onclose = function (evt) {
            console.log("WebSocket connection closed");
            console.log(evt);
        };
        conn.onmessage = function (evt) {
          console.log("WebSocket onmessage");
          console.log(evt);
        };
  })

  return (
  <div className="container-fluid">
    <Page control={root} />
    <div className="row">
      <div className="col">
        <LoadingButton />
        <User userId="42" />
      </div>
    </div>
  </div>);
}

export default App;
