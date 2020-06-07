import React from 'react';
//import logo from './logo.svg';
import Node from './components/Node'
import './pglet.scss';

import LoadingButton from './components/LoadingButton';
import { useSelector } from 'react-redux';

const App = () => {

  const root = useSelector(state => state.controls[0]);

  return (
  <div>
    <Node control={root} />
    <LoadingButton />
  </div>);
}

export default App;
