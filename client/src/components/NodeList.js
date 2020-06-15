import React from 'react'
import Node from './Node'
import Row from './Row'
import Column from './Column'
import Text from './Text'
import Button from './Button'

const NodeList = ({ controls }) => {

    //console.log(`render NodeList: ${id}`);

    const controlTypes = {
        'Node': Node,
        'Row': Row,
        'Column': Column,
        'Text': Text,
        'Button': Button
    }

    const renderChild = control => {
        const ControlType = controlTypes[control.t];
        return <ControlType key={control.i} control={control} />
    }

    return controls.map(renderChild);
}

export default NodeList