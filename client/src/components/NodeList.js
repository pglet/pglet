import React from 'react'
import Node from './Node'

const NodeList = ({ id, controls }) => {

    //console.log(`render NodeList: ${id}`);

    const controlTypes = {
        'Node': Node
    }

    const renderChild = control => {
        const ControlType = controlTypes[control.t];
        return <ControlType key={control.i} control={control} />
    }

    return controls.map(renderChild);
}

export default NodeList