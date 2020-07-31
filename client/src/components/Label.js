import React from 'react'

const Label = React.memo(({control}) => {

  console.log(`render Label: ${control.i}`);

  return <div className="text">{control.text}</div>;
})

export default Label