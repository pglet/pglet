import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch, useSelector, shallowEqual } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { MessageBar, MessageBarType, IMessageBarProps, MessageBarButton, IButtonProps, useTheme } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { getThemeColor, defaultPixels } from './Utils'

export const Message = React.memo<IControlProps>(({ control }) => {

  const ws = useContext(WebSocketContext);
  const dispatch = useDispatch();
  const theme = useTheme();

  const handleDismiss = (actionName: string) => {

    const val = "false"

    let payload: any = {}
    if (control.f) {
      // binding redirect
      const p = control.f.split('|')
      payload["i"] = p[0]
      payload[p[1]] = val
    } else {
      // unbound control
      payload["i"] = control.i
      payload["visible"] = val
    }

    dispatch(changeProps([payload]));
    ws.updateControlProps([payload]);
    ws.pageEventFromWeb(control.i, 'dismiss', control.data ? `${control.data}|${actionName}` : actionName)
  }

  const buttons = useSelector<any, IButtonProps[]>((state: any) =>
    (control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId]))
      .filter((oc: any) => oc.t === 'button')
      .map((oc: any) => ({
        key: oc.i,
        text: oc.text ? oc.text : oc.action,
        onClick: () => handleDismiss(oc.action ? oc.action : oc.text)
      })), shallowEqual);

  let barType = MessageBarType.info; // info
  switch (control.type ? control.type.toLowerCase() : '') {
    case 'error': barType = MessageBarType.error; break;
    case 'blocked': barType = MessageBarType.blocked; break;
    case 'severewarning': barType = MessageBarType.severeWarning; break;
    case 'success': barType = MessageBarType.success; break;
    case 'warning': barType = MessageBarType.warning; break;
  }

  const props: IMessageBarProps = {
    messageBarType: barType,
    isMultiline: control.multiline === 'true',
    truncated: control.truncated === 'true',
    styles: {
      root: {
        width: control.width !== undefined ? defaultPixels(control.width) : undefined,
        height: control.height !== undefined ? defaultPixels(control.height) : undefined,
        padding: control.padding !== undefined ? defaultPixels(control.padding) : undefined,
        margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined
      },
    }
  };

  if (control.icon) {
    props.messageBarIconProps = {
      iconName: control.icon
    }

    if (control.iconcolor) {
      props.messageBarIconProps!.styles = {
        root: {
          color: getThemeColor(theme, control.iconcolor)
        }
      }
    }
  }

  if (control.dismiss) {
    props.onDismiss = () => handleDismiss("");

    if (control.dismissicon) {
      props.dismissIconProps = {
        iconName: control.dismissicon
      }

      if (control.dismissiconcolor) {
        props.dismissIconProps!.styles = {
          root: {
            color: getThemeColor(theme, control.dismissiconcolor) + "!important"
          }
        }
      }
    }
  }

  if (buttons.length > 0) {
    props.actions = <div>
      {buttons.map(buttonProps => (<MessageBarButton {...buttonProps} />))}
    </div>
  }

  return <MessageBar {...props}>{control.value}</MessageBar>;
})