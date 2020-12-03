import React from 'react'
import MyStack from './Stack'
import { Textbox } from './Textbox'
import { MyDropdown } from './Dropdown'
import { MyCheckbox } from './Checkbox'
import { Progress } from './Progress'
import Button from './Button'
import Text from './Text'

const ControlsList = ({ controls }) => {

    //console.log(`render ControlsList: ${id}`);

    const controlTypes = {
        'textbox': Textbox,
        'checkbox': MyCheckbox,
        'dropdown': MyDropdown,
        'progress': Progress,
        'text': Text,
        'button': Button,
        'stack': MyStack,
    }

    const renderChild = control => {
        const ControlType = controlTypes[control.t];
        return <ControlType key={control.i} control={control} />
    }

    return controls.map(renderChild);
}

export default ControlsList