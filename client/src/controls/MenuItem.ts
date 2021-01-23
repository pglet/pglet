import {
    IContextualMenuProps,
    ICommandBarItemProps,
    ContextualMenuItemType } from '@fluentui/react';

export function getMenuProps(state: any, parent: any, menuClick: any, itemClick: any): any {
    const itemControls = parent.c.map((childId: any) =>
        state.page.controls[childId])
        .filter((ic: any) => ic.t === 'item' && ic.visible !== "false");
    
    if (itemControls.length === 0) {
      return null
    }

    let menuProps: any = {
      items: Array<any>(),
      onItemClick: menuClick
    };

    for(let i = 0; i < itemControls.length; i++) {
      let item: ICommandBarItemProps = {
        key: itemControls[i].key ? itemControls[i].key : itemControls[i].i,
        text: itemControls[i].text ? itemControls[i].text : (itemControls[i].key ? itemControls[i].key : itemControls[i].i),
        secondaryText: itemControls[i].secondarytext ? itemControls[i].secondarytext : undefined,
        href: itemControls[i].url ? itemControls[i].url : undefined,
        target: itemControls[i].newwindow === 'true' ? '_blank' : undefined,
        disabled: itemControls[i].disabled === 'true' ? true : undefined,
        split: itemControls[i].split === 'true' ? true : undefined
      };
      if (itemControls[i].icon) {
        item.iconProps = {
          iconName: itemControls[i].icon
        }
      }
      if (itemControls[i].divider === 'true') {
        item.itemType = ContextualMenuItemType.Divider;
        item.key = "divider_" + itemControls[i].i;
      }
      if (itemClick != null) {
          item.onClick = itemClick
      }
      const subMenuProps = getMenuProps(state, itemControls[i], menuClick, itemClick);
      if (subMenuProps !== null) {
        item.subMenuProps = subMenuProps
      }
      menuProps.items.push(item)
    }

    return menuProps;
  }