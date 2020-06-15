import React from 'react'
import { shallowEqual, useSelector } from 'react-redux'
import NodeList from './NodeList'

const Column = React.memo(({control}) => {

  console.log(`render Column: ${control.i}`);

  const childControls = useSelector(state => control.c.map(childId => state.page.controls[childId]), shallowEqual);

  return <div className="col"><NodeList controls={childControls} /></div>;
})

export default Column