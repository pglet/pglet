import React from 'react'

const Textbox = React.memo(({control}) => {

  console.log(`render Textbox: ${control.i}`);

  return <input type="text" value={control.value} />;
})

export default Textbox