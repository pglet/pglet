import './pglet.scss';
import React from 'react';
import ReactDOM from 'react-dom';
import { configureStore } from '@reduxjs/toolkit'
import rootReducer from './rootReducer'
import WebSocketProvider from './WebSocket';
import { Provider } from 'react-redux'
import App from './App';
import * as serviceWorker from './serviceWorker';

const store = configureStore({
  reducer: rootReducer
});

console.log(store.getState());

ReactDOM.render(
  <React.StrictMode>
    <Provider store={store}>
      <WebSocketProvider>
        <App />
      </WebSocketProvider>
    </Provider>
  </React.StrictMode>,
  document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
