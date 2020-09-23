import React from 'react'

const Button = React.memo(({control}) => {

  console.log(`render Button: ${control.i}`);

  return <button type="button" class="btn btn-primary">{control.text}</button>;
})

export default Button