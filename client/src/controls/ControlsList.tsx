import React from 'react'
import { IControlsListProps } from './IControlsListProps'
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
import { MyLink } from './Link'
import { Tabs } from './Tabs'
import { Toolbar } from './Toolbar'
import { MyNav } from './Nav'
import { Grid } from './Grid'
import { Icon } from './Icon'
import { Message } from './Message'
import { MyDialog } from './Dialog'

export const ControlsList: React.FunctionComponent<IControlsListProps> = ({ controls, parentDisabled }) => {

    //console.log(`render ControlsList:`, controls);

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
        'button': Button,
        'stack': MyStack,
        'tabs': Tabs,
        'toolbar': Toolbar,
        'nav': MyNav,
        'grid': Grid,
        'message': Message,
        'dialog': MyDialog,
    }

    const renderChild = (control: any) => {
        if (control.visible === "false") {
            return null;
        }
        const ControlType = controlTypes[control.t];
        if (ControlType === null) {
            console.log(`Unknown control type: ${control.t}`)
        }
        return <ControlType key={control.i} control={control} parentDisabled={parentDisabled} />
    }

    return controls.map(renderChild);
}