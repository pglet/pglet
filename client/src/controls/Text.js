import React from 'react'

const Text = React.memo(({control}) => {

  console.log(`render Text: ${control.i}`);

  return <div className="text">{control.value}</div>;
})

export default Text