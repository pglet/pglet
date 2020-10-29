import React from 'react'

const Button = React.memo(({control}) => {

  console.log(`render Button: ${control.i}`);

  const handleClick = e => {
    console.log(control.event)
  }

  return <button type="button" className="btn btn-primary" onClick={handleClick}>{control.text}</button>;
})

export default Button