import React from 'react'
import { IControlsListProps } from './IControlsListProps'
import { MyStack } from './Stack'
import { Textbox } from './Textbox'
import { MyDropdown } from './Dropdown'
import { MyCheckbox } from './Checkbox'
import { Progress } from './Progress'
import { Button } from './Button'
import { MyText } from './Text'

export const ControlsList: React.FunctionComponent<IControlsListProps> = ({ controls }) => {

    //console.log(`render ControlsList: ${id}`);

    const controlTypes: any = {
        'textbox': Textbox,
        'checkbox': MyCheckbox,
        'dropdown': MyDropdown,
        'progress': Progress,
        'text': MyText,
        'button': Button,
        'stack': MyStack,
    }

    const renderChild = (control: any) => {
        const ControlType = controlTypes[control.t];
        return <ControlType key={control.i} control={control} />
    }

    return controls.map(renderChild);
}