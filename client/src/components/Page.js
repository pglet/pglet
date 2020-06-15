import React from 'react'
import { shallowEqual, useSelector } from 'react-redux'
import NodeList from './NodeList'

const Page = React.memo(({control}) => {

  console.log(`render page: ${control.i}`);

  const childControls = useSelector(state => control.c.map(childId => state.page.controls[childId]), shallowEqual);

  return <NodeList controls={childControls} />
})

export default Page