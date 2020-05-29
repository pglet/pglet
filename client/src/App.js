import React from 'react';
//import logo from './logo.svg';
import Node from './components/Node'
import './App.css';
import './pglet.scss';

import LoadingButton from './components/LoadingButton';

const App = () => {
  return (
  <div>
    <Node id={0} />
    <LoadingButton />
  </div>);
}

export default App;
