import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { TextField } from 'office-ui-fabric-react';

const Textbox = React.memo<ITextboxProps>(({control}) => {

  console.log(`render Textbox: ${control.i}`);

  const ws = useContext(WebSocketContext);

  const dispatch = useDispatch();
  
  const handleChange = (event: React.FormEvent<HTMLInputElement | HTMLTextAreaElement>, newValue?: string) => {

    const payload = [
      {
        i: control.i,
        "value": newValue
      }
    ];

    dispatch(changeProps(payload));
    ws!.updateControlProps(payload);
  }

  return <TextField value={control.value} onChange={handleChange} />;
})

interface ITextboxProps {
  control: any;
}

export default Textbox