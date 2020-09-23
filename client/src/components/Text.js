import React from 'react'

const Text = React.memo(({control}) => {

  console.log(`render Text: ${control.i}`);

  return <div className="text">{control.text}</div>;
})

export default Text