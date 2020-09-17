import React from 'react'
import { shallowEqual, useSelector, useDispatch } from 'react-redux'
import { increment, createNode, addChild, removeChild, deleteNode, toggleExpand } from '../features/page/pageSlice'
import ControlsList from './ControlsList'

const Node = React.memo(({control}) => {

  //console.log(`render node: ${control.id}`);

  const id = control.i;
  const parentId = control.p;

  const childControls = useSelector(state => control.c.map(childId => state.page.controls[childId]), shallowEqual);
  
  //console.log(childControls);

  const dispatch = useDispatch();

  const handleIncrementClick = () => {
    dispatch(increment(id));
  }

  const handleAddChildClick = e => {
    e.preventDefault()

    const childId = dispatch(createNode(id)).nodeId
    dispatch(addChild({
      nodeId: id,
      childId
    }));
  }

  const handleRemoveClick = e => {
    e.preventDefault()

    dispatch(removeChild({
      nodeId: parentId,
      childId: id
    }))
    dispatch(deleteNode(id))
  }

  const toggleExpandHandler = e => {
    dispatch(toggleExpand(id));
  }

  return (
    <div>
      <span onClick={toggleExpandHandler}>Counter: {control.counter}</span>
      {' '}
      <button onClick={handleIncrementClick}>
        +
      </button>
      {' '}
      {parentId !== null &&
        <a href="#" onClick={handleRemoveClick} // eslint-disable-line jsx-a11y/anchor-is-valid
           style={{ color: 'lightgray', textDecoration: 'none' }}>
          ×
        </a>
      }
      <ul>
        {control.expanded ? <ControlsList id={id} controls={childControls} /> : ""}
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