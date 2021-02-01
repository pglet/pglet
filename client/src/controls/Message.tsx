import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch, useSelector, shallowEqual } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { MessageBar, MessageBarType, IMessageBarProps, MessageBarButton, IButtonProps } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const Message = React.memo<IControlProps>(({ control }) => {

  //console.log(`render Textbox: ${control.i}`);
  const ws = useContext(WebSocketContext);
  const dispatch = useDispatch();

  const handleDismiss = (actionName: string) => {
    const payload = [
      {
        i: control.i,
        "visible": "false"
      }
    ];

    dispatch(changeProps(payload));
    ws.updateControlProps(payload);
    ws.pageEventFromWeb(control.i, 'dismiss', actionName)
  }

  const buttons = useSelector<any, IButtonProps[]>((state: any) => control.c.map((childId: any) =>
    state.page.controls[childId])
      .filter((oc: any) => oc.t === 'button')
      .map((oc: any) => ({
        key: oc.i,
        text: oc.text ? oc.text : oc.i,
        onClick: () => handleDismiss(oc.action ? oc.action : oc.text ? oc.text : oc.i)
      })), shallowEqual);    

  let barType = 0; // info
  if (control.error === 'true') {
    barType = MessageBarType.error;
  } else if (control.blocked === 'true') {
    barType = MessageBarType.blocked;
  } else if (control.severewarning === 'true') {
    barType = MessageBarType.severeWarning;
  } else if (control.success === 'true') {
    barType = MessageBarType.success;
  } else if (control.warning === 'true') {
    barType = MessageBarType.warning;
  }

  const props: IMessageBarProps = {
    messageBarType: barType,
    isMultiline: control.multiline === 'true',
    truncated: control.truncated === 'true',
    styles: {
      root: {
        width: control.width !== undefined ? control.width : undefined,
        height: control.height !== undefined ? control.height : undefined,
        padding: control.padding !== undefined ? control.padding : undefined,
        margin: control.margin !== undefined ? control.margin : undefined
      },
    }
  };

  if (control.dismiss) {
    props.onDismiss = () => handleDismiss("");
  }

  if (buttons.length > 0) {
    props.actions = <div>
      {buttons.map(buttonProps => (<MessageBarButton {...buttonProps}/>))}
    </div>
  }

  return <MessageBar {...props}>{control.value}</MessageBar>;
})