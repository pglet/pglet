import { ICommandBarItemProps, ContextualMenuItemType, Theme} from '@fluentui/react';
import { IWebSocket } from '../WebSocket';
import { getThemeColor } from './Utils'

export function getMenuProps(state: any, parent: any, parentDisabled: boolean, ws: IWebSocket, theme: Theme): any {
    const childControls = parent.children !== undefined ? parent.children
    : parent.c.map((childId: any) => state.page.controls[childId]);

    const itemControls = childControls.filter((ic: any) => ic.t === 'item' && ic.visible !== "false");

    if (itemControls.length === 0) {
        return null
    }

    let menuProps: any = {
        items: Array<any>()
    };

    for (let i = 0; i < itemControls.length; i++) {
        let disabled = (itemControls[i].disabled === 'true') || parentDisabled;

        let item: ICommandBarItemProps = {
            key: itemControls[i].i,
            text: itemControls[i].text ? itemControls[i].text : itemControls[i].i,
            secondaryText: itemControls[i].secondarytext ? itemControls[i].secondarytext : undefined,
            href: itemControls[i].url ? itemControls[i].url : undefined,
            title: itemControls[i].title ? itemControls[i].title : undefined,
            target: itemControls[i].newwindow === 'true' ? '_blank' : undefined,
            disabled: disabled,
            split: itemControls[i].split === 'true' ? true : undefined
        };
        if (itemControls[i].icon) {
            item.iconProps = {
                iconName: itemControls[i].icon
            }
            item.iconOnly = itemControls[i].icononly === 'true' ? true : false;

            if (itemControls[i].iconcolor !== undefined) {
                item.iconProps!.styles = {
                    root: {
                        color: getThemeColor(theme, itemControls[i].iconcolor)
                    }
                }
                item.buttonStyles = {
                    rootHovered: {
                        '.ms-Button-icon': {
                            color: getThemeColor(theme, itemControls[i].iconcolor)
                        }
                    },
                    rootPressed: {
                        '.ms-Button-icon': {
                            color: getThemeColor(theme, itemControls[i].iconcolor)
                        }
                    }
                }
            }
        }
        if (itemControls[i].divider === 'true') {
            item.itemType = ContextualMenuItemType.Divider;
            item.key = "divider_" + itemControls[i].i;
        }

        item.onClick = () => {
            ws.pageEventFromWeb(itemControls[i].i, 'click', itemControls[i].data)
        }

        const subMenuProps = getMenuProps(state, itemControls[i], disabled, ws, theme);
        if (subMenuProps !== null) {
            item.subMenuProps = subMenuProps
        }
        menuProps.items.push(item)
    }

    return menuProps;
}