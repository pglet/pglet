import { INCREMENT, TOGGLE_EXPAND, ADD_CHILD, REMOVE_CHILD, CREATE_NODE, DELETE_NODE } from '../actions/pageActions'

const initialState = {
  "name": "test-1",
  "controls": {
    "0": {
      "c": [
        "1"
      ],
      "i": "0",
      "p": "",
      "t": "Page"
    },
    "1": {
      "c": [
        "2",
        "3"
      ],
      "i": "1",
      "p": "0",
      "t": "Row"
    },
    "2": {
      "c": [
        "4",
        "myBtn"
      ],
      "i": "2",
      "p": "1",
      "t": "Column"
    },
    "3": {
      "c": [
        "5"
      ],
      "i": "3",
      "p": "1",
      "t": "Column"
    },
    "4": {
      "i": "4",
      "p": "2",
      "t": "Text",
      "text": "Hello, world!"
    },
    "5": {
      "i": "5",
      "p": "3",
      "t": "Button",
      "text": "Click me!"
    },
    "myBtn": {
      "i": "myBtn",
      "p": "2",
      "t": "Button",
      "text": "Cancel"
    }
  }
}

const childIds = (state, action) => {
  switch (action.type) {
    case ADD_CHILD:
      return [ ...state, action.childId ]
    case REMOVE_CHILD:
      return state.filter(id => id !== action.childId)
    default:
      return state
  }
}

const node = (state, action) => {
  switch (action.type) {
    case CREATE_NODE:
      return {
        i: action.nodeId,
        p: action.parentId,
        t: 'Node',
        c: [],
        counter: 0,
        expanded: true
      }
    case INCREMENT:
      return {
        ...state,
        counter: state.counter + 1
      }
    case TOGGLE_EXPAND:
      return {
        ...state,
        expanded: !state.expanded
      }
    case ADD_CHILD:
    case REMOVE_CHILD:
      return {
        ...state,
        c: childIds(state.c, action)
      }
    default:
      return state
  }
}

const getAllDescendantIds = (state, nodeId) => (
  state[nodeId].c.reduce((acc, childId) => (
    [ ...acc, childId, ...getAllDescendantIds(state, childId) ]
  ), [])
)

const deleteMany = (state, ids) => {
  state = { ...state }
  ids.forEach(id => delete state.controls[id])
  return state
}

export default (state = initialState, action) => {

  console.log(action);
  // console.log(state);

  const { nodeId } = action
  if (typeof nodeId === 'undefined') {
    return state
  }

  if (action.type === DELETE_NODE) {
    const descendantIds = getAllDescendantIds(state.controls, nodeId)
    //console.log(descendantIds);
    return deleteMany(state, [ nodeId, ...descendantIds ])
  }

  state.controls[nodeId] = node(state.controls[nodeId], action);

  //console.log(state);

  return state;
}
