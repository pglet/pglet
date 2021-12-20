import React from 'react'
import { shallowEqual, useSelector, useDispatch } from 'react-redux'
import { ControlsList } from './ControlsList'
import useTitle from '../hooks/useTitle'
import { Stack, IStackProps, IStackTokens } from '@fluentui/react';
import { Signin } from './Signin'
import { ISigninProps, IPageProps } from './Control.types'
import { WebSocketContext } from '../WebSocket';
import { changeProps } from '../slices/pageSlice'
import { defaultPixels, getWindowHash, isFalse, isTrue } from './Utils'

export const Page = React.memo<IPageProps>(({ control, pageName, updateTheme }) => {

  const ws = React.useContext(WebSocketContext);
  const dispatch = useDispatch();

  // page title
  let title = `${pageName} - pglet`;
  if (control.title) {
    title = control.title
  }
  useTitle(title)

  const data = {
    fireUpdateHashEvent: true
  }

  function updateHash(hash: string) {
    if (data.fireUpdateHashEvent) {
      const payload: any = {
        i: "page",
        hash: hash
      }
  
      dispatch(changeProps([payload]));
      ws.updateControlProps([payload]);
      ws.pageEventFromWeb("page", 'hashChange', hash);
    }

    data.fireUpdateHashEvent = true;
  }

  React.useEffect(() => {

    // handle resize
    let resizeTimeout:any = null;
    function handleResize() {
      clearTimeout(resizeTimeout);
      resizeTimeout = setTimeout(() => {
        //console.log("window size:", window.innerHeight, window.innerWidth);
        const payload: any = [{
          i: "page",
          width: String(window.innerWidth)
        },
        {
          i: "page",
          height: String(window.innerHeight)
        }];
    
        dispatch(changeProps(payload));
        ws.updateControlProps(payload);
        ws.pageEventFromWeb("page", 'resize', hash);
      }, 250)
    }

    // theme
    updateTheme(control.theme, control.themeprimarycolor, control.themetextcolor, control.themebackgroundcolor);
    
    const hash = getWindowHash();
    const pageHash = control.hash !== undefined ? control.hash : "";

    if (pageHash !== hash) {
      window.location.hash = pageHash ? "#" + pageHash : "";
      data.fireUpdateHashEvent = false;
    }

    const handleHashChange = (e: any) => {
      updateHash(getWindowHash());
    }

    window.addEventListener("hashchange", handleHashChange);
    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener("hashchange", handleHashChange);
      window.removeEventListener("resize", handleResize);
    }
    // eslint-disable-next-line
  }, [control, ws]);

  const childControls = useSelector((state: any) => control.c.map((childId: string) => state.page.controls[childId]), shallowEqual);

  if (isFalse(control.visible)) {
    return null;
  }

  let disabled = isTrue(control.disabled)

  // stack props
  const stackProps: IStackProps = {
    verticalFill: control.verticalfill ? isTrue(control.verticalfill) : undefined,
    horizontalAlign: control.horizontalalign === '' ? undefined : (control.horizontalalign ? control.horizontalalign : "start"),
    verticalAlign: control.verticalalign === '' ? undefined : (control.verticalalign ? control.verticalalign : "start"),
    styles: {
      root: {
        width: control.width ? defaultPixels(control.width) : "100%",
        height: control.height !== undefined ? defaultPixels(control.height) : undefined,
        padding: control.padding ? defaultPixels(control.padding) : "10px",
        margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined
      }
    },
  };

  const stackTokens: IStackTokens = {
    childrenGap: control.gap ? control.gap : 10
  }

  const authProviders = control.signin ? control.signin.split(",").map((s:string) => s.trim().toLowerCase()) : [];
  const signinGroups = isTrue(control.signingroups)

  const handleDismiss = () => {
    const payload: any = {
      i: "page",
      signin: ''
    }

    dispatch(changeProps([payload]));
    ws.updateControlProps([payload]);
    ws.pageEventFromWeb("page", 'dismissSignin', "");
  }

  let signinProps: ISigninProps = {
    signinOptions: {
      gitHubEnabled: authProviders.includes("github") || authProviders.includes("*"),
      gitHubGroupScope: signinGroups,
      azureEnabled: authProviders.includes("azure") || authProviders.includes("*"),
      azureGroupScope: signinGroups,
      googleEnabled: authProviders.includes("google") || authProviders.includes("*"),
      googleGroupScope: signinGroups
    },
    onDismiss: isTrue(control.signinallowdismiss) ? handleDismiss : undefined
  }

  return <>
      <Stack tokens={stackTokens} {...stackProps}>
        <ControlsList controls={childControls} parentDisabled={disabled} />
      </Stack>
      { authProviders.length > 0 &&
        <Signin {...signinProps} />
      }
    </>
})