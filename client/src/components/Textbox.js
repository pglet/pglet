import React from 'react'
import { useDispatch } from 'react-redux'
import { changeProps } from '../features/page/pageSlice'

const Textbox = React.memo(({control}) => {

  console.log(`render Textbox: ${control.i}`);

  const dispatch = useDispatch();
  
  const handleChange = e => {
    dispatch(changeProps({
        nodeId: control.i,
        newProps: {
            "value": e.target.value
        }
    }));
  }

  return <input type="text" value={control.value} onChange={handleChange} />;
})

export default Textbox