import { ICommandBarItemProps, ContextualMenuItemType, Theme} from '@fluentui/react';
import { IWebSocket } from '../WebSocket';
import { getThemeColor, getId } from './Utils'

export function getMenuProps(state: any, parent: any, parentDisabled: boolean, ws: IWebSocket, theme: Theme, inverted: boolean): any {
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
            split: itemControls[i].split === 'true',
            checked: itemControls[i].checked === 'true',
            className: getId(itemControls[i].i)
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

        if (inverted) {
            const darkerBkg = "linear-gradient(rgba(0,0,0,0.1),rgba(0,0,0,0.1))";
            const menuColor = getThemeColor(theme, theme.isInverted ? "neutralPrimary" : "neutralLight");
            const whiteColor = getThemeColor(theme, theme.isInverted ? "black" : "white");
    
            item.buttonStyles = {
                root: {
                    backgroundColor: "inherit",
                    color: menuColor
                },
                rootHovered: {
                    //backgroundColor: 'red',
                    background: darkerBkg,
                    color: whiteColor
                },
                rootPressed: {
                    background: darkerBkg,
                    color: whiteColor
                },
                rootExpanded: {
                    background: darkerBkg,
                    color: whiteColor
                },
                rootExpandedHovered: {
                    background: darkerBkg,
                    color: whiteColor
                },
                rootChecked: {
                    backgroundColor: getThemeColor(theme, "white") + "!important"
                },
    
                icon: {
                    color: menuColor,
                },
                iconHovered: {
                    color: whiteColor
                },
                iconPressed: {
                    color: whiteColor
                },
                iconExpanded: {
                    color: whiteColor + "!important"
                },
    
                menuIcon: {
                    color: menuColor,
                    background: "transparent"
                },
                menuIconHovered: {
                    color: whiteColor,
                    ".ms-Button-icon": {
                        color: whiteColor,
                    },                 
                },
                menuIconPressed: {
                    color: whiteColor
                },
                menuIconExpanded: {
                    color: whiteColor + "!important"
                },            
                menuIconExpandedHovered: {
                    color: whiteColor,
                    background: "transparent"
                },
    
                splitButtonDivider: {
                    background: darkerBkg
                },
                splitButtonMenuButton: {
                    background: "inherit",
                    ":hover": {
                        color: whiteColor,
                        background: darkerBkg,
                    },
                    ":hover .ms-Button-icon": {
                        color: whiteColor,
                    },                
                    ":active": {
                        background: darkerBkg,
                    },
                    ":active .ms-Button-icon": {
                        color: whiteColor,
                    },                 
                },
                splitButtonMenuIcon: {
                    color: menuColor
                },
                splitButtonMenuButtonExpanded: {
                    background: darkerBkg,
                    ":hover": {
                        background: darkerBkg,
                    },
                    ":active": {
                        background: darkerBkg,
                    },
                    ".ms-Button-icon": {
                        color: whiteColor,
                    },                
                },
            }
        }

        item.onClick = () => {
            ws.pageEventFromWeb(itemControls[i].i, 'click', itemControls[i].data)
        }

        const subMenuProps = getMenuProps(state, itemControls[i], disabled, ws, theme, inverted);
        if (subMenuProps !== null) {
            item.subMenuProps = subMenuProps
        }
        menuProps.items.push(item)
    }

    return menuProps;
}