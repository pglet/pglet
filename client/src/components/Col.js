import React from 'react'
import { shallowEqual, useSelector } from 'react-redux'
import ControlsList from './ControlsList'

const Col = React.memo(({control}) => {

  console.log(`render Col: ${control.i}`);

  const childControls = useSelector(state => control.c.map(childId => state.page.controls[childId]), shallowEqual);

  return <div className="col"><ControlsList controls={childControls} /></div>;
})

export default Col