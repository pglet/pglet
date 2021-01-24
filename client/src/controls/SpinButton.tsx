import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { SpinButton, ISpinButtonProps } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const MySpinButton = React.memo<IControlProps>(({control, parentDisabled}) => {

  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  const dispatch = useDispatch();

  const handleIncrementDecrement = (value: string, event?: React.MouseEvent<HTMLElement> | React.KeyboardEvent<HTMLElement>) => {
    handleChange(value);
  }

  const handleChange = (value: string) => {
    console.log(value);
    // const payload = [
    //   {
    //     i: control.i,
    //     "value": value
    //   }
    // ];

    // dispatch(changeProps(payload));
    // ws.updateControlProps(payload);
    // ws.pageEventFromWeb(control.i, 'change', String(value))
  }

  //console.log(`render Text: ${control.i}`);

  // https://developer.microsoft.com/en-us/fluentui#/controls/web/references/ifontstyles#IFontStyles

  const props: ISpinButtonProps = {
    //value: control.value ? control.value : undefined,
    label: control.label ? control.label : undefined,
    min: control.min ? parseInt(control.min) : undefined,
    max: control.max ? parseInt(control.max) : undefined,
    step: control.step ? parseInt(control.step) : undefined,
    disabled: disabled,
    styles: {
      root: {
        width: control.width !== undefined ? control.width : undefined,
        height: control.height !== undefined ? control.height : undefined,
        padding: control.padding !== undefined ? control.padding : undefined,
        margin: control.margin !== undefined ? control.margin : undefined
      }
    }
  };

  if (control.icon) {
    props.iconProps = {
      iconName: control.icon
    }
  }

  return <SpinButton {...props}></SpinButton>;
})