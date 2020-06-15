import React from 'react'
import { shallowEqual, useSelector } from 'react-redux'
import NodeList from './NodeList'

const Row = React.memo(({control}) => {

  console.log(`render Row: ${control.i}`);

  const childControls = useSelector(state => control.c.map(childId => state.page.controls[childId]), shallowEqual);

  return <div className="row"><NodeList controls={childControls} /></div>;
})

export default Row