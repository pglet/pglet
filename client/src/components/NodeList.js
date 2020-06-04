import React from 'react'
import Node from './Node'

const NodeList = React.memo(({ id, controls }) => {

    //console.log(`render NodeList: ${id}`);

    const controlTypes = {
        'Node': Node
    }

    const renderChild = control => {
        const ControlType = controlTypes[control.type];
        return (
            <li key={control.id}>
                <ControlType control={control} />
            </li>
        )
    }

    return controls.map(renderChild);
})

export default NodeList