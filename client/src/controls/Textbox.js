import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { TextField } from 'office-ui-fabric-react';

const Textbox = React.memo(({control}) => {

  console.log(`render Textbox: ${control.i}`);

  const ws = useContext(WebSocketContext);

  const dispatch = useDispatch();
  
  const handleChange = e => {

    const payload = [
      {
        i: control.i,
        "value": e.target.value
      }
    ];

    dispatch(changeProps(payload));
    ws.updateControlProps(payload);
  }

  return <TextField value={control.value} onChange={handleChange} />;
})

export default Textbox