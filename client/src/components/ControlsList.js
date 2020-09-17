import React from 'react'
import Row from './Row'
import Col from './Col'
import Textbox from './Textbox'
import Button from './Button'
import Label from './Label'

const ControlsList = ({ controls }) => {

    //console.log(`render ControlsList: ${id}`);

    const controlTypes = {
        'row': Row,
        'col': Col,
        'textbox': Textbox,
        'label': Label,
        'button': Button
    }

    const renderChild = control => {
        const ControlType = controlTypes[control.t];
        return <ControlType key={control.i} control={control} />
    }

    return controls.map(renderChild);
}

export default ControlsList