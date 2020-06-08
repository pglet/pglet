import React from 'react';
//import logo from './logo.svg';
import './pglet.scss';
import { useSelector } from 'react-redux';
import Node from './components/Node'
import LoadingButton from './components/LoadingButton';
import User from './components/User'

const App = () => {

  const root = useSelector(state => state.page.controls[0]);

  return (
  <div>
    <Node control={root} />
    <LoadingButton />
    <User userId="1" />
  </div>);
}

export default App;
