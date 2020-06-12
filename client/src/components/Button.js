import React from 'react'

const Button = React.memo(({control}) => {

  console.log(`render Button: ${control.i}`);

  return <button>{control.text}</button>;
})

export default Button