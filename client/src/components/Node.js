import React from 'react'
import { useSelector, useDispatch } from 'react-redux'
import * as Actions from '../actions'

const Node = ({id, parentId}) => {
  const node = useSelector(state => state[id]);
  const dispatch = useDispatch();

  const handleIncrementClick = () => {
    dispatch(Actions.increment(id));
  }

  const handleAddChildClick = e => {
    e.preventDefault()

    const childId = dispatch(Actions.createNode()).nodeId
    dispatch(Actions.addChild(id, childId));
  }

  const handleRemoveClick = e => {
    e.preventDefault()

    dispatch(Actions.removeChild(parentId, id))
    dispatch(Actions.deleteNode(id))
  }

  const toggleExpand = e => {
    dispatch(Actions.toggleExpand(id));
  }

  const renderChild = childId => {
    return (
      <li key={childId}>
        <Node id={childId} parentId={id} />
      </li>
    )
  }

  return (
    <div>
      <span onClick={toggleExpand}>Counter: {node.counter}</span>
      {' '}
      <button onClick={handleIncrementClick}>
        +
      </button>
      {' '}
      {typeof parentId !== 'undefined' &&
        <a href="#" onClick={handleRemoveClick} // eslint-disable-line jsx-a11y/anchor-is-valid
           style={{ color: 'lightgray', textDecoration: 'none' }}>
          ×
        </a>
      }
      <ul>
        {node.expanded ? node.childIds.map(renderChild) : ""}
        <li key="add">
          <a href="#" // eslint-disable-line jsx-a11y/anchor-is-valid
            onClick={handleAddChildClick}
          >
            Add child
          </a>
        </li>
      </ul>
    </div>
  )
}

export default Node