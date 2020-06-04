import React from 'react'
import { shallowEqual, useSelector, useDispatch } from 'react-redux'
import * as Actions from '../actions'
import NodeList from './NodeList'

const Node = React.memo(({control}) => {

  console.log(`render node: ${control.id}`);

  const id = control.id;
  const parentId = control.parentId;

  const childControls = useSelector(state => control.childIds.map(childId => state[childId]), shallowEqual);
  
  //console.log(childControls);

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

  return (
    <div>
      <span onClick={toggleExpand}>Counter: {control.counter}</span>
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
        {control.expanded ? <NodeList id={id} controls={childControls} /> : ""}
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
})

export default Node