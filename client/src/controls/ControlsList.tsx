import React from 'react'
import { IControlsListProps } from './IControlsListProps'
import { MyStack } from './Stack'
import { Textbox } from './Textbox'
import { MyDropdown } from './Dropdown'
import { MyChoiceGroup } from './ChoiceGroup'
import { MyCheckbox } from './Checkbox'
import { Progress } from './Progress'
import { MySpinner } from './Spinner'
import { Button } from './Button'
import { MyText } from './Text'
import { Tabs } from './Tabs'
import { Toolbar } from './Toolbar'
import { MyNav } from './Nav'

export const ControlsList: React.FunctionComponent<IControlsListProps> = ({ controls, parentDisabled }) => {

    //console.log(`render ControlsList: ${id}`);

    const controlTypes: any = {
        'textbox': Textbox,
        'checkbox': MyCheckbox,
        'dropdown': MyDropdown,
        'choicegroup': MyChoiceGroup,
        'progress': Progress,
        'spinner': MySpinner,
        'text': MyText,
        'button': Button,
        'stack': MyStack,
        'tabs': Tabs,
        'toolbar': Toolbar,
        'nav': MyNav,
    }

    const renderChild = (control: any) => {
        if (control.visible === "false") {
            return null;
        }
        const ControlType = controlTypes[control.t];
        return <ControlType key={control.i} control={control} parentDisabled={parentDisabled} />
    }

    return controls.map(renderChild);
}