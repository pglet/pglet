import React from 'react'
import { shallowEqual, useSelector } from 'react-redux'
import ControlsList from './ControlsList'
import useTitle from '../hooks/useTitle'

const Page = React.memo(({ control }) => {

  console.log(`render page: ${control.i}`);

  var title = "Pglet"
  if (control.title) {
    title = control.title
  }

  useTitle(title)

  const childControls = useSelector(state => control.c.map(childId => state.page.controls[childId]), shallowEqual);

  return <ControlsList controls={childControls} />
})

export default Page