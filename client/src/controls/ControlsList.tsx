import React from 'react'
import { MessageBar, MessageBarType } from '@fluentui/react'
import { IControlsListProps } from './Control.types'
import { MyStack } from './Stack'
import { Textbox } from './Textbox'
import { Searchbox } from './Searchbox'
import { MySpinButton } from './SpinButton'
import { MyDropdown } from './Dropdown'
import { MyChoiceGroup } from './ChoiceGroup'
import { MyCheckbox } from './Checkbox'
import { MyToggle } from './Toggle'
import { Progress } from './Progress'
import { MySpinner } from './Spinner'
import { MySlider } from './Slider'
import { Button } from './Button'
import { MyText } from './Text'
import { MyImage } from './Image'
import { MyLink } from './Link'
import { Tabs } from './Tabs'
import { Toolbar } from './Toolbar'
import { MyNav } from './Nav'
import { Grid } from './Grid'
import { Icon } from './Icon'
import { Message } from './Message'
import { MyDialog } from './Dialog'
import { MyPanel } from './Panel'
import { IFrame } from './IFrame'
import { MyVerticalBarChart } from './VerticalBarChart'

export const ControlsList: React.FunctionComponent<IControlsListProps> = ({ controls, parentDisabled }) => {

    const controlTypes: any = {
        'textbox': Textbox,
        'searchbox': Searchbox,
        'icon': Icon,
        'checkbox': MyCheckbox,
        'toggle': MyToggle,
        'dropdown': MyDropdown,
        'choicegroup': MyChoiceGroup,
        'progress': Progress,
        'spinner': MySpinner,
        'slider': MySlider,
        'text': MyText,
        'spinbutton': MySpinButton,
        'link': MyLink,
        'image': MyImage,
        'button': Button,
        'stack': MyStack,
        'tabs': Tabs,
        'toolbar': Toolbar,
        'nav': MyNav,
        'grid': Grid,
        'message': Message,
        'dialog': MyDialog,
        'panel': MyPanel,
        'iframe': IFrame,
        'verticalbarchart': MyVerticalBarChart,
    }

    const renderChild = (control: any) => {
        if (control.visible === "false") {
            return null;
        }
        const ControlType = controlTypes[control.t];
        if (!ControlType) {
            const props = Object.getOwnPropertyNames(control)
                .filter(p => p.length > 1)
                .map(p => `${p}="${control[p]}"`).join(' ');
            return <MessageBar key={control.i} messageBarType={MessageBarType.error} messageBarIconProps={ { iconName: 'WebComponents'} }
                isMultiline><b>Unknown control:</b> {`${control.t} ${props}`}</MessageBar>
        }
        return <ControlType key={control.i} control={control} parentDisabled={parentDisabled} />
    }

    return controls.map(renderChild);
}