import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { changeProps } from '../slices/pageSlice'
import { useDispatch, shallowEqual, useSelector } from 'react-redux'
import { Nav, INavProps, INavLink, mergeStyles } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const MyNav = React.memo<IControlProps>(({ control, parentDisabled }) => {

  //console.log(`render Button: ${control.i}`);

  const dispatch = useDispatch();

  const ws = useContext(WebSocketContext);

  const navItems = useSelector<any, any>((state: any) => {

    function getNavLinks(parent: any): any {
      const itemControls = (parent.children !== undefined ? parent.children : parent.c.map((childId: any) => state.page.controls[childId]))
        .filter((ic: any) => ic.t === 'item' && ic.visible !== "false");

      if (itemControls.length === 0) {
        return []
      }

      let items = [];

      for (let i = 0; i < itemControls.length; i++) {
        let disabled = (itemControls[i].disabled === 'true') || parentDisabled;

        let item: any = {
          id: itemControls[i].i,
          key: itemControls[i].key ? itemControls[i].key : itemControls[i].text,
          name: itemControls[i].text ? itemControls[i].text : itemControls[i].key,
          url: itemControls[i].url ? itemControls[i].url : undefined,
          target: itemControls[i].newwindow === 'true' ? '_blank' : undefined,
          disabled: disabled,
          isExpanded: itemControls[i].expanded === "true",
          collapseByDefault: itemControls[i].expanded === 'false', // groups only
        };

        item.links = getNavLinks(itemControls[i]);

        if (itemControls[i].icon !== undefined) {
          item.iconProps = {
            iconName: itemControls[i].icon
          }

          if (itemControls[i].iconcolor !== undefined && !disabled) {
            item.iconProps.className = mergeStyles({
              color: itemControls[i].iconcolor + '!important'
            });
          }
        }

        items.push(item)
      }

      return items;
    }
    return getNavLinks(control)
  }, shallowEqual)

  let navProps: INavProps = {
    groups: navItems,
    styles: {
      root: {
        width: control.width !== undefined ? control.width : undefined,
        height: control.height !== undefined ? control.height : undefined,
        padding: control.padding !== undefined ? control.padding : undefined,
        margin: control.margin !== undefined ? control.margin : undefined
      }
    }
  };

  const handleExpandLink = (ev?: React.MouseEvent<HTMLElement>, item?: INavLink) => {
    //console.log("EXPAND:", item!.isExpanded!.toString())

    const selectedKey = item!.key as string
    const eventName = item!.isExpanded ? "collapse" : "expand";

    const payload = [
      {
        i: item!.id,
        "expanded": (!item!.isExpanded!).toString()
      }
    ];

    dispatch(changeProps(payload));
    ws.updateControlProps(payload);
    ws.pageEventFromWeb(control.i, eventName, selectedKey)
  }

  const handleLinkClick = (ev?: React.MouseEvent<HTMLElement>, item?: INavLink) => {

    //console.log("ITEM:", item!.links!.length);

    const selectedKey = item!.key as string
    if (selectedKey === undefined) {
      return
    }

    const payload = [
      {
        i: control.i,
        "value": selectedKey
      }
    ];

    dispatch(changeProps(payload));
    ws.updateControlProps(payload);
    ws.pageEventFromWeb(control.i, 'change', selectedKey)
  }

  if (control.value) {
    navProps.selectedKey = control.value;
  }

  return <Nav {...navProps} onLinkClick={handleLinkClick} onLinkExpandClick={handleExpandLink} />;
})