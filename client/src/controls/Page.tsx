import React, { useEffect, useContext } from 'react'
import { shallowEqual, useSelector } from 'react-redux'
import { ControlsList } from './ControlsList'
import useTitle from '../hooks/useTitle'
import { Stack, IStackProps, IStackTokens } from '@fluentui/react';
import { IControlProps } from './IControlProps'
import { WebSocketContext } from '../WebSocket';

export const Page = React.memo<IControlProps>(({control}) => {

  //console.log(`render page: ${control.i}`);

  const ws = useContext(WebSocketContext);

  // page title
  let title = "Pglet";
  if (control.title) {
    title = control.title
  }
  useTitle(title)

  useEffect(() => {
    // https://danburzo.github.io/react-recipes/recipes/use-effect.html
    // https://codedaily.io/tutorials/72/Creating-a-Reusable-Window-Event-Listener-Hook-with-useEffect-and-useCallback
    const handleWindowClose = (e: any) => {
      console.log('zzaede');
      ws.pageEventFromWeb(control.i, 'close', control.data);
    }
    window.addEventListener("beforeunload", handleWindowClose);
    return () => window.removeEventListener("beforeunload", handleWindowClose);
  }, [control, ws]);

  // stack props
  const stackProps: IStackProps = {
    verticalFill: control.verticalFill ? control.verticalFill : false,
    horizontalAlign: control.horizontalalign ? control.horizontalalign : "start",
    verticalAlign: control.verticalalign ? control.verticalalign : "start",
    styles: {
      root: {
        width: control.width ? control.width : "100%",
        height: control.height !== undefined ? control.height : undefined,
        padding: control.padding ? control.padding : "10px",
        margin: control.margin !== undefined ? control.margin : undefined
      }
    },
  };

  const stackTokens: IStackTokens = {
    childrenGap: control.gap ? control.gap : 10
  }

  const childControls = useSelector((state: any) => control.c.map((childId: string) => state.page.controls[childId]), shallowEqual);

  return <Stack tokens={stackTokens} {...stackProps}>
    <ControlsList controls={childControls} />
  </Stack>
})