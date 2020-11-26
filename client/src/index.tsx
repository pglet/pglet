import React from 'react';
import ReactDOM from 'react-dom';
import { App } from './App';
import { FluentCustomizations } from '@uifabric/fluent-theme';
import { Customizer, mergeStyles } from 'office-ui-fabric-react';
import * as serviceWorker from './serviceWorker';
import { Provider } from 'react-redux'
import rootReducer from './rootReducer'
import { configureStore } from '@reduxjs/toolkit'
import { WebSocketProvider } from './WebSocket';
import { initializeIcons } from '@uifabric/icons';

initializeIcons();

// Inject some global styles
mergeStyles({
  selectors: {
    ':global(body), :global(html), :global(#root)': {
      margin: 0,
      padding: 0,
      height: '100vh'
    }
  }
});

const store = configureStore({
  reducer: rootReducer
});

ReactDOM.render(
  <Customizer {...FluentCustomizations}>
    <Provider store={store}>
      <WebSocketProvider>
        <App />
      </WebSocketProvider>
    </Provider>
  </Customizer>,
  document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
